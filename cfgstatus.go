// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// A configuration is a production rule of the grammar together with
// a mark (dot) showing how much of that rule has been processed so far.
// Configurations also contain a follow-set which is a list of terminal
// symbols which are allowed to immediately follow the end of the rule.
// Every configuration is recorded as an instance of the following:
type cfgstatus int

const (
	COMPLETE cfgstatus = iota
	INCOMPLETE
)

var cfgstatus_name = [...]string{
	COMPLETE:   "COMPLETE",
	INCOMPLETE: "IMCOMPLETE",
}
