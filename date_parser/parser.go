package date_parser

import (
	"fmt"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"time"
)

var Lexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Whitespace", `\s+`},
	{"Operator", `[<>]=?|==|!=`},
	{"And", `(?i)\band\b`},
	{"Or", `(?i)\bor\b`},
	{"Not", `(?i)\bnot\b`},
	{"Date", `[A-Za-z]{3}[0-9]{1,2}`},   // e.g., Jun13
	{"Ident", `[a-zA-Z_][a-zA-Z0-9_]*`}, // e.g., cash, Cafe
})

type Expr struct {
	Or []*OrExpr `@@ ( "or" @@ )*`
}

type OrExpr struct {
	And []*AndExpr `@@ ( "and" @@ )*`
}

type AndExpr struct {
	Comparison *Comparison `  @@`
	Tag        *Tag        `| @@`
}

type Comparison struct {
	Op   string `@Operator`
	Date string `@Date`
}

type Tag struct {
	Not  bool   `[ @"not" ]`
	Name string `@Ident`
}

type ParsedResult struct {
	AndTerms    []string
	OrTerms     []string
	AndNotTerms []string
	LessThan    string
	GreaterThan string
}

func removeString(slice []string, target string) []string {
	var i int
	for _, s := range slice {
		if s != target {
			slice[i] = s
			i++
		}
	}
	return slice[:i]
}
func ExtractTerms(input string) (ParsedResult, error) {
	var parser = participle.MustBuild[Expr](
		participle.Lexer(Lexer),
		participle.Elide("Whitespace"),
		participle.CaseInsensitive("And", "Or", "Not"),
	)

	expr, err := parser.ParseString("", input)
	if err != nil {
		return ParsedResult{}, err
	}
	result := ParsedResult{}

	for i, orExpr := range expr.Or {
		isOrGroup := i > 0 // everything after the first is from an "or" group

		for _, andExpr := range orExpr.And {
			switch {
			case andExpr.Comparison != nil:
				if andExpr.Comparison.Op == "<" {
					result.LessThan = andExpr.Comparison.Date
				} else if andExpr.Comparison.Op == ">" {
					result.GreaterThan = andExpr.Comparison.Date
				}
				if !isOrGroup {
					result.AndTerms = append(result.AndTerms, andExpr.Comparison.Date)
				} else {
					result.OrTerms = append(result.OrTerms, andExpr.Comparison.Date)
				}
			case andExpr.Tag != nil:
				tag := andExpr.Tag.Name
				if andExpr.Tag.Not {
					result.AndNotTerms = append(result.AndNotTerms, tag)
				} else {
					if !isOrGroup {
						result.AndTerms = append(result.AndTerms, tag)
					} else {
						result.OrTerms = append(result.OrTerms, tag)
					}
				}
			}
		}
	}
	if result.LessThan != "" {
		result.AndTerms = removeString(result.AndTerms, result.LessThan)
	}
	if result.GreaterThan != "" {
		result.AndTerms = removeString(result.AndTerms, result.GreaterThan)
	}
	if len(result.AndTerms) == 0 {
		result.AndTerms = nil
	}

	return result, nil
}

func GenerateSQLWhereDate(input string, ope string) (string, error) {
	// Validate input length
	if len(input) < 4 || len(input) > 5 {
		return "", fmt.Errorf("invalid input format")
	}

	// Parse input like "Jun15"
	monthAbbrev := input[:3]
	dayStr := input[3:]

	// Get current year
	year := time.Now().Year()

	// Parse full date string in local (UTC-5) timezone
	dateStr := fmt.Sprintf("%s %s %d", monthAbbrev, dayStr, year)
	loc := time.FixedZone("UTC-5", -5*60*60)

	// Parse date in local timezone
	t, err := time.ParseInLocation("Jan 2 2006", dateStr, loc)
	if err != nil {
		return "", fmt.Errorf("failed to parse date: %v", err)
	}

	// Convert to UTC
	utcTime := t.UTC()
	if ope == "<" {
		utcTime = utcTime.Add(24 * time.Hour)
	}
	// Format to ISO 8601 for SQL
	isoDate := utcTime.Format(time.RFC3339)

	// Construct SQL WHERE clause
	whereClause := isoDate

	return whereClause, nil
}
