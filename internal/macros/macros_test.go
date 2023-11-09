// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"fmt"
	"testing"
)

func TestParseMacroExpression(t *testing.T) {
	variables := make(map[string]string)
	variables["a"] = "true"
	variables["b"] = "true"
	for _, expr := range []string{
		"a",
		"a&b",
		"a&(b|c)",
		"!a",
		"!a&!b",
		"(a)&!(b)",
		"a&b&!(!(!(c)))",
	} {
		// Convert the postfix expression to infix
		infixExpression, err := parseMacroExpression([]byte(expr), variables)
		if err != nil {
			fmt.Printf("%q: error: %+v\n", expr, err)
		} else {
			// Print the infix expression
			fmt.Printf("%q: %q\n", expr, string(infixExpression))
		}
	}
}
