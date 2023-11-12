// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"github.com/mdhender/lemon/internal/sets"
	"sort"
	"unicode"
	"unicode/utf8"
)

type symbol struct {
	name       string      // Name of the symbol
	index      int         // Index number for this symbol
	type_      symbol_type // Symbols are all either TERMINALS or NTs */ // was 'ty
	rule       *rule       // Linked list of rules of this (if an NT)
	fallback   *symbol     // fallback token in case this token doesn't parse
	prec       int         // Precedence if defined (-1 otherwise)
	assoc      e_assoc     // Associativity if precedence is defined
	firstset   *sets.Set   // First-set for all rules of this symbol // was ch
	lambda     bool        // True if NT and can generate an empty string
	useCnt     int         // Number of times used
	destructor string      // Code which executes whenever this symbol is popped from the stack during error processing
	destLineno int         // Line number for start of destructor.  Set to -1 for duplicate destructors.
	datatype   string      // The data type of information held by this object. Only used if type==NONTERMINAL
	dtnum      int         // The data type number.  In the parser, the value stack is a union.  The .yy%d element of this union is the correct data type for this object
	bContent   bool        // True if this symbol ever carries content - if it is ever more than just syntax

	// The following fields are used by MULTITERMINALs only
	nsubsym int       // Number of constituent symbols in the MULTI
	subsym  []*symbol // Array of constituent symbols
}

// create a global symbol table
var x2a = make(map[string]*symbol)

// Symbol_new returns a pointer to the (terminal or nonterminal) named symbol.
// Create a new symbol if this is the first time "x" has been seen.
//
// Note on the index. We assume that symbols are never deleted, so the index
// is simply the position the symbol appeared in the grammar. The first symbol
// is 0, the second 1, etc.
func Symbol_new(name string) *symbol {
	sp := x2a[name]
	if sp == nil {
		sp = &symbol{
			name:   name,
			index:  len(x2a),
			prec:   -1,
			assoc:  UNK,
			lambda: LEMON_FALSE,
		}
		r, _ := utf8.DecodeRuneInString(name)
		if unicode.IsUpper(r) {
			sp.type_ = TERMINAL
		} else {
			sp.type_ = NONTERMINAL
		}
		x2a[name] = sp
	}
	sp.useCnt++
	return sp
}

func Symbol_arrayOf() []*symbol {
	var symbols []*symbol
	for _, sym := range x2a {
		symbols = append(symbols, sym)
	}
	//for j := 0; j < len(symbols); j++ {
	//	fmt.Printf("symar: %3d: %3d: %q\n", j, symbols[j].index, symbols[j].name)
	//}
	return symbols
}

func Symbol_count() int {
	return len(x2a)
}

func Symbol_find(name string) *symbol {
	return x2a[name]
}

// Symbol_sortedSlice has the side effect of changing the symbol indexes.
func Symbol_sortedSlice() []*symbol {
	symbols := Symbol_arrayOf()
	sort.Slice(symbols, func(i, j int) bool {
		return symbols[i].less(symbols[j])
	})
	for index, sym := range symbols {
		sym.index = index
	}
	return symbols
}

// typeOf returns an integer 1..3 based on the type of symbol.
// Multi-terminals have a type of 3, non-terminals 2, and terminals are 1.
func (s *symbol) typeOf() int {
	if s.type_ == MULTITERMINAL {
		return 3
	} else if s.name[0] > 'Z' {
		// non-terminal and the "{default}" symbol
		return 2
	}
	// terminal
	return 1
}

/*
** Symbols that begin with upper case letters (terminals or tokens)
** must sort before symbols that begin with lower case letters
** (non-terminals).  And MULTITERMINAL symbols (created using the
** %token_class directive) must sort at the very end. Other than
** that, the order does not matter.
 */
func (s *symbol) less(s2 *symbol) bool {
	a, b := s.typeOf(), s2.typeOf()
	if a > b {
		return false
	} else if a == b {
		// We find experimentally that leaving the symbols in their original
		// order (the order they appeared in the grammar file) gives the
		// smallest parser tables in SQLite.
		// If you don't care, just `return s.name < s2.name`
		return s.index < s2.index
	}
	return true
}
