// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"testing"
)

func TestParse(t *testing.T) {
	symtab := make(map[string]string)
	symtab["a"] = "true"
	symtab["b"] = "true"

	lem := &lemon{
		filename:          "example.y",
		printPreprocessed: false,
	}
	Parse(lem, symtab)
	if lem.errorcnt != 0 {
		t.Errorf("parse: want 0 errors, got %d\n", lem.errorcnt)
	}
	if lem.nrule == 0 {
		t.Errorf("parse: want >0 rules, got %d\n", lem.nrule)
	}
	lem.errsym = Symbol_find("error")

	/* Count and index the symbols of the grammar */
	Symbol_new("{default}")
	lem.nsymbol = Symbol_count()
	lem.symbols = Symbol_sortedSlice()

	i := lem.nsymbol
	for lem.symbols[i-1].type_ == MULTITERMINAL {
		i--
	}
	if lem.symbols[i-1].name != "{default}" {
		t.Errorf("parse: want symbols[%d] to be %q, got %q\n", i-1, "{default}", lem.symbols[i-1].name)
	}
	lem.nsymbol = i - 1
	for i = 1; isupper(lem.symbols[i].name[0]); i++ {
		//
	}
	lem.nterminal = i

}
