package ll1

import (
	"fmt"

	"golang.org/x/exp/ebnf"
)

type Empty struct{}

func (Empty) Clone() Expr           { return Empty{} }
func (Empty) Equal(other Expr) bool { _, ok := other.(Empty); return ok }
func (Empty) NewFromEBNF(e ebnf.Expression) (Expr, error) {
	t, ok := e.(*ebnf.Token)
	if !ok {
		return nil, fmt.Errorf("input is not a token: %w", ErrInvalidArgument)
	}
	return Empty{}.NewFromToken(t)
}
func (Empty) NewFromToken(token *ebnf.Token) (Expr, error) {
	if token.String != "" {
		return nil, fmt.Errorf("input must be an empty token: %w", ErrInvalidArgument)
	}
	return Empty{}, nil
}
func (Empty) String() string             { return `""` }
func (Empty) terminal()                  {}
func (Empty) templateArgs() templateArgs { return templateArgs{} }
