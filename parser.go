package ll1

import (
	"bytes"
	"fmt"
	"go/format"
	"text/template"

	"golang.org/x/exp/ebnf"
)

const utf8ExtraImport = "unicode/utf8"

var parserTmpl = template.Must(template.New("").ParseFiles("parser.go.tmpl"))

type tmpl struct {
	goBuildTags  []string // GoBuildTags.
	packageName  string   // PackageName (expected "main").
	start        string   // Start symbol name.
	typePrefix   string   // TypePrefix.
	extraImports []string // ExtraImports packages.
	table        []tableRow
	terminals    []Terminal
	names        []Name
}

type tmplOptions struct {
	GoBuildTags []string // GoBuildTags.
	ByteNames   map[byte]string
	RuneNames   map[rune]string
	TokenNames  map[string]string
	RangeNames  map[struct{ Lo, Hi rune }]string
}

func newTemplate(grammar ebnf.Grammar, start string, opts *tmplOptions) (*tmpl, error) {
	t := &tmpl{
		goBuildTags: opts.GoBuildTags,
		packageName: "main",
		start:       start,
		typePrefix:  "symbol",
	}
	g, err := NewGrammarFromEBNF(grammar)
	if err != nil {
		return nil, err
	}
	terminals, err := g.terminals(start, true)
	if err != nil {
		return nil, err
	}
	t.terminals = terminals
	names, err := g.names(start, true)
	if err != nil {
		return nil, err
	}
	t.names = names
	return t, nil
}

func (t *tmpl) ExecuteTemplate() ([]byte, error) {
	var buf bytes.Buffer
	buf.Grow(3000) // wc -c parser.go.tmpl + C
	if err := parserTmpl.ExecuteTemplate(&buf, "parser.go.tmpl", t); err != nil {
		return nil, err
	}
	return format.Source(buf.Bytes())
}

type templateArgs struct {
	parserGoVars []parserGoVar
	lexCases     []lexCase
}

type parserGoVar struct {
	src     string
	comment string
}

type lexCase struct {
	terminal    Terminal
	src         string
	comment     string
	advance     int
	advanceRune bool
}

type tableRow struct {
	key  string
	cols []tableCol
}

type tableCol struct {
	key string
	val int
}

type tmplName struct {
	index int
	name  string
	rhs   []int
	expr  Expr
}

func (n tmplName) RhsArg(i int) string {
	return ""
}

func (n tmplName) String() string {
	return fmt.Sprintf("%d. %s = %s", n.index, n.name, n.expr.String())
}
