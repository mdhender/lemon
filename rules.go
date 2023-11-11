// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// Each production rule in the grammar is stored in the following
// structure.
type rule struct {
	lhs         *symbol   // Left-hand side of the rule
	lhsalias    string    // Alias for the LHS (NULL if none)
	lhsStart    bool      // True if left-hand side is the start symbol
	ruleline    int       // Line number for the rule
	nrhs        int       // Number of RHS symbols
	rhs         []*symbol // The RHS symbols
	rhsalias    []string  // An alias for each RHS symbol (NULL if none)
	line        int       // Line number at which code begins
	code        string    // The code executed when this rule is reduced
	codePrefix  []byte    // Setup code before code[] above
	codeSuffix  []byte    // Breakdown code after code[] above
	precsym     *symbol   // Precedence symbol for this rule
	index       int       // An index number for this rule
	iRule       int       // Rule number as used in the generated tables
	noCode      bool      // True if this rule has no associated C code
	codeEmitted bool      // True if the code has been emitted already
	canReduce   bool      // True if this rule is ever reduced
	doesReduce  bool      // Reduce actions occur after optimization
	neverReduce bool      // Reduce is theoretically possible, but prevented  by actions or other outside implementation
	nextlhs     *rule     // Next rule with the same LHS
	next        *rule     // Next rule in the global list
}
