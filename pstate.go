// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

type pstate struct {
	filename        string    // Name of the input file
	tokenlineno     int       // Linenumber at which current token starts
	errorcnt        int       // Number of errors so far
	tokenstart      []byte    // Text of current token
	gp              *lemon    // Global state vector
	state           e_state   // The state of the parser
	fallback        *symbol   // The fallback token
	tkclass         *symbol   // Token class symbol
	lhs             *symbol   // Left-hand side of current rule
	lhsalias        string    // Alias for the LHS
	nrhs            int       // Number of right-hand side symbols seen
	rhs             []*symbol // RHS symbols
	alias           []string  // Aliases for each RHS symbol (or NULL)
	prevrule        *rule     // Previous rule parsed
	declkeyword     string    // Keyword of a declaration
	declargslot     *string   // Where the declaration argument should be put. originally a pointer to char buffer (char**)
	declArgSlotBuf  []byte    // oh boy
	declArgSlotSym  *symbol   // oh boy
	insertLineMacro bool      // Add #line before declaration insert
	decllinenoslot  *int      // Where to write declaration line number
	declassoc       e_assoc   // Assign this association to decl arguments
	preccounter     int       // Assign this precedence to decl arguments
	firstrule       *rule     // Pointer to first rule in the grammar
	lastrule        *rule     // Pointer to the most recently parsed rule
	debug           bool      // mdhender
}
