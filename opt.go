package ll1

import (
	"fmt"

	"golang.org/x/exp/ebnf"
)

type Opt struct{ body Expr }

func (o Opt) Clone() Expr { return Opt{body: o.body.Clone()} }
func (o Opt) Equal(other Expr) bool {
	otherOpt, ok := other.(Opt)
	return ok && o.EqualOpt(otherOpt)
}
func (o Opt) EqualOpt(other Opt) bool { return o.body.Equal(other.body) }
func (Opt) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	opt, ok := expr.(*ebnf.Option)
	if !ok {
		return nil, fmt.Errorf("input must be ebnf option: %v", ErrInvalidArgument)
	}
	return Opt{}.NewFromOption(opt)
}
func (Opt) NewFromOption(opt *ebnf.Option) (Expr, error) {
	body, err := NewFromEBNF(opt.Body)
	if err != nil {
		return nil, err
	}
	return Opt{}.NewFromBody(body)
}
func (Opt) NewFromBody(body Expr) (Expr, error) {
	switch body := body.(type) {
	case Terminal: // Simplify: Opt!<(terminal) -> OptT!<(terminal).
		return OptT{}.NewFromBody(body)
	case Rep: // Simplify: [{x}] -> {x}.
		return body, nil
	case Opt: // Simplify: [[x]] -> [x].
		return body, nil
	}
	if _, ok := body.(Terminal); ok { // Simplify: <Opt!(terminal)> => <OptT!(terminal)>.
		return OptT{Opt{body}}, nil
	}
	return Opt{body}, nil
}
func (o Opt) String() string { return fmt.Sprintf("[%s]", o.body) }

// OptT is an Opt where all elements are Terminal by construction.
type OptT struct{ opt Opt }

func (o OptT) Clone() Expr { return OptT{o.opt.Clone().(Opt)} }
func (o OptT) Equal(other Expr) bool {
	otherOptT, ok := other.(OptT)
	return ok && o.opt.Equal(otherOptT.opt)
}
func (OptT) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	opt, err := Opt{}.NewFromEBNF(expr)
	if err != nil {
		return nil, err
	}
	optT, ok := opt.(OptT)
	if !ok {
		return nil, fmt.Errorf("optional must be comprised only of terminals: %w", ErrInvalidArgument)
	}
	return optT, nil
}
func (o OptT) NewFromBody(body Terminal) (Terminal, error) {
	opt, err := Opt{}.NewFromBody(body)
	if err != nil {
		return nil, err
	}
	return opt.(OptT), nil
}
func (o OptT) String() string             { return o.opt.String() }
func (OptT) terminal()                    {}
func (o OptT) templateArgs() templateArgs { return o.opt.body.(Terminal).templateArgs() }
