// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

// this file contains routines to scan the grammar file.

import (
	"bytes"
	"fmt"
)

// scanCComment reads C style comments and returns the comment.
func scanCComment(input []byte) (comment []byte) {
	if bytes.HasPrefix(input, []byte{'/', '*'}) {
		pos := bytes.Index(input[2:], []byte{'*', '/'})
		if pos == -1 { // comment ran past end-of-input!
			return input
		}
		pos = pos + 4 // include the leading and trailing delimiters
		return input[:pos]
	}
	return nil
}

// scanCPPComment reads C++ style comments and returns the comment.
// For C++ comments, the comment runs to the end-of-line but does not include it.
func scanCPPComment(input []byte) (comment []byte) {
	if bytes.HasPrefix(input, []byte{'/', '/'}) {
		pos := bytes.IndexByte(input, '\n')
		if pos == -1 { // comment ran to end of input
			return input
		}
		return input[:pos]
	}
	return nil
}

// scanCodeBlock scans a block of C code runs from { to }.
func scanCodeBlock(input []byte) ([]byte, error) {
	if len(input) == 0 || input[0] != '{' {
		return nil, nil
	}
	pos, level := 0, 0
	for pos < len(input) {
		if input[pos] == '{' {
			pos, level = pos+1, level+1
		} else if input[pos] == '}' {
			pos, level = pos+1, level-1
			if level == 0 {
				break
			}
		} else if input[pos] == '\'' || input[pos] == '"' { // char or string literal
			quote := input[pos]
			pos = pos + 1
			for pos < len(input) {
				if input[pos] == '\\' && pos+1 < len(input) {
					pos++
				} else if input[pos] == quote {
					break
				}
				pos++
			}
			if pos == len(input) || input[pos] != quote { // unterminated
				if quote == '"' {
					return input, fmt.Errorf("unterminated string literal")
				}
				return input, fmt.Errorf("unterminated char literal")
			}
			pos++
		} else if comments := scanCComment(input[pos:]); len(comments) != 0 { // C style comment
			pos += len(comments)
		} else if comments = scanCPPComment(input[pos:]); len(comments) != 0 { // C++ style comment
			pos += len(comments)
		} else {
			pos++
		}
	}
	if pos == len(input) || level != 0 {
		// unterminated code block
		return input, fmt.Errorf("unterminated code block")
	}
	return input[:pos], nil
}

// scanStringLiteral scans a quoted string.
// The returned string includes the open and close quotes.
// NOTE: this function DOES NOT recognize backslash as an escape character!
func scanStringLiteral(input []byte) []byte {
	if len(input) == 0 || input[0] != '"' {
		return nil
	}
	pos, quote := 1, input[0]
	for pos < len(input) && input[pos] != quote {
		pos++
	}
	if !(pos < len(input)) { // missing closing quote
		return input
	}
	return input[:pos+1] // include the terminating quote
}
