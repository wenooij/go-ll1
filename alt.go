package ll1

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/exp/ebnf"
	"golang.org/x/text/unicode/rangetable"
)

type Alt struct{ body []Expr }

func (a Alt) Clone() Expr {
	res := &Alt{body: make([]Expr, 0, len(a.body))}
	for _, e := range a.body {
		res.body = append(res.body, e.Clone())
	}
	return res
}
func (a Alt) Equal(other Expr) bool {
	otherAlt, ok := other.(Alt)
	return ok && a.EqualAlt(otherAlt)
}
func (a Alt) EqualAlt(other Alt) bool {
	if len(a.body) != len(other.body) {
		return false
	}
	for i := range a.body {
		if !a.body[i].Equal(other.body[i]) {
			return false
		}
	}
	return true
}

// NewFromEBNF creates a new alternative expression or simplified form.
//
// The type of alternative is determined based on inspecting the contents and can be one of:
//   - `Alt` when contents are not all terminals.
//   - `AltT` when contents are terminal but not all rangeable (contains tokens).
//   - `AltRange` when contents are terminal and rangeable.
//
// A simplified form may be returned:
//   - `Token`
//   - `Rune`
//   - `Byte`
//   - `Empty`
func (Alt) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	switch expr := expr.(type) {
	case ebnf.Alternative:
		return Alt{}.NewFromAlternative(expr)
	default:
		return nil, fmt.Errorf("input must be an alternative: %v", ErrInvalidArgument)
	}
}
func (Alt) NewFromAlternative(alt ebnf.Alternative) (Expr, error) {
	switch len(alt) {
	case 0: // This is ambiguous, return an error.
		return nil, fmt.Errorf("input must be a valid ebnf alternative: %v", ErrInvalidArgument)
	case 1: // Simplify: <alt!(a)> => a.
		return NewFromEBNF(alt[0])
	}
	elems := make([]Expr, 0, len(alt))
	for _, e := range alt { // Simplify: a|b|(c|d) => a|b|c|d.
		v, err := NewFromEBNF(e)
		if err != nil {
			return nil, fmt.Errorf("invalid expression in alternative: %w", err)
		}
		switch v := v.(type) {
		case Alt:
			elems = append(elems, v.body...)
		case AltT:
			elems = append(elems, v.alt.body...)
		default:
			elems = append(elems, v)
		}
	}
	return Alt{}.newFromElemsUnchecked(elems...)
}

func (Alt) NewFromElems(elems ...Expr) (Expr, error) {
	switch len(elems) {
	case 0: // This is ambiguous, return an error.
		return nil, fmt.Errorf("invalid input elements: %v", ErrInvalidArgument)
	case 1: // Simplify: <alt!(a)> => a.
		return elems[0], nil
	}
	return Alt{}.newFromElemsUnchecked(elems...)
}

func (Alt) newFromElemsUnchecked(elems ...Expr) (Expr, error) {
	for _, e := range elems { // Use the most specific expression.
		switch e.(type) {
		case Terminal:
		default:
			return Alt{elems}, nil
		}
	}
	return AltT{Alt{elems}}, nil // Simplify: <Alt!(terminals)> => <AltT!(terminals)>.
}

func (a Alt) String() string {
	var sb strings.Builder
	for i, e := range a.body {
		if i > 0 {
			sb.WriteString(" | ")
		}
		fmt.Fprint(&sb, e)
	}
	return sb.String()
}

// AltT is an Alt where all elements are Terminal by construction.
type AltT struct{ alt Alt }

func (a AltT) Clone() Expr { return AltT{a.alt.Clone().(Alt)} }
func (a AltT) Equal(other Expr) bool {
	otherAltT, ok := other.(AltT)
	return ok && a.alt.EqualAlt(otherAltT.alt)
}

func (AltT) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	alt, err := Alt{}.NewFromEBNF(expr)
	if err != nil {
		return nil, err
	}
	altT, ok := alt.(AltT)
	if !ok {
		return nil, fmt.Errorf("alternative must be comprised only of terminals: %w", ErrInvalidArgument)
	}
	return altT, nil
}
func (AltT) NewFromBody(body ...Terminal) (Terminal, error) {
	expr := make([]Expr, 0, len(body))
	for _, e := range body {
		expr = append(expr, e)
	}
	return AltT{Alt{body: expr}}, nil
}

// RangeTable tries to make a unicode RangeTable for the AltT which is possible when it does not have any Tokens.
// This can lead to optimized lexing for this Terminal.
func (r AltT) RangeTable() (*unicode.RangeTable, error) {
	runes := []rune{}
	rts := []*unicode.RangeTable{}
	for _, e := range r.alt.body {
		switch e := e.(type) {
		case Byte:
			runes = append(runes, rune(e.bv))
		case Rune:
			runes = append(runes, e.rv)
		case Range:
			e.Each(func(r rune) { runes = append(runes, r) })
		default:
			return nil, fmt.Errorf("invalid expression type in alternative: %w", ErrInvalidArgument)
		}
	}
	rt := rangetable.New(runes...)
	if len(rts) > 0 {
		rt = rangetable.Merge(append(rts, rt)...)
	}
	return rt, nil
}

func (a AltT) String() string { return a.alt.String() }
func (AltT) terminal()        {}
func (a AltT) templateArgs() templateArgs {
	var t templateArgs
	for _, e := range a.alt.body {
		args := e.(Terminal).templateArgs()
		t.lexCases = append(t.lexCases, args.lexCases...)
		t.parserGoVars = append(t.parserGoVars, args.parserGoVars...)
	}
	return t
}
