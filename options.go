// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// global variables that are updated via command line switches
var (
	// azDefine is the symbol table for macro definitions
	azDefine map[string]string
	nDefine  int
	// outputDir is the path to create output files in
	outputDir = "."
	// user_templatename is the name of the template file to load
	user_templatename string
)

type option_type int

const (
	OPT_FLAG option_type = iota + 1
	OPT_FSTR
)

type s_options struct {
	type_   option_type
	label   string
	pBool   *bool
	pFunc   func(string) bool
	message string
}

func OptInit(args []string, options []*s_options) bool {
	programName := args[0]
	errorsFound := 0
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "+") {
			if !handleflags(arg, options) {
				errorsFound++
			}
		} else if strings.HasPrefix(arg, "-") {
			if !handleflags(arg, options) {
				errorsFound++
			}
		} else if k, v, ok := strings.Cut(arg, "="); ok {
			if !handleswitch(k, v, options) {
				errorsFound++
			}
		}
	}
	if errorsFound != 0 {
		OptPrint(programName, options)
	}
	return errorsFound == 0
}

func OptPrint(programName string, options []*s_options) {
	_, _ = fmt.Fprintf(os.Stderr, "Valid command line options for %q are:\n", programName)
	for _, opt := range options {
		switch opt.type_ {
		case OPT_FLAG:
			_, _ = fmt.Fprintf(os.Stderr, "  -%s         %s\n", opt.label, opt.message)
		case OPT_FSTR:
			_, _ = fmt.Fprintf(os.Stderr, "  -%s<string> %s\n", opt.label, opt.message)
		}
	}
}

// process a flag command line argument.
// returns true if there were no error with the argument.
func handleflags(arg string, options []*s_options) bool {
	var opt *s_options
	for j := 0; j < len(options) && opt == nil; j++ {
		if arg[1:] == options[j].label {
			opt = options[j]
			break
		}
	}

	if opt == nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: undefined option %q\n", arg)
		return false
	} else if opt.pBool == nil && opt.pFunc == nil {
		// ignore this option
		return true
	}

	switch opt.type_ {
	case OPT_FLAG: // boolean flag
		if opt.pBool == nil {
			panic("assert(opt.pBool != nil)")
		}
		*opt.pBool = arg[0] == '-'
		fmt.Printf("set %q to %v\n", arg, *opt.pBool)
		return true
	case OPT_FSTR: // string function flag
		if opt.pFunc == nil {
			panic("assert(opt.pFunc != nil)")
		}
		return opt.pFunc(arg[2:])
	}
	_, _ = fmt.Fprintf(os.Stderr, "error: command line syntax error: missing argument on switch %q.\n", arg)
	return false
}

/*
** Process a command line switch which has an argument.
 */
func handleswitch(arg, val string, options []*s_options) bool {
	var opt *s_options
	for j := 0; j < len(options) && opt == nil; j++ {
		if arg[1:] == options[j].label {
			opt = options[j]
		}
	}

	if opt == nil {
		_, _ = fmt.Fprintf(os.Stderr, "error: undefined option %q\n", arg)
		return false
	} else if opt.pBool == nil && opt.pFunc == nil {
		// ignore this option
		return true
	}

	switch opt.type_ {
	case OPT_FLAG:
		_, _ = fmt.Fprintf(os.Stderr, "error: option %q requires an argument.\n", arg)
		return false
	case OPT_FSTR:
		if opt.pFunc == nil {
			panic("assert(opt.pFunc != nil)")
		}
		return opt.pFunc(val)
	}
	return false
}

// add the macro defined to the azDefine array.
func handle_D_option(macroName string) bool {
	azDefine[macroName] = "true"
	nDefine = len(azDefine)
	return true
}

// update the path to the output directory
func handle_d_option(path string) bool {
	outputDir = filepath.Clean(path)
	return true
}

// update the path to the user template file
func handle_T_option(path string) bool {
	user_templatename = filepath.Clean(path)
	return true
}
