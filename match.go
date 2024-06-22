package ll1

import (
	"fmt"
	"unicode/utf8"
)

func MatchString(e Expr, s string) bool {
	switch e := e.(type) {
	case Empty:
		return s == ""
	case Byte:
		return len(s) == 1 && s[0] == e.bv
	case Rune:
		r, size := utf8.DecodeRuneInString(s)
		return r == e.rv && size == len(s)
	case Range:
		if len(s) == 0 {
			return false
		}
		switch lo := e.lo.(type) {
		case Byte:
			if s[0] < lo.bv {
				return false
			}
		case Rune:
			r, _ := utf8.DecodeRuneInString(s)
			if r < lo.rv {
				return false
			}
		default:
			panic(fmt.Errorf("invalid range lo: %w", ErrInvalidArgument))
		}
		switch hi := e.hi.(type) {
		case Byte:
			if s[0] > hi.bv {
				return false
			}
		case Rune:
			r, _ := utf8.DecodeRuneInString(s)
			if r > hi.rv {
				return false
			}
		default:
			panic(fmt.Errorf("invalid range hi: %w", ErrInvalidArgument))
		}
		return true
	case Token:
		return s == e.text
	default:
		panic("not implemented")
	}
}

func (g *Grammar) MatchEmpty(e Expr) bool {
	v := map[string]bool{}
	empty, _ := g.matchEmptyVisitor(v, e)
	return empty
}

func (g *Grammar) matchEmptyVisitor(visited map[string]bool, e Expr) (empty, ok bool) {
	switch e := e.(type) {
	case Empty:
		return true, true
	case Byte, Rune, Range, RepT, OptT:
		return false, true
	case Opt:
		return g.matchEmptyVisitor(visited, e.body)
	case Rep:
		return g.matchEmptyVisitor(visited, e.body)
	case AltT:
		for _, e := range e.alt.body {
			if empty, ok := g.matchEmptyVisitor(visited, e); ok && empty {
				return true, true
			}
		}
		return false, true
	case Seq:
		for _, e := range e.elems {
			if empty, ok := g.matchEmptyVisitor(visited, e); ok && !empty {
				return false, true
			}
		}
		return true, true
	case Name:
		if _, ok := visited[e.id]; ok {
			return false, false
		}
		visited[e.id] = true
		return g.matchEmptyVisitor(visited, g.prods[e.id].expr)
	default:
		return false, true
	}
}
