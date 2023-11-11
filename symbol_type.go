// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// Symbols (terminals and nonterminals) of the grammar are stored in the following
type symbol_type int

const (
	TERMINAL symbol_type = iota
	NONTERMINAL
	MULTITERMINAL
)

var symbol_type_names = [...]string{
	TERMINAL:      "TERMINAL",
	NONTERMINAL:   "NONTERMINAL",
	MULTITERMINAL: "MULTITERMINAL",
}

func (st symbol_type) String() string {
	return symbol_type_names[st]
}
