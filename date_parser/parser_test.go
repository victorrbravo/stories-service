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

func TestExtractTerms_CtgFunction(t *testing.T) {
	res, err := ExtractTerms(`ctg- food,drink,travel`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d: %+v", len(res.Functions), res.Functions)
	}
	f := res.Functions[0]
	if f.Name != "ctg" {
		t.Errorf("expected function name 'ctg', got %q", f.Name)
	}
	want := []string{"food", "drink", "travel"}
	if len(f.Params) != len(want) {
		t.Fatalf("expected %d params, got %d: %v", len(want), len(f.Params), f.Params)
	}
	for i, w := range want {
		if f.Params[i] != w {
			t.Errorf("param[%d]: expected %q, got %q", i, w, f.Params[i])
		}
	}
}

func TestExtractTerms_CtgFunctionWithOtherTerms(t *testing.T) {
	res, err := ExtractTerms(`coffee and ctg- food,drink,travel or news`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res.Functions) != 1 {
		t.Fatalf("expected 1 function, got %d: %+v", len(res.Functions), res.Functions)
	}
	f := res.Functions[0]
	if f.Name != "ctg" {
		t.Errorf("expected function name 'ctg', got %q", f.Name)
	}
	if got := f.Params; len(got) != 3 || got[0] != "food" || got[1] != "drink" || got[2] != "travel" {
		t.Errorf("expected params [food drink travel], got %v", got)
	}

	foundCoffee := false
	for _, term := range res.AndTerms {
		if term == "coffee" {
			foundCoffee = true
			break
		}
	}
	if !foundCoffee {
		t.Errorf("expected 'coffee' in AndTerms, got %+v", res.AndTerms)
	}

	foundNews := false
	for _, term := range res.OrTerms {
		if term == "news" {
			foundNews = true
			break
		}
	}
	if !foundNews {
		t.Errorf("expected 'news' in OrTerms, got %+v", res.OrTerms)
	}
}
