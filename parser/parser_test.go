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
		participle.UseLookahead(2),
	)
	input := `term1 OR term2 AND NOT term3`

	ast, err := parser.ParseString("", input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Parsed successfully.")
	extractSimpleTerms(ast, 0)
}

func TestExtractTerms_Function(t *testing.T) {
	res, err := ExtractTerms(`max- amount,2 and bill`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d: %+v", len(res.Functions), res.Functions)
	}
	f := res.Functions[0]
	if f.Name != "max" {
		t.Errorf("expected function name 'max', got %q", f.Name)
	}
	if len(f.Params) != 2 || f.Params[0] != "amount" || f.Params[1] != "2" {
		t.Errorf("expected params [amount 2], got %v", f.Params)
	}
}

func TestExtractTerms_FunctionMixedWithTerms(t *testing.T) {
	res, err := ExtractTerms(`foo AND max- amount,2 OR bar`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d: %+v", len(res.Functions), res.Functions)
	}
	if res.Functions[0].Name != "max" {
		t.Errorf("expected function name 'max', got %q", res.Functions[0].Name)
	}
	if got := res.Functions[0].Params; len(got) != 2 || got[0] != "amount" || got[1] != "2" {
		t.Errorf("expected params [amount 2], got %v", got)
	}
}
