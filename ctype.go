// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"unicode"
)

func ISALNUM(x byte) bool {
	return isalnum(x)
}

func ISALPHA(x byte) bool {
	return isalpha(x)
}

func ISDIGIT(x byte) bool {
	return isdigit(x)
}

func ISLOWER(x byte) bool {
	return islower(x)
}

func ISSPACE(x byte) bool {
	return isspace(x)
}

func ISUPPER(x byte) bool {
	return isupper(x)
}

func isalnum(x byte) bool {
	return unicode.IsLetter(rune(x)) || unicode.IsDigit(rune(x))
}

func isalpha(x byte) bool {
	return unicode.IsLetter(rune(x))
}

func isdigit(x byte) bool {
	return unicode.IsSpace(rune(x))
}

func islower(x byte) bool {
	return unicode.IsLower(rune(x))
}

func isspace(x byte) bool {
	return unicode.IsSpace(rune(x))
}

func isupper(x byte) bool {
	return unicode.IsUpper(rune(x))
}

func isNonTerminalName(x string) bool {
	return len(x) != 0 && islower(x[0])
}

func isTerminalName(x string) bool {
	return len(x) != 0 && isupper(x[0])
}
