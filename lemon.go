// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// The state vector for the entire parser generator is recorded as
// follows.  (LEMON uses no global variables and makes little use of
// static variables.  Fields in the following structure can be thought
// of as begin global variables in the program.)
type lemon struct {
	sorted            []*state  // Table of states sorted by state number
	rule              *rule     // List of all rules
	startRule         *rule     // First rule
	nstate            int       // Number of states
	nxstate           int       // nstate with tail degenerate states removed
	nrule             int       // Number of rules
	nruleWithAction   int       // Number of rules with actions
	nsymbol           int       // Number of terminal and nonterminal symbols
	nterminal         int       // Number of terminal symbols
	minShiftReduce    int       // Minimum shift-reduce action value
	errAction         int       // Error action value
	accAction         int       // Accept action value
	noAction          int       // No-op action value
	minReduce         int       // Minimum reduce action
	maxAction         int       // Maximum action value of any kind
	symbols           []*symbol // Sorted array of pointers to symbols
	errorcnt          int       // Number of errors
	errsym            *symbol   // The error symbol
	wildcard          *symbol   // Token that matches anything
	name              string    // Name of the generated parser
	arg               string    // Declaration of the 3th argument to parser
	ctx               string    // Declaration of 2nd argument to constructor
	tokentype         string    // Type of terminal symbols in the parser stack
	vartype           string    // The default type of non-terminal symbols
	start             string    // Name of the start symbol for the gram
	stacksize         string    // Size of the parser stack
	include           string    // Code to put at the start of the C file
	error             string    // Code to execute when an error is seen
	overflow          string    // Code to execute on a stack overflow
	failure           string    // Code to execute on parser failure
	accept            string    // Code to execute when the parser excepts
	extracode         string    // Code appended to the generated file
	tokendest         string    // Code to execute to destroy token data
	vardest           string    // Code for the default non-terminal destructor
	filename          string    // Name of the input file
	outname           string    // Name of the current output file
	tokenprefix       string    // A prefix added to token names in the .h file
	nconflict         int       // Number of parsing conflicts
	nactiontab        int       // Number of entries in the yy_action[] table
	nlookaheadtab     int       // Number of entries in yy_lookahead[]
	tablesize         int       // Total table size of all tables in bytes
	basisflag         bool      // Print only basis configurations
	printPreprocessed bool      // Show preprocessor output on stdout
	has_fallback      bool      // True if any %fallback is seen in the grammar
	nolinenosflag     bool      // True if #line statements should not be printed
	argv0             string    // Name of the program
}
