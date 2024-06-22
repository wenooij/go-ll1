package ll1

import (
	"fmt"

	"golang.org/x/exp/ebnf"
)

type Rep struct{ body Expr }

func (r Rep) Clone() Expr { return Rep{body: r.body.Clone()} }
func (r Rep) Equal(other Expr) bool {
	otherRep, ok := other.(Rep)
	return ok && r.EqualRep(otherRep)
}
func (r Rep) EqualRep(other Rep) bool { return r.body.Equal(other.body) }
func (Rep) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	opt, ok := expr.(*ebnf.Repetition)
	if !ok {
		return nil, fmt.Errorf("input must be a repitition: %v", ErrInvalidArgument)
	}
	body, err := NewFromEBNF(opt.Body)
	if err != nil {
		return nil, err
	}
	if _, ok := body.(Terminal); ok { // Simplify: <Rep!(terminal)> => <RepT!(terminal)>.
		return RepT{Rep{body}}, nil
	}
	return Rep{body}, nil
}
func (r Rep) String() string { return fmt.Sprintf("{%s}", r.body) }

// RepT is a Rep where the body is a Terminal by construction.
type RepT struct{ rep Rep }

func (r RepT) Clone() Expr { return RepT{r.rep.Clone().(Rep)} }
func (r RepT) Equal(other Expr) bool {
	otherRepT, ok := other.(RepT)
	return ok && r.rep.EqualRep(otherRepT.rep)
}
func (RepT) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	rep, err := Rep{}.NewFromEBNF(expr)
	if err != nil {
		return nil, err
	}
	repT, ok := rep.(RepT)
	if !ok {
		return nil, fmt.Errorf("rep body must be a terminal: %w", ErrInvalidArgument)
	}
	return repT, nil
}

func (r RepT) String() string             { return r.rep.String() }
func (RepT) terminal()                    {}
func (r RepT) templateArgs() templateArgs { return r.rep.body.(Terminal).templateArgs() }
