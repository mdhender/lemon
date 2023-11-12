// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"io"
	"math/rand"
	"sort"
)

// Each production rule in the grammar is stored in the following
// structure.
type rule struct {
	lhs         *symbol   // Left-hand side of the rule
	lhsalias    string    // Alias for the LHS (NULL if none)
	lhsStart    bool      // True if left-hand side is the start symbol
	ruleline    int       // Line number for the rule
	nrhs        int       // Number of RHS symbols
	rhs         []*symbol // The RHS symbols
	rhsalias    []string  // An alias for each RHS symbol (NULL if none)
	line        int       // Line number at which code begins
	code        string    // The code executed when this rule is reduced
	codePrefix  []byte    // Setup code before code[] above
	codeSuffix  []byte    // Breakdown code after code[] above
	precsym     *symbol   // Precedence symbol for this rule
	index       int       // An index number for this rule
	iRule       int       // Rule number as used in the generated tables
	noCode      bool      // True if this rule has no associated C code
	codeEmitted bool      // True if the code has been emitted already
	canReduce   bool      // True if this rule is ever reduced
	doesReduce  bool      // Reduce actions occur after optimization
	neverReduce bool      // Reduce is theoretically possible, but prevented  by actions or other outside implementation
	nextlhs     *rule     // Next rule with the same LHS
	next        *rule     // Next rule in the global list
}

func (r *rule) length() int {
	n := 0
	for r != nil {
		r, n = r.next, n+1
	}
	return n
}

func (r *rule) shuffle() *rule {
	// temporarily turn the linked list into a flat list
	var list []*rule
	for node := r; node != nil; node = node.next {
		list = append(list, node)
	}
	// shuffle
	rand.Shuffle(len(list), func(i, j int) {
		list[i], list[j] = list[j], list[i]
	})
	// append a nil node to make the re-linking easier
	list = append(list, nil)
	// update the links in the flat list
	for index, node := range list {
		if node != nil {
			node.next = list[index+1]
		}
	}
	// and return the now-shuffles list
	return list[0]
}

// Rule_sort implements a merge sort a list of rules in order of increasing iRule value
func (r *rule) sort() *rule {
	// temporarily turn the linked list into a flat list
	var list []*rule
	for node := r; node != nil; node = node.next {
		list = append(list, node)
	}
	// sort the flat list in order of increasing iRule value
	sort.Slice(list, func(i, j int) bool {
		return list[i].iRule < list[j].iRule
	})
	// append a nil node to make the re-linking easier
	list = append(list, nil)
	// update the links in the flat list
	for index, node := range list {
		if node != nil {
			node.next = list[index+1]
		}
	}
	// and return the now-sorted list
	return list[0]
}

// Reprint duplicates the input file without comments and without actions on rules
func Reprint(w io.Writer, lemp *lemon) {
	fmt.Printf("// Reprint of input file \"%s\".\n// Symbols:\n", lemp.filename)
	maxlen := 10
	for i := 0; i < lemp.nsymbol; i++ {
		sp := lemp.symbols[i]
		if len(sp.name) > maxlen {
			maxlen = len(sp.name)
		}
	}
	ncolumns := 76 / (maxlen + 5)
	if ncolumns < 1 {
		ncolumns = 1
	}

	// print headings as comments
	skip := (lemp.nsymbol + ncolumns - 1) / ncolumns
	for i := 0; i < skip; i++ {
		_, _ = fmt.Fprintf(w, "//")
		for j := i; j < lemp.nsymbol; j += skip {
			sp := lemp.symbols[j]
			if !(sp.index == j) {
				panic("assert(sp.index == j)")
			}
			_, _ = fmt.Fprintf(w, " %3d %-*.*s", j, maxlen, maxlen, sp.name)
		}
		_, _ = fmt.Fprintf(w, "\n")
	}

	// print the rules
	for rp := lemp.rule; rp != nil; rp = rp.next {
		rp.print(w)
		_, _ = fmt.Fprintf(w, ".")
		if rp.precsym != nil {
			_, _ = fmt.Fprintf(w, " [%s]", rp.precsym.name)
		}
		// if rp.code != "" {
		//   fmt.Printf("\n    %s", rp.code)
		//
		_, _ = fmt.Fprintf(w, "\n")
	}
}

// print the text of a rule
func (r *rule) print(out io.Writer) {
	_, _ = fmt.Fprintf(out, "%s", r.lhs.name)
	// if rp.lhsalias != "" {
	//     _, _ = fmt.Fprintf(out, "(%s)", rp.lhsalias)
	// }
	_, _ = fmt.Fprintf(out, " ::=")
	for i := 0; i < r.nrhs; i++ {
		sp := r.rhs[i]
		if sp.type_ == MULTITERMINAL {
			_, _ = fmt.Fprintf(out, " %s", sp.subsym[0].name)
			for j := 1; j < sp.nsubsym; j++ {
				_, _ = fmt.Fprintf(out, "|%s", sp.subsym[j].name)
			}
		} else {
			_, _ = fmt.Fprintf(out, " %s", sp.name)
		}
		// if rp.rhsalias[i] != nil {
		//     _, _ = fmt.Fprintf(out, "(%s)", rp.rhsalias[i])
		// }
	}
}

// print a single rule
func (r *rule) printCursor(fp io.Writer, iCursor int) {
	_, _ = fmt.Fprintf(fp, "%s ::=", r.lhs.name)
	for i := 0; i <= r.nrhs; i++ {
		if i == iCursor {
			_, _ = fmt.Fprintf(fp, " *")
		}
		if i == r.nrhs {
			break
		}
		sp := r.rhs[i]
		if sp.type_ == MULTITERMINAL {
			_, _ = fmt.Fprintf(fp, " %s", sp.subsym[0].name)
			for j := 1; j < sp.nsubsym; j++ {
				_, _ = fmt.Fprintf(fp, "|%s", sp.subsym[j].name)
			}
		} else {
			_, _ = fmt.Fprintf(fp, " %s", sp.name)
		}
	}
}
