package date_parser

import (
	"testing"
)

func TestExpression_ExtractTerms(t *testing.T) {
	// ... rest of the test code remains the same ...

	//terms := ExtractTerms(ast)
	//fmt.Printf("%+v\n", terms)
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

	// "bill" should still be parsed as a regular AND term.
	foundBill := false
	for _, term := range res.AndTerms {
		if term == "bill" {
			foundBill = true
			break
		}
	}
	if !foundBill {
		t.Errorf("expected 'bill' in AndTerms, got %+v", res.AndTerms)
	}
}

func TestExtractTerms_FunctionAcrossOr(t *testing.T) {
	res, err := ExtractTerms(`foo and max- amount,2 or bar`)
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
