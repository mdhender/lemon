// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// Each state of the generated parser's finite state machine is encoded
// as an instance of the following structure.
type state struct {
	id                int     // globally unique identifier for this state?
	bp                *config // The basis configurations for this state
	cfp               *config // All configurations in this set
	statenum          int     // Sequential number for this state
	ap                *action // List of actions for this state
	nTknAct, nNtAct   int     // Number of actions on terminals and nonterminals
	iTknOfst, iNtOfst int     // yy_action[] offset for terminals and nonterms
	iDfltReduce       int     // Default action is to REDUCE by this rule
	pDfltReduce       *rule   // The default REDUCE rule.
	autoReduce        bool    // True if this is an auto-reduce state
}
