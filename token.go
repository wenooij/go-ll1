package ll1

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/exp/ebnf"
)

type Token struct{ text string }

func (t Token) Clone() Expr { return Token{text: t.text} }
func (t Token) Equal(other Expr) bool {
	otherToken, ok := other.(Token)
	return ok && t.text == otherToken.text
}
func (Token) NewFromEBNF(e ebnf.Expression) (Expr, error) {
	token, ok := e.(*ebnf.Token)
	if !ok {
		return nil, fmt.Errorf("input must be a token: %w", ErrInvalidArgument)
	}
	return Token{}.NewFromToken(token)
}
func (Token) NewFromToken(token *ebnf.Token) (Expr, error) {
	switch {
	case len(token.String) == 0: // Simplify: use Empty where possible.
		return Empty{}.NewFromToken(token)
	case len(token.String) == 1: // Simplify: use Byte where possible.
		return Byte{}.NewFromToken(token)
	case utf8.RuneCountInString(token.String) == 1: // Simplify: use Rune where possible.
		return Rune{}.NewFromToken(token)
	}
	return Token{token.String}, nil
}
func (t Token) String() string { return fmt.Sprintf("%q", t.text) }
func (t Token) terminal()      {}
func (t Token) templateArgs() templateArgs {
	return templateArgs{
		lexCases: []lexCase{{
			src:     fmt.Sprintf("s[:%d] == %q", len(t.text), t.text),
			comment: t.text,
		}},
	}
}
