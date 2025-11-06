package parser

import (
	"fmt"
	"log"
	"strings"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

// Define the lexer using participle's expected [][]string format
var logicalLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Keyword", `(?i)\bAND\b|\bOR\b|\bNOT\b`},
	{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`},
	{"Whitespace", `[ \t\n\r]+`},
	{"Punct", `[\(\)]`},
	{"Operator", `[<>]=?|==|!=`},
	{"Date", `[A-Za-z]{3}[0-9]{1,2}`}, // e.g., Jun13
	{"Quoted", `"[a-zA-Z_][a-zA-Z0-9_\s]*"`},
})

// Grammar struct definitions
type Expression struct {
	Or []*OrExpr `@@ ( "OR" @@ )*`
}

type OrExpr struct {
	And []*AndExpr `@@ ( "AND" @@ )*`
}

type AndExpr struct {
	Not        *NotExpr    `@@`
	Comparison *Comparison `| @@`
}

type NotExpr struct {
	Not  bool  `( @"NOT" )?`
	Term *Term `@@`
}

type Comparison struct {
	Op   string `@Operator`
	Date string `@Date`
}

type Term struct {
	Ident  string      `  @Ident`
	Nested *Expression `| "(" @@ ")"`
}

// Parser instance

type ParserResult struct {
	Ands    []string
	Ors     []string
	Nots    []string
	Greater string
	Less    string
}

func ExtractTerms(expression string) (ParserResult, error) {
	var parser = participle.MustBuild[Expression](
		participle.Lexer(logicalLexer),
		participle.Elide("Whitespace"),        // ✅ This line makes whitespace ignored!
		participle.CaseInsensitive("Keyword"), // Optional: make AND/OR/NOT case-insensitive
	)
	input := expression

	ast, err := parser.ParseString("", input)
	if err != nil {
		return ParserResult{}, err
	}

	fmt.Println("Parsed successfully.")
	p := extractSimpleTerms(ast, 0)
	return p, nil
}

func extractSimpleTerms(e *Expression, indent int) ParserResult {
	AndTerms := []string{}
	AndNotTerms := []string{}
	var OrTerms []string
	var greater string
	var less string

	var terms *[]string
	for i, or := range e.Or {
		terms = &OrTerms
		if i > 0 {
			fmt.Println(strings.Repeat("  ", indent) + "OR")
		}
		for j, and := range or.And {
			if j > 0 {
				fmt.Println(strings.Repeat("  ", indent+1) + "AND")
				OrTerms = append(OrTerms)
				terms = &AndTerms
			}
			if and.Not.Not {
				fmt.Println(strings.Repeat("  ", indent+2) + "NOT")
				terms = &AndNotTerms
			} else if and.Comparison != nil {
				if and.Comparison.Op == ">" {
					greater = and.Comparison.Date
				} else if and.Comparison.Op == "<" {
					less = and.Comparison.Date
				}
			}
			printTerm(and.Not.Term, indent+3, terms)
		}
	}
	if len(AndTerms) < 2 {
		if len(OrTerms) > 0 {
			AndTerms = append(AndTerms, OrTerms[0])
		} else {
			log.Println("Error AndTerms ONLY ONE OR ZERO TERM: ", AndTerms)
		}
	}
	p := ParserResult{
		Ands:    AndTerms,
		Ors:     OrTerms,
		Nots:    AndNotTerms,
		Greater: greater,
		Less:    less,
	}
	return p
}

func printTerm(t *Term, indent int, terms *[]string) {
	if t.Ident != "" {
		fmt.Println(strings.Repeat("  ", indent)+"TERM:", t.Ident)
		*terms = append(*terms, t.Ident)
	} else if t.Nested != nil {
		fmt.Println(strings.Repeat("  ", indent) + "(nested)")
		extractSimpleTerms(t.Nested, indent+1)
	}
}
