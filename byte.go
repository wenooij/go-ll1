package ll1

import (
	"fmt"

	"golang.org/x/exp/ebnf"
)

type Byte struct{ bv byte }

func (b Byte) Clone() Expr { return Byte{bv: b.bv} }
func (b Byte) Equal(other Expr) bool {
	otherByte, ok := other.(Byte)
	return ok && b.bv == otherByte.bv
}
func (Byte) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	token, ok := expr.(*ebnf.Token)
	if !ok {
		return nil, fmt.Errorf("input is not a token: %w", ErrInvalidArgument)
	}
	return Byte{}.NewFromToken(token)
}
func (Byte) NewFromToken(token *ebnf.Token) (Expr, error) {
	switch len(token.String) {
	case 0: // Simplify: Use Empty where possible.
		return Empty{}.NewFromToken(token)
	case 1:
		return Byte{bv: token.String[0]}, nil
	default:
		return nil, fmt.Errorf("input is not a valid byte token: %w", ErrInvalidArgument)
	}
}
func (b Byte) String() string   { return fmt.Sprintf("%q", string(b.bv)) }
func (b Byte) GoString() string { return fmt.Sprintf("%#v", b.bv) }
func (b Byte) terminal()        {}
func (b Byte) Byte() byte       { return b.bv }
func (b Byte) Rune() rune       { return rune(b.bv) }
func (b Byte) templateArgs() templateArgs {
	return templateArgs{
		lexCases: []lexCase{{
			src:     fmt.Sprintf("s[0] == %q", b.bv),
			comment: string(b.bv),
		}},
	}
}
