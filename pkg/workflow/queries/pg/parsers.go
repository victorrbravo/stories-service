package pg

import (
	"fmt"
	"github.com/bzick/tokenizer"
)

const (
	TEquality = 1
	TDot      = 2
	TMath     = 3
)

func CheckWhereExpression(s string) {
	parser := tokenizer.New()
	parser.DefineTokens(TEquality, []string{"<", "=", "<=", "==", ">=", ">", "!="})
	parser.DefineTokens(TDot, []string{"."})
	parser.DefineTokens(TMath, []string{"+", "-", "/", "*", "%"})
	//parser.DefineStringToken(tokenizer.TokenKey(1), `"`, " ").SetEscapeSymbol(tokenizer.BackSlash)

	// create tokens stream
	stream := parser.ParseString(s)
	defer stream.Close()

	// iterate over each token
	for stream.IsValid() {
		token := stream.CurrentToken()
		if token.Key() == TEquality {
			fmt.Println("Equality |", token.ValueString(), "|")
		} else {
			fmt.Println("token |", token.ValueString(), "|")
		}

		stream.GoNext()
	}

}

// configure tokenizer
