package ll1

import (
	"fmt"
	"regexp"

	"golang.org/x/exp/ebnf"
)

var validNamePattern = regexp.MustCompile(`^([\p{L}_][\p{L}\p{N}_]*)$`)

type Name struct {
	id string
}

func (n Name) Clone() Expr { return Name{n.id} }
func (n Name) Equal(other Expr) bool {
	otherName, ok := other.(Name)
	return ok && n.id == otherName.id
}
func (Name) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	name, ok := expr.(*ebnf.Name)
	if !ok || !validNamePattern.MatchString(name.String) {
		return nil, fmt.Errorf("input must be valid ebnf Name: %w", ErrInvalidArgument)
	}
	return Name{name.String}, nil
}
func (n Name) String() string { return n.id }
