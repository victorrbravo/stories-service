package parser

import (
	"fmt"
	"testing"

	"github.com/alecthomas/participle/v2"
)

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

func TestExpression_ExtractTerms(t *testing.T) {
	// ... rest of the test code remains the same ...
	var parser = participle.MustBuild[Expression](
		participle.Lexer(logicalLexer),
		participle.Elide("Whitespace"),        // ✅ This line makes whitespace ignored!
		participle.CaseInsensitive("Keyword"), // Optional: make AND/OR/NOT case-insensitive
	)
	input := `term1 OR term2 AND NOT term3`

	ast, err := parser.ParseString("", input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsed successfully.")
	extractSimpleTerms(ast, 0)
}
