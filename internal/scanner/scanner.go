// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package scanner

// this file contains routines to scan the grammar file.

import (
	"bytes"
	"fmt"
	"os"
	"unicode"
)

//// mdhender
//// todo: left off moving ~/parser.go/parseSingleToken into this package.
//// todo: am thinking it should not be moved at all.
//func Scan(filename string, input []byte) error {
//	var ps struct {
//		gp          any
//		errorcnt    int
//		filename    string
//		state       e_state
//		tokenlineno int
//		tokenstart  []byte
//	}
//	ps.state = INITIALIZE
//
//	pos, lineno, startline := 0, 1, 0
//	for pos < len(input) {
//		if input[pos] == '\n' {
//			lineno++ /* Keep track of the line number */
//		}
//		if isspace(input[pos]) { /* Skip all white space */
//			pos++
//			continue
//		} else if comments := scanCPPComment(input[pos:]); len(comments) != 0 { // skip c++ style comments
//			pos += len(comments)
//			continue
//		} else if comments := scanCComment(input[pos:]); len(comments) != 0 { // skip c style comments
//			lineno = bytes.Count(comments, []byte{'\n'})
//			pos += len(comments)
//			continue
//		}
//
//		tokenStart := pos                                              /* Mark the beginning of the token */
//		ps.tokenlineno = lineno                                        /* Line number on which token begins */
//		if literal := scanStringLiteral(input[pos:]); literal != nil { /* String literals */
//			lineno += bytes.Count(literal, []byte{'\n'})
//			pos += len(literal)
//			if len(literal) == 1 || literal[len(literal)-1] != '"' {
//				fprintf(os.Stderr, "%s:%d: string starting on this line is not terminated before the end of the file.", ps.filename, startline)
//				ps.errorcnt++
//			}
//			ps.tokenstart = literal
//		} else if codeBlock := scanCodeBlock(input[pos:]); codeBlock != nil { /* A block of C code */
//			lineno += bytes.Count(codeBlock, []byte{'\n'})
//			pos += len(codeBlock)
//			if len(codeBlock) == 1 || codeBlock[len(codeBlock)-1] != '}' {
//				fprintf(os.Stderr, "%s:%d: C code starting on this line is not terminated before the end of the file.", ps.filename, ps.tokenlineno)
//				ps.errorcnt++
//			}
//			ps.tokenstart = codeBlock
//		} else if isalnum(input[pos]) { /* Identifiers */
//			pos += 1
//			for pos < len(input) && (input[pos] == '_' || isalnum(input[pos])) {
//				pos++
//			}
//			ps.tokenstart = input[tokenStart:pos]
//		} else if bytes.HasPrefix(input[pos:], []byte{':', ':', '='}) { /* The operator "::=" */
//			pos += 3
//			ps.tokenstart = input[tokenStart:pos]
//		} else if len(input[pos:]) > 1 && input[pos] == '/' && isalpha(input[pos+1]) {
//			pos += 1
//			for pos < len(input) && (input[pos] == '_' || isalnum(input[pos])) {
//				pos++
//			}
//			ps.tokenstart = input[tokenStart:pos]
//		} else if len(input[pos:]) > 1 && input[pos] == '|' && isalpha(input[pos+1]) {
//			pos += 1
//			for pos < len(input) && (input[pos] == '_' || isalnum(input[pos])) {
//				pos++
//			}
//			ps.tokenstart = input[tokenStart:pos]
//		} else { // all other (one character) operators */
//			pos++
//			ps.tokenstart = input[tokenStart:pos]
//		}
//		// and parse the token
//		parseSingleToken(&ps)
//	}
//	gp.rule = ps.firstrule
//	gp.errorcnt = ps.errorcnt
//
//	return nil
//}

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
func scanCodeBlock(input []byte) []byte {
	if len(input) == 0 || input[0] != '{' {
		return nil
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
		} else if input[pos] == '\'' { // char literal
			for pos = pos + 1; pos < len(input) && !(input[pos] == '\'' && input[pos-1] != '\\'); pos++ {
				//
			}
		} else if input[pos] == '"' { // string literal
			for pos = pos + 1; pos < len(input) && !(input[pos] == '"' && input[pos-1] != '\\'); pos++ {
				//
			}
		} else if comments := scanCComment(input[pos:]); len(comments) != 0 { // C style comment
			pos += len(comments)
		} else if comments = scanCPPComment(input[pos:]); len(comments) != 0 { // C++ style comment
			pos += len(comments)
		}
	}
	return input[:pos]
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

// fprintf is a helper to print to stderr.
func fprintf(fp *os.File, format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, format, args...)
}

func isalnum(x byte) bool {
	return unicode.IsLetter(rune(x)) || unicode.IsDigit(rune(x))
}

func isalpha(x byte) bool {
	return unicode.IsLetter(rune(x))
}

func isspace(x byte) bool {
	return unicode.IsSpace(rune(x))
}
