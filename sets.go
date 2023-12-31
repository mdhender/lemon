// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import "github.com/mdhender/lemon/internal/sets"

// FindFirstSets finds all non-terminals which will generate lambda (a/k/a the empty string).
// It then goes back and compute the first sets of every non-terminal. This is
// repeated until no more first sets are found.
//
// The first set is the set of all terminal symbols which can begin a string
// generated by that non-terminal.
func FindFirstSets(lemp *lemon) {
	// initialize all lambdas to false
	for i := 0; i < lemp.nsymbol; i++ {
		lemp.symbols[i].lambda = false
	}
	// create a set for each symbol that will hold a flag for the first set of each terminal.
	for i := lemp.nterminal; i < lemp.nsymbol; i++ {
		lemp.symbols[i].firstset = sets.New(lemp.nterminal + 1)
	}

	// first compute all lambdas.
	for foundLambda := true; foundLambda; {
		foundLambda = false
		for rp := lemp.rule; rp != nil; rp = rp.next {
			if rp.lhs.lambda {
				continue
			}
			var i int
			for i = 0; i < rp.nrhs; i++ {
				sp := rp.rhs[i]
				if !(sp.type_ == NONTERMINAL || sp.lambda == false) {
					panic("assert(sp.type == NONTERMINAL || sp.lambda == false)")
				}
				if sp.lambda == false {
					break
				}
			}
			if i == rp.nrhs {
				foundLambda, rp.lhs.lambda = true, true
			}
		}
	}

	// now compute all first sets. repeat this until a pass completes with no sets added.
	for setsAdded := true; setsAdded; {
		setsAdded = false
		for rp := lemp.rule; rp != nil; rp = rp.next {
			s1 := rp.lhs
			for i := 0; i < rp.nrhs; i++ {
				s2 := rp.rhs[i]
				if s2.type_ == TERMINAL {
					if s1.firstset.Add(s2.index) {
						setsAdded = true
					}
					break
				} else if s2.type_ == MULTITERMINAL {
					for j := 0; j < s2.nsubsym; j++ {
						if s1.firstset.Add(s2.subsym[j].index) {
							setsAdded = true
						}
					}
					break
				} else if s1 == s2 {
					if s1.lambda == false {
						break
					}
				} else {
					if s1.firstset.Union(s2.firstset) {
						setsAdded = true
					}
					if s2.lambda == false {
						break
					}
				}
			}
		}
	}
}
