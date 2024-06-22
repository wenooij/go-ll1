package ll1

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/exp/ebnf"
)

type Expr interface {
	fmt.Stringer
	Clone() Expr // Clone returns a deep clone of the Expr.
	Equal(Expr) bool
	NewFromEBNF(ebnf.Expression) (Expr, error) // NewFromEBNF creates a new simplified Expr from the ebnf Expression.
}

func NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	switch expr := expr.(type) {
	case *ebnf.Token:
		return Token{}.NewFromToken(expr)
	case *ebnf.Range:
		return Range{}.NewFromRange(expr)
	case *ebnf.Option:
		return Opt{}.NewFromOption(expr)
	case *ebnf.Repetition:
		return Rep{}.NewFromEBNF(expr)
	case *ebnf.Name:
		return Name{}.NewFromEBNF(expr)
	case *ebnf.Group: // Simplify: (x) => x.
		return NewFromEBNF(expr.Body)
	case ebnf.Sequence:
		return Seq{}.NewFromSequence(expr)
	case ebnf.Alternative:
		return Alt{}.NewFromAlternative(expr)
	default:
		return nil, fmt.Errorf("not a valid expression %T: %w", expr, ErrInvalidArgument)
	}
}

type Terminal interface {
	Expr
	// len returns the min and max length for terminals.
	// If closed is true max is set to the value of the interval.
	// If closed is false and max is 0 there is no max length limit.
	terminal()
	templateArgs() templateArgs
}

func NewTFromEBNF(e ebnf.Expression) (Terminal, error) {
	switch e := e.(type) {
	case *ebnf.Token:
		var t Terminal
		switch {
		case e.String == "":
			t = Empty{}
		case len(e.String) == 1:
			t = Byte{}
		case utf8.RuneCountInString(e.String) == 1:
			t = Rune{}
		default:
			t = Token{}
		}
		expr, err := t.NewFromEBNF(e)
		if err != nil {
			return nil, err
		}
		return expr.(Terminal), nil
	case *ebnf.Range:
		r, err := Range{}.NewFromEBNF(e)
		if err != nil {
			return nil, err
		}
		return r.(Terminal), nil
	case ebnf.Alternative:
		alt, err := Alt{}.NewFromEBNF(e)
		if err != nil {
			return nil, err
		}
		t, ok := alt.(Terminal)
		if !ok {
			return nil, fmt.Errorf("alternative is not a terminal: %w", ErrInvalidArgument)
		}
		return t, nil
	default: // Unhandled type.
		return nil, fmt.Errorf("input is not a valid ebnf expression: %w", ErrInvalidArgument)
	}
}
