// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"os"
)

// parse the command line and do it...
func main() {
	var basisflag bool
	var compress bool
	var mhflag bool
	var noResort bool
	var nolinenosflag bool
	var printPP bool
	var quiet bool
	var rpflag bool
	showHelp := false
	var showPrecedenceConflict bool
	var sqlFlag bool
	var statistics bool
	var version bool

	// command line flags
	options := []*s_options{
		{type_: OPT_FLAG, label: "h", pBool: &showHelp, message: "Print a usage summary and exit."},
		{type_: OPT_FLAG, label: "b", pBool: &basisflag, message: "Print only the basis in report."},
		{type_: OPT_FLAG, label: "c", pBool: &compress, message: "Don't compress the action table."},
		{type_: OPT_FSTR, label: "d", pFunc: handle_d_option, message: "Output directory.  Default '.'"},
		{type_: OPT_FSTR, label: "D", pFunc: handle_D_option, message: "Define an %ifdef macro."},
		{type_: OPT_FLAG, label: "E", pBool: &printPP, message: "Print input file after preprocessing."},
		{type_: OPT_FLAG, label: "g", pBool: &rpflag, message: "Print grammar without actions."},
		{type_: OPT_FLAG, label: "m", pBool: &mhflag, message: "Output a makeheaders compatible file."},
		{type_: OPT_FLAG, label: "l", pBool: &nolinenosflag, message: "Do not print #line statements."},
		{type_: OPT_FLAG, label: "p", pBool: &showPrecedenceConflict, message: "Show conflicts resolved by precedence rules"},
		{type_: OPT_FLAG, label: "q", pBool: &quiet, message: "(Quiet) Don't print the report file."},
		{type_: OPT_FLAG, label: "r", pBool: &noResort, message: "Do not sort or renumber states."},
		{type_: OPT_FLAG, label: "s", pBool: &statistics, message: "Print parser stats to standard output."},
		{type_: OPT_FLAG, label: "S", pBool: &sqlFlag, message: "Generate the *.sql file describing the parser tables."},
		{type_: OPT_FLAG, label: "x", pBool: &version, message: "Print the version number."},
		{type_: OPT_FSTR, label: "T", pFunc: handle_T_option, message: "Specify a template file."},
		{type_: OPT_FSTR, label: "f", message: "Ignored.  (Placeholder for '-f' compiler options.)"},
		{type_: OPT_FSTR, label: "I", message: "Ignored.  (Placeholder for '-I' compiler options.)"},
		{type_: OPT_FSTR, label: "O", message: "Ignored.  (Placeholder for '-O' compiler options.)"},
		{type_: OPT_FSTR, label: "W", message: "Ignored.  (Placeholder for '-W' compiler options.)"},
	}

	if !OptInit(os.Args, options) {
		os.Exit(1)
	} else if showHelp {
		OptPrint(os.Args[0], options)
		os.Exit(0)
	} else if version {
		fmt.Printf("Lemon version 1.0\n")
		os.Exit(0)
	}
}
