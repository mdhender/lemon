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
	for i, expr := range []string{
		"a",
		"a && b",
		"a && (b || c)",
		"!a",
		"!a && !b",
		"( a ) && ! ( b )",
		"a || b || (((c)))",
		"!(c)",
		"!(!c)",
		"!(!(c))",
	} {
		fmt.Printf("%2d: %q\n", i+1, expr)
		result, err := parseExpression([]byte(expr))
		if err != nil {
			fmt.Printf("%2d: --- error: %+v\n", i+1, err)
		} else {
			// print the result
			for _, tok := range result {
				fmt.Printf("%2d: ... token %s\n", i+1, tok)
			}
		}
	}
}
