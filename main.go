// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"flag"
	"fmt"
	"github.com/mdhender/lemon/internal/sets"
	"os"
)

// parse the command line and do it...
func main() {
	var compress bool
	var mhflag bool
	var noResort bool
	var quiet bool
	rpflag := false
	var showPrecedenceConflict bool
	var sqlFlag bool
	var statistics bool
	var version bool

	var macdefs macroSymbolTable = make(map[string]string)

	lem := &lemon{
		argv0: os.Args[0],
	}

	flag.BoolVar(&lem.basisflag, "b", lem.basisflag, "Print only the basis in report.")
	flag.BoolVar(&lem.nolinenosflag, "l", lem.nolinenosflag, "Do not print #line statements.")
	flag.BoolVar(&lem.printPreprocessed, "E", lem.printPreprocessed, "Print input file after preprocessing.")

	flag.BoolVar(&compress, "c", compress, "Don't compress the action table.")
	flag.BoolVar(&rpflag, "g", rpflag, "Print grammar without actions.")
	flag.BoolVar(&mhflag, "m", mhflag, "Output a makeheaders compatible file.")
	flag.BoolVar(&showPrecedenceConflict, "p", showPrecedenceConflict, "Show conflicts resolved by precedence rules")
	flag.BoolVar(&quiet, "q", quiet, "(Quiet) Don't print the report file.")
	flag.BoolVar(&noResort, "r", noResort, "Do not sort or renumber states.")
	flag.BoolVar(&statistics, "s", statistics, "Print parser stats to standard output.")
	flag.BoolVar(&sqlFlag, "S", sqlFlag, "Generate the *.sql file describing the parser tables.")
	flag.BoolVar(&version, "x", version, "Print the version number.")
	flag.StringVar(&outputDir, "d", outputDir, "Output directory.")
	flag.StringVar(&lem.filename, "i", lem.filename, "Grammar file to process.")
	flag.StringVar(&user_templatename, "T", user_templatename, "Specify a template file.")
	flag.Var(macdefs, "D", "Define macro.")
	//{type_: OPT_FSTR, label: "f", message: "Ignored.  (Placeholder for '-f' compiler options.)"},
	//{type_: OPT_FSTR, label: "I", message: "Ignored.  (Placeholder for '-I' compiler options.)"},
	//{type_: OPT_FSTR, label: "O", message: "Ignored.  (Placeholder for '-O' compiler options.)"},
	//{type_: OPT_FSTR, label: "W", message: "Ignored.  (Placeholder for '-W' compiler options.)"},
	flag.Parse()
	if version {
		fmt.Printf("Lemon version 1.0\n")
		os.Exit(0)
	}
	argsLeftOver := flag.Args()
	if lem.filename == "" && len(argsLeftOver) != 0 {
		lem.filename, argsLeftOver = argsLeftOver[0], argsLeftOver[1:]
	}
	if len(argsLeftOver) != 0 {
		for _, arg := range argsLeftOver {
			_, _ = fmt.Fprintf(os.Stderr, "error: unknown option %q.\n", arg)
		}
		os.Exit(1)
	}
	if lem.filename == "" {
		_, _ = fmt.Fprintf(os.Stderr, "error: missing grammar file name on command line.\n")
		os.Exit(1)
	}

	// initialize the machine
	// Strsafe_init() // mdhender - no longer needed
	// Symbol_init()// mdhender - no longer needed
	// State_init()// mdhender - no longer needed
	Symbol_new("$")

	// parse the input file
	Parse(lem, macdefs)
	if lem.errorcnt != 0 {
		_, _ = fmt.Fprintf(os.Stderr, "error: parse failed with %d errors.\n", lem.errorcnt)
		os.Exit(1)
	} else if lem.nrule == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "error: grammar file contains no rules.\n")
		os.Exit(1)
	}
	lem.errsym = Symbol_find("error")

	// count and index the symbols of the grammar
	Symbol_new("{default}")
	lem.nsymbol = Symbol_count()
	lem.symbols = Symbol_sortedSlice()

	i := lem.nsymbol
	for i > 1 && lem.symbols[i-1].type_ == MULTITERMINAL {
		i--
	}
	if lem.symbols[i-1].name != "{default}" {
		_, _ = fmt.Fprintf(os.Stderr, "error: internal error: want symbols[%d] to be %q, got %q\n", i-1, "{default}", lem.symbols[i-1].name)
		os.Exit(1)
	}
	lem.nsymbol = i - 1
	// count the number of terminal symbols and update the index
	for i = 1; isupper(lem.symbols[i].name[0]); i++ {
		//
	}
	lem.nterminal = i

	// Assign sequential rule numbers.  Start with 0.  Put rules that have no
	// reduce action C-code associated with them last, so that the switch()
	// statement that selects reduction actions will have a smaller jump table.
	rulesIndex := 0
	for rp := lem.rule; rp != nil; rp = rp.next {
		if rp.code != "" {
			rp.iRule = rulesIndex
			rulesIndex = rulesIndex + 1
		} else {
			// negative rule index means no reduce action
			rp.iRule = -1
		}
	}
	lem.nruleWithAction = rulesIndex
	// now update the index for the rules with no reduce action,
	// putting them at the end of the list by resetting their index
	for rp := lem.rule; rp != nil; rp = rp.next {
		if rp.iRule < 0 {
			rp.iRule = rulesIndex
			rulesIndex = rulesIndex + 1
		}
	}
	lem.startRule = lem.rule
	lem.rule = lem.rule.sort()

	/* Generate a reprint of the grammar, if requested on the command line */
	if rpflag {
		Reprint(os.Stdout, lem)
	} else {
		/* Initialize the size for all follow and first sets */
		sets.SetSize(lem.nterminal + 1)

		/* Find the precedence for every production rule (that has one) */
		FindRulePrecedences(lem.rule)

		///* Compute the lambda-nonterminals and the first-sets for every
		// ** nonterminal */
		//FindFirstSets(&lem);
		//
		///* Compute all LR(0) states.  Also record follow-set propagation
		// ** links so that the follow-set can be computed later */
		//lem.nstate = 0;
		//FindStates(&lem);
		//lem.sorted = State_arrayof();
		//
		///* Tie up loose ends on the propagation links */
		//FindLinks(&lem);
		//
		///* Compute the follow set of every reducible configuration */
		//FindFollowSets(&lem);
		//
		///* Compute the action tables */
		//FindActions(&lem);
		//
		///* Compress the action tables */
		//if (compress == 0) {
		//	CompressTables(&lem);
		//}
		//
		///* Reorder and renumber the states so that states with fewer choices
		// ** occur at the end.  This is an optimization that helps make the
		// ** generated parser tables smaller. */
		//if (noResort == 0) {
		//	ResortStates(&lem);
		//}
		//
		///* Generate a report of the parser generated.  (the "y.output" file) */
		//if (!quiet) {
		//	ReportOutput(&lem);
		//}
		//
		///* Generate the source code for the parser */
		//ReportTable(&lem, mhflag, sqlFlag);
		//
		///* Produce a header file for use by the scanner.  (This step is
		// ** omitted if the "-m" option is used because makeheaders will
		// ** generate the file for us.) */
		//if (!mhflag) {
		//	ReportHeader(&lem);
		//}
	}
}
