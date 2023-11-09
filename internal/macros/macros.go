// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"bytes"
	"fmt"
	"unicode"
)

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

type expr struct {
	kind     string
	value    []byte
	variable []byte
	and, or  struct {
		lh, rh *expr
	}
	not   *expr
	group *expr
}

func (e *expr) String() string {
	if e == nil {
		return "nil"
	} else if e.value != nil {
		return fmt.Sprintf("value(%q)", string(e.value))
	} else if e.variable != nil {
		return fmt.Sprintf("variable(%q)", string(e.variable))
	} else {
		return "expr{...}"
	}
}

func isalpha(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

func isident(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || ('0' <= ch && ch <= '9') || ch == '_'
}

func evalMacro(input []byte, symtab map[string]string) (bool, error) {

	return false, fmt.Errorf("not implemented")
}

func parseExpression(input []byte) ([]token, error) {
	var postfix, infix, operators []token

	// parse the input (which is a postfix boolean expression) into tokens
	for pos := 0; pos < len(input); {
		if unicode.IsSpace(rune(input[pos])) {
			pos++
		} else if isalpha(input[pos]) {
			start := pos
			for pos < len(input) && (isident(input[pos])) {
				pos++
			}
			postfix = append(postfix, token{start: start + 1, kind: cVariable, value: input[start:pos]})
		} else if input[pos] == '(' {
			postfix = append(postfix, token{start: pos + 1, kind: cGroupStart})
			pos++
		} else if input[pos] == ')' {
			postfix = append(postfix, token{start: pos + 1, kind: cGroupEnd})
			pos++
		} else if input[pos] == '!' {
			postfix = append(postfix, token{start: pos + 1, kind: cNot})
			pos++
		} else if bytes.HasPrefix(input[pos:], []byte{'&', '&'}) {
			postfix = append(postfix, token{start: pos + 1, kind: cAnd})
			pos += 2
		} else if bytes.HasPrefix(input[pos:], []byte{'|', '|'}) {
			postfix = append(postfix, token{start: pos + 1, kind: cOr})
			pos += 2
		} else {
			text := input[pos:]
			if len(text) > 8 {
				text = text[:8]
			}
			return nil, fmt.Errorf("%d: unexpected input %q", pos+1, string(text))
		}
	}

	// convert the postfix expression into infix expression
	for _, tok := range postfix {
		switch tok.kind {
		case cAnd:
			operators = append(operators, tok)
		case cGroupEnd:
			if len(operators) == 0 {
				return nil, fmt.Errorf("%d: unexpected ')'", tok.start)
			}
			op := operators[len(operators)-1]
			if op.kind == cGroupStart {
				operators = operators[:len(operators)-1]
			} else {
				fmt.Printf("operators")
				for _, op := range operators {
					switch op.kind {
					case cAnd:
						fmt.Print(" &&")
					case cGroupEnd:
						fmt.Print(" ) ")
					case cGroupStart:
						fmt.Print(" ( ")
					case cNot:
						fmt.Print(" ! ")
					case cOr:
						fmt.Print(" ||")
					default:
						fmt.Print(op)
					}
				}
				fmt.Println("")
				return nil, fmt.Errorf("%d: unbalanced ')'", tok.start)
			}
		case cGroupStart:
			operators = append(operators, tok)
		case cNot:
			operators = append(operators, tok)
		case cOr:
			operators = append(operators, tok)
		case cVariable:
			infix = append(infix, tok)
			if len(operators) != 0 {
				op := operators[len(operators)-1]
				if op.kind != cGroupStart {
					infix, operators = append(infix, op), operators[:len(operators)-1]
				}
			}
		}
	}
	for len(operators) != 0 {
		op := operators[len(operators)-1]
		if op.kind == cGroupStart {
			return nil, fmt.Errorf("%d: unbalanced '('", op.start)
		}
		infix, operators = append(infix, op), operators[:len(operators)-1]
	}

	return infix, nil
}

// variables is the symbol table for macros.
func parseMacroExpression(input []byte, variables map[string]string) ([]byte, error) {
	// Create a stack to store the operands and operators
	var stack []byte
	var operators []byte

	var priorCell cell

	// Iterate over the postfix expression from left to right
	for pos, token := range input {
		// If the token is an operand, push it onto the stack
		if token == '(' {
			operators = append(operators, token)
		} else if token == ')' {
			if len(operators) == 0 {
				return nil, fmt.Errorf("%d: unexpected ')'", pos+1)
			} else if operators[len(operators)-1] != '(' {
				fmt.Printf("error: unbalanced ')': operators are %q\n", string(operators))
				return nil, fmt.Errorf("%d: unbalanced ')'", pos+1)
			}
			operators = operators[:len(operators)-1]
		} else if token == '&' {
			switch priorCell {
			case cTrue, cFalse, cVariable:
				// okay
			case cNone, cAnd, cOr, cNot:
				return nil, fmt.Errorf("%d: unexpected operator %q", pos+1, string(token))
			}
			operators = append(operators, token)
			priorCell = cAnd
		} else if token == '|' {
			switch priorCell {
			case cTrue, cFalse, cVariable:
				// okay
			case cNone, cAnd, cOr, cNot:
				return nil, fmt.Errorf("%d: unexpected operator %q", pos+1, string(token))
			}
			operators = append(operators, token)
			priorCell = cOr
		} else if token == '!' {
			switch priorCell {
			case cNone, cAnd, cOr, cNot:
				// okay
			case cTrue, cFalse, cVariable:
				return nil, fmt.Errorf("%d: unexpected operator %q", pos+1, string(token))
			}
			operators = append(operators, token)
			priorCell = cNot
		} else if 'a' <= token && token <= 'z' {
			switch priorCell {
			case cTrue, cFalse, cVariable:
				return nil, fmt.Errorf("%d: unexpected variable %q", pos+1, string(token))
			case cNone, cAnd, cOr, cNot:
				// okay
			}
			var variable byte
			if _, ok := variables[string(token)]; ok {
				//variable = 'T'
				variable = token
			} else {
				//variable = 'F'
				variable = token
			}
			stack = append(stack, variable)
			if len(operators) != 0 && operators[len(operators)-1] != '(' {
				stack, operators = append(stack, operators[len(operators)-1]), operators[:len(operators)-1]
			}
			priorCell = cVariable
		} else {
			return nil, fmt.Errorf("%d: unknown token %q", pos+1, string(token))
		}
	}

	for len(operators) != 0 {
		if operators[len(operators)-1] == '(' {
			return nil, fmt.Errorf("%d: unbalanced '('", len(input))
		}
		stack, operators = append(stack, operators[len(operators)-1]), operators[:len(operators)-1]
	}

	return stack, nil
}

//func isOperand(token byte) bool {
//	return token == 't' || token == 'f'
//}

//func calculate(operand1, operand2, operator string) string {
//	var result string
//	switch operator {
//	case "+":
//		result = operand1 + operand2
//	case "-":
//		result = operand1 - operand2
//	case "*":
//		result = operand1 * operand2
//	case "/":
//		result = operand1 / operand2
//	}
//	return result
//}
