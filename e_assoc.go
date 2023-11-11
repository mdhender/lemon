// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

type e_assoc int

const (
	LEFT e_assoc = iota
	RIGHT
	NONE
	UNK
)

var e_assoc_names = [...]string{
	LEFT:  "LEFT",
	RIGHT: "RIGHT",
	NONE:  "NONE",
	UNK:   "UNK",
}

func (e e_assoc) String() string {
	return e_assoc_names[e]
}
