package ll1

import (
	"fmt"
	"unicode/utf8"

	"golang.org/x/exp/ebnf"
)

// Range represents a closed interval between Bytes and Runes.
// lo and hi may either Bytes and Runes or both.
type Range struct {
	lo Terminal
	hi Terminal
}

func (r Range) Clone() Expr {
	return Range{
		lo: r.lo.Clone().(Terminal),
		hi: r.lo.Clone().(Terminal),
	}
}
func (r Range) Equal(other Expr) bool {
	otherRange, ok := other.(Range)
	return ok && r.lo.Equal(otherRange.lo) && r.hi.Equal(otherRange.hi)
}

func (Range) NewFromEBNF(expr ebnf.Expression) (Expr, error) {
	rng, ok := expr.(*ebnf.Range)
	if !ok {
		return nil, fmt.Errorf("input must be ebnf Range: %w", ErrInvalidArgument)
	}
	return Range{}.NewFromRange(rng)
}
func (Range) NewFromRange(rng *ebnf.Range) (Expr, error) {
	lo, err := Rune{}.NewFromToken(rng.Begin)
	if err != nil {
		return nil, err
	}
	hi, err := Rune{}.NewFromToken(rng.End)
	if err != nil {
		return nil, err
	}
	loOrd := lo.(interface{ Rune() rune }).Rune()
	hiOrd := hi.(interface{ Rune() rune }).Rune()
	switch {
	case loOrd > hiOrd:
		return nil, fmt.Errorf("invalid range: %w", ErrInvalidArgument)
	case loOrd == hiOrd: // Simplify: a … a => a.
		return lo, nil
	}
	return Range{lo.(Terminal), hi.(Terminal)}, nil
}
func (r Range) Each(f func(rune)) {
	for rv := r.lo.(interface{ Rune() rune }).Rune(); rv <= r.hi.(interface{ Rune() rune }).Rune(); rv++ {
		f(rv)
	}
}
func (r Range) String() string { return fmt.Sprintf("%s … %s", r.lo, r.hi) }
func (Range) terminal()        {}
func (r Range) templateArgs() templateArgs {
	_, loRune := r.lo.(Rune)
	_, hiRune := r.hi.(Rune)
	switch {
	case loRune && hiRune:
		advanceLo := utf8.RuneLen(r.lo.(interface{ Rune() rune }).Rune())
		advanceHi := utf8.RuneLen(r.hi.(interface{ Rune() rune }).Rune())
		return templateArgs{
			lexCases: []lexCase{{
				src:         fmt.Sprintf("s[:%d] >= %#v && s[:%d] <= %#v", advanceLo, r.lo, advanceHi, r.hi),
				comment:     r.String(),
				advanceRune: true,
			}},
		}
	case !loRune && hiRune:
		advanceHi := utf8.RuneLen(r.hi.(interface{ Rune() rune }).Rune())
		return templateArgs{
			lexCases: []lexCase{{
				src:     fmt.Sprintf("s[0] >= %s && s[0] <= 0x80", r.lo),
				comment: fmt.Sprintf("%s … 0x80", r.lo),
				advance: 1,
			}, {
				src:         fmt.Sprintf("s[0] >= 0x80 && s[:%d] < %s", advanceHi, r.hi),
				comment:     fmt.Sprintf("0x80 … %s", r.hi),
				advanceRune: true,
			}},
		}
	case !loRune: // && !hiRune
		return templateArgs{
			lexCases: []lexCase{{
				src:     fmt.Sprintf("s[0] >= %#v && s[0] <= %#v", r.lo, r.hi),
				comment: r.String(),
			}},
		}
	default:
		panic(ErrInvalidArgument)
	}
}
