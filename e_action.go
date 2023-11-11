// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

type e_action int

// the order of e_action is important. the order is used when sorting actions.
const (
	SHIFT e_action = iota
	ACCEPT
	REDUCE
	ERROR
	SSCONFLICT  // A shift/shift conflict
	SRCONFLICT  // Was a reduce, but part of a conflict
	RRCONFLICT  // Was a reduce, but part of a conflict
	SH_RESOLVED // Was a shift.  Precedence resolved conflict
	RD_RESOLVED // Was reduce.  Precedence resolved conflict
	NOT_USED    // Deleted by compression
	SHIFTREDUCE // Shift first, then reduce
)

var e_action_names = [...]string{
	SHIFT:       "SHIFT",
	ACCEPT:      "ACCEPT",
	REDUCE:      "REDUCE",
	ERROR:       "ERROR",
	SSCONFLICT:  "SSCONFLICT",
	SRCONFLICT:  "SRCONFLICT",
	RRCONFLICT:  "RRCONFLICT",
	SH_RESOLVED: "SH_RESOLVED",
	RD_RESOLVED: "RD_RESOLVED",
	NOT_USED:    "NOT_USED",
	SHIFTREDUCE: "SHIFTREDUCE",
}

func (e e_action) String() string {
	return e_action_names[e]
}
