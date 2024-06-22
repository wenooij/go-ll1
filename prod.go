package ll1

import (
	"fmt"

	"golang.org/x/exp/ebnf"
)

type Prod struct {
	g    *Grammar
	name Name
	expr Expr
}

func (g *Grammar) newProdFromProduction(prod *ebnf.Production) (*Prod, error) {
	p := &Prod{g: g}
	p.name = Name{prod.Name.String}
	expr, err := NewFromEBNF(prod.Expr)
	if err != nil {
		return nil, err
	}
	p.expr = expr
	if existingProd, ok := g.prods[prod.Name.String]; ok {
		if err := existingProd.merge(p); err != nil {
			return nil, err
		}
	}
	if g.prods == nil {
		g.prods = make(map[string]*Prod)
	}
	g.prods[prod.Name.String] = p
	return p, nil
}

func (p *Prod) merge(other *Prod) error {
	if p.name.id != other.name.id {
		return fmt.Errorf("input name must match: %w", ErrInvalidArgument)
	}
	panic("not implemented")
}
