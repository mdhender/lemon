// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// global variables that are updated via command line switches
var (
	// azDefine is the symbol table for macro definitions
	azDefine map[string]string
	nDefine  int
	// outputDir is the path to create output files in
	outputDir = "."
	// user_templatename is the name of the template file to load
	user_templatename string
)
