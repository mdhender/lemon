// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// Hash a configuration
func confighash(a *config) (h uint64) {
	return h*571 + uint64(a.rp.index*37+a.dot)
}
