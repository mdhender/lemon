// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/lemon/internal/macros"
	"os"
)

// The text in the input is part of the argument to an %ifdef or %ifndef.
// Evaluate the text as a boolean expression.  Return true or false.
//
// input is something like "var && (var || !var)"
func eval_preprocessor_boolean(input []byte, lineno int, symtab map[string]string) bool {
	result, err := macros.EvalExpression(input, symtab)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%%if syntax error on line %d.\n", lineno)
		_, _ = fmt.Fprintf(os.Stderr, "   %v\n", err)
		os.Exit(1)
	}
	return result
}
