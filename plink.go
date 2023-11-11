// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// A followset propagation link indicates that the contents of one
// configuration followset should be propagated to another whenever
// the first changes.
type plink struct {
	cfp  *config // The configuration to which linked
	next *plink  // The next propagate link
}
