package ll1

import (
	"fmt"
	"slices"
	"strings"

	"golang.org/x/exp/ebnf"
)

type Seq struct{ elems []Expr }

func (s Seq) Clone() Expr {
	elems := make([]Expr, 0, len(s.elems))
	for _, e := range s.elems {
		elems = append(elems, e.Clone())
	}
	return Seq{elems}
}
func (s Seq) Equal(other Expr) bool {
	otherSeq, ok := other.(Seq)
	return ok && slices.EqualFunc(s.elems, otherSeq.elems, func(e1, e2 Expr) bool { return e1.Equal(e2) })
}
func (Seq) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	switch expr := expr.(type) {
	case ebnf.Sequence:
		return Seq{}.NewFromSequence(expr)
	default:
		return nil, fmt.Errorf("input must be a sequence: %w", ErrInvalidArgument)
	}
}
func (Seq) NewFromSequence(seq ebnf.Sequence) (Expr, error) {
	switch len(seq) {
	case 0: // Ambiguous case, return error.
		return nil, fmt.Errorf("invalid input sequence: %w", ErrInvalidArgument)
	case 1: // Simplify: seq!<(a)> => a.
		return NewFromEBNF(seq[0])
	}
	var elems []Expr
	for _, e := range seq {
		v, err := NewFromEBNF(e)
		if err != nil {
			return nil, fmt.Errorf("invalid element in sequence: %w", err)
		}
		elems = append(elems, v)
	}
	return Seq{}.newFromElemsUnchecked(elems...)
}
func (Seq) NewFromElems(elems ...Expr) (Expr, error) {
	switch len(elems) {
	case 0: // Ambiguous case, return error.
		return nil, fmt.Errorf("invalid input elements: %w", ErrInvalidArgument)
	case 1: // Simplify: seq!<(a)> => a.
		return elems[0], nil
	}
	return Seq{}.newFromElemsUnchecked(elems...)
}
func (Seq) newFromElemsUnchecked(elems ...Expr) (Expr, error) {
	return Seq{elems}, nil
}
func (s Seq) String() string {
	var sb strings.Builder
	for i, e := range s.elems {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%s", e)
	}
	return sb.String()
}
