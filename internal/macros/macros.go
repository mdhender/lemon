// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"bytes"
	"fmt"
	"unicode"
)

// EvalExpression evaluates a macro expression. A macro is considered
// to be true if it is defined with -D on the command line; otherwise
// it is considered to be false. `symtab` is the table of macros that
// were defined on the command line.
func EvalExpression(input []byte, symtab map[string]string) (bool, error) {
	expression, err := parseExpression(input)
	if err != nil {
		return false, err
	}
	result, err := evalExpression(expression, symtab)
	if err != nil {
		return false, err
	}
	return result, nil
}

func evalExpression(expression []token, symtab map[string]string) (bool, error) {
	if len(expression) == 0 {
		return false, nil
	}
	var stack []bool
	for _, tok := range expression {
		switch tok.kind {
		case cTrue:
			stack = append(stack, true)
		case cFalse:
			stack = append(stack, false)
		case cVariable:
			if _, ok := symtab[string(tok.value)]; ok {
				stack = append(stack, true)
			} else {
				stack = append(stack, false)
			}
		case cAnd:
			if len(stack) < 2 {
				return false, fmt.Errorf("%d: unexpected '&&'", tok.start)
			}
			rh := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lh := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, lh && rh)
		case cOr:
			if len(stack) < 2 {
				return false, fmt.Errorf("%d: unexpected '||'", tok.start)
			}
			rh := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			lh := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, lh || rh)
		case cNot:
			if len(stack) == 0 {
				return false, fmt.Errorf("%d: unexpected '!'", tok.start)
			}
			lh := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			stack = append(stack, !lh)
		default:
			return false, fmt.Errorf("%d: unexpected %q", tok.start, tok)
		}
	}
	if len(stack) == 0 {
		// should never happen
		return false, fmt.Errorf("0: empty stack")
	} else if len(stack) != 1 {
		return false, fmt.Errorf("0: error evaluating")
	}
	return stack[0], nil
}

func parseExpression(input []byte) ([]token, error) {
	// parse the input (which is an infix boolean expression) into tokens
	var infix []token
	for pos := 0; pos < len(input); {
		if unicode.IsSpace(rune(input[pos])) {
			pos++
		} else if isalpha(input[pos]) {
			start := pos
			for pos < len(input) && (isident(input[pos])) {
				pos++
			}
			infix = append(infix, token{start: start + 1, kind: cVariable, value: input[start:pos]})
		} else if input[pos] == '(' {
			infix = append(infix, token{start: pos + 1, kind: cGroupStart})
			pos++
		} else if input[pos] == ')' {
			infix = append(infix, token{start: pos + 1, kind: cGroupEnd})
			pos++
		} else if input[pos] == '!' {
			infix = append(infix, token{start: pos + 1, kind: cNot})
			pos++
		} else if bytes.HasPrefix(input[pos:], []byte{'&', '&'}) {
			infix = append(infix, token{start: pos + 1, kind: cAnd})
			pos += 2
		} else if bytes.HasPrefix(input[pos:], []byte{'|', '|'}) {
			infix = append(infix, token{start: pos + 1, kind: cOr})
			pos += 2
		} else {
			text := input[pos:]
			if len(text) > 8 {
				text = text[:8]
			}
			return nil, fmt.Errorf("%d: unexpected input %q", pos+1, string(text))
		}
	}

	// convert infix to postfix
	var postfix, operators []token
	for _, tok := range infix {
		switch tok.kind {
		case cVariable:
			postfix = append(postfix, tok)
			if len(operators) != 0 {
				op := operators[len(operators)-1]
				if op.kind != cGroupStart {
					// pop the operator and push it onto the postfix stack
					postfix, operators = append(postfix, op), operators[:len(operators)-1]
				}
			}
		case cAnd, cOr, cNot:
			operators = append(operators, tok)
		case cGroupStart:
			operators = append(operators, tok)
		case cGroupEnd:
			var op token
			for len(operators) != 0 {
				op, operators = operators[len(operators)-1], operators[:len(operators)-1]
				if op.kind == cGroupStart {
					break
				}
				postfix = append(postfix, op)
			}
			if op.kind != cGroupStart {
				return nil, fmt.Errorf("%d: unbalanced ')'", tok.start)
			}
		}
	}
	for len(operators) != 0 {
		op := operators[len(operators)-1]
		if op.kind == cGroupStart {
			return nil, fmt.Errorf("%d: unbalanced '('", op.start)
		}
		postfix = append(postfix, op)
		operators = operators[:len(operators)-1]
	}

	return postfix, nil
}

type cell int

const (
	cNone cell = iota
	cVariable
	cAnd
	cOr
	cNot
	cGroupStart
	cGroupEnd
	cTrue
	cFalse
)

func (c cell) String() string {
	switch c {
	case cNone:
		return "none"
	case cVariable:
		return "variable"
	case cAnd:
		return "and"
	case cOr:
		return "or"
	case cNot:
		return "not"
	case cGroupStart:
		return "groupStart"
	case cGroupEnd:
		return "groupEnd"
	case cTrue:
		return "true"
	case cFalse:
		return "false"
	}
	panic(fmt.Sprintf("assert(cell != %d)", c))
}

type token struct {
	start int
	kind  cell
	value []byte
}

func (t token) String() string {
	return fmt.Sprintf("{%d %s %q}", t.start, t.kind, string(t.value))
}

func isalpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isident(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || ch == '_'
}
