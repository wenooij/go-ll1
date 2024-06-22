package ll1

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/exp/ebnf"
)

type Rune struct{ rv rune }

func (r Rune) Clone() Expr { return Rune{rv: r.rv} }
func (r Rune) Equal(other Expr) bool {
	otherRune, ok := other.(Rune)
	return ok && r.rv == otherRune.rv
}
func (Rune) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	token, ok := expr.(*ebnf.Token)
	if !ok {
		return nil, fmt.Errorf("input is not an ebnf rune token: %w", ErrInvalidArgument)
	}
	return Rune{}.NewFromToken(token)
}
func (Rune) NewFromToken(token *ebnf.Token) (Expr, error) {
	switch len(token.String) {
	case 0: // Simplify: use Empty where possible.
		return Empty{}.NewFromToken(token)
	case 1: // Simplify: use Byte where possible.
		return Byte{}.NewFromToken(token)
	}
	rv, size := utf8.DecodeRuneInString(token.String)
	if size == 0 || len(token.String) != size || rv == utf8.RuneError {
		return nil, fmt.Errorf("input is not a valid rune: %w", ErrInvalidArgument)
	}
	return Rune{rv: rv}, nil
}
func (r Rune) String() string   { return fmt.Sprintf("%q", string(r.rv)) }
func (r Rune) GoString() string { return fmt.Sprintf("%#v", r.rv) }
func (Rune) terminal()          {}
func (r Rune) Rune() rune       { return r.rv }
func (r Rune) templateArgs() templateArgs {
	return templateArgs{
		lexCases: []lexCase{{
			src:     fmt.Sprintf("s[:%d] == %q", utf8.RuneLen(r.rv), r.rv),
			comment: string(r.rv),
		}},
	}
}
