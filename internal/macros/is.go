// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"bytes"
	"unicode"
)

// ismacro is a helper function to detect macro definitions in the input grammar.
// a line defines a macro if it has the prefix followed by a space or end of input.
func ismacro(input []byte, macro string) bool {
	if !bytes.HasPrefix(input, []byte(macro)) {
		return false
	} else if len(input) > len(macro) && !unicode.IsSpace(rune(input[len(macro)])) {
		return false
	}
	return true
}
