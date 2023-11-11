// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// Every shift or reduce operation is stored as one of the following
type action struct {
	sp    *symbol // The look-ahead symbol
	type_ e_action
	x     struct { // union
		stp *state // The new state, if a shift
		rp  *rule  // The rule, if a reduce
	}
	spOpt   *symbol // SHIFTREDUCE optimization to this symbol
	next    *action // Next action for this state
	collide *action // Next action with the same hash

	seq int // mdhender: added to support sorting actions
}
