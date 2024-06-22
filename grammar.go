package ll1

import (
	"fmt"

	"golang.org/x/exp/ebnf"
)

type Grammar struct {
	prods map[string]*Prod
}

func NewGrammarFromEBNF(grammar ebnf.Grammar) (*Grammar, error) {
	g := &Grammar{}
	for _, p := range grammar {
		if _, err := g.newProdFromProduction(p); err != nil {
			return nil, fmt.Errorf("failed to create production %s: %w", p.Name.String, err)
		}
	}
	return g, nil
}

func (g *Grammar) terminals(start string, recursive bool) ([]Terminal, error) {
	var terminals []Terminal

	visitedEmpty := false
	visitedBytes := map[Byte]struct{}{}
	visitedRunes := map[Rune]struct{}{}
	visitedTokens := map[Token]struct{}{}
	visitedRanges := map[Range]struct{}{}
	visitedNames := map[string]struct{}{}

	startProd, ok := g.prods[start]
	if !ok {
		return nil, fmt.Errorf("start does not appear in grammar: %w", ErrInvalidArgument)
	}

	queue := []Expr{startProd.name}

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		switch e := e.(type) {
		case Empty:
			if visitedEmpty {
				continue
			}
			visitedEmpty = true
			terminals = append(terminals, e)
		case Byte:
			if _, ok := visitedBytes[e]; ok {
				continue
			}
			visitedBytes[e] = struct{}{}
			terminals = append(terminals, e)
		case Rune:
			if _, ok := visitedRunes[e]; ok {
				continue
			}
			visitedRunes[e] = struct{}{}
			terminals = append(terminals, e)
		case Token:
			if _, ok := visitedTokens[e]; ok {
				continue
			}
			visitedTokens[e] = struct{}{}
			terminals = append(terminals, e)
		case Range:
			if _, ok := visitedRanges[e]; ok {
				continue
			}
			visitedRanges[e] = struct{}{}
			terminals = append(terminals, e)
		case OptT:
			queue = append(queue, e.opt.body)
		case RepT:
			queue = append(queue, e.rep.body)
		case AltT:
			queue = append(queue, e.alt.body...)
		case Opt:
			queue = append(queue, e.body)
		case Rep:
			queue = append(queue, e.body)
		case Alt:
			queue = append(queue, e.body...)
		case Name:
			if _, ok := visitedNames[e.id]; ok {
				continue
			}
			visitedNames[e.id] = struct{}{}
			if recursive {
				queue = append(queue, g.prods[e.id].expr)
			}
		case Seq:
			queue = append(queue, e.elems...)
		}
	}
	return terminals, nil
}

func (g *Grammar) names(start string, recursive bool) ([]Name, error) {
	var names []Name
	visitedNames := map[string]struct{}{}

	startProd, ok := g.prods[start]
	if !ok {
		return nil, fmt.Errorf("start does not appear in grammar: %w", ErrInvalidArgument)
	}

	queue := []Expr{startProd.name}

	for len(queue) > 0 {
		e := queue[0]
		queue = queue[1:]

		switch e := e.(type) {
		case OptT:
			queue = append(queue, e.opt.body)
		case RepT:
			queue = append(queue, e.rep.body)
		case AltT:
			queue = append(queue, e.alt.body...)
		case Opt:
			queue = append(queue, e.body)
		case Rep:
			queue = append(queue, e.body)
		case Alt:
			queue = append(queue, e.body...)
		case Name:
			if _, ok := visitedNames[e.id]; ok {
				continue
			}
			names = append(names, e)
			visitedNames[e.id] = struct{}{}
			if recursive {
				queue = append(queue, g.prods[e.id].expr)
			}
		case Seq:
			queue = append(queue, e.elems...)
		}
	}
	return names, nil
}

func (g *Grammar) first() map[string][]Terminal {
	first := make(map[string][]Terminal, len(g.prods))
	for hasChanges := true; ; hasChanges = false {
		for name, p := range g.prods {
			_, changes := g.addFirstExpr(first, name, p.expr)
			hasChanges = hasChanges || changes
		}
		if !hasChanges {
			break
		}
	}
	return first
}

func (g *Grammar) addFirstExpr(first map[string][]Terminal, name string, expr Expr) (empty, changes bool) {
	add := func(name string, t Terminal) (empty, changes bool) {
		for _, e := range first[name] {
			if g.MatchEmpty(e) {
				empty = true
			}
			if e.Equal(t) {
				return empty, false
			}
		}
		first[name] = append(first[name], t)
		return empty, true
	}
	switch expr := expr.(type) {
	case Terminal:
		return add(name, expr)
	case Opt:
		hasEmpty1, hasChanges1 := add(name, Empty{})
		hasEmpty2, hasChanges2 := g.addFirstExpr(first, name, expr.body)
		return hasEmpty1 || hasEmpty2, hasChanges1 || hasChanges2
	case Rep:
		hasEmpty1, hasChanges1 := add(name, Empty{})
		hasEmpty2, hasChanges2 := g.addFirstExpr(first, name, expr.body)
		return hasEmpty1 || hasEmpty2, hasChanges1 || hasChanges2
	case Alt:
		var hasEmpty, hasChanges bool
		for _, e := range expr.body {
			empty, changes := g.addFirstExpr(first, name, e)
			hasEmpty = hasEmpty || empty
			hasChanges = hasChanges || changes
		}
		return hasEmpty, hasChanges
	case Seq:
		var hasEmpty, hasChanges bool
		for _, e := range expr.elems {
			empty, changes := g.addFirstExpr(first, name, e)
			if !empty {
				return false, changes
			}
			hasEmpty = true
			hasChanges = hasChanges || changes
		}
		return hasEmpty, hasChanges
	case Name:
		var hasEmpty, hasChanges bool
		for _, t := range first[expr.id] {
			empty, changes := add(name, t)
			hasEmpty = hasEmpty || empty
			hasChanges = hasChanges || changes
		}
		return hasEmpty, hasChanges
	default:
		panic(fmt.Errorf("unexpected Expr %T", expr))
	}
}
