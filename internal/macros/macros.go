// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"fmt"
)

type cell int

const (
	cNone cell = iota
	cValue
	cVariable
	cAnd
	cOr
	cNot
)

// variables is the symbol table for macros.
func parseMacroExpression(postfixExpression []byte, variables map[string]string) ([]byte, error) {
	// Create a stack to store the operands and operators
	var stack []byte
	var operators []byte

	var priorCell cell

	// Iterate over the postfix expression from left to right
	for pos, token := range postfixExpression {
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
			case cValue, cVariable:
				// okay
			case cNone, cAnd, cOr, cNot:
				return nil, fmt.Errorf("%d: unexpected operator %q", pos+1, string(token))
			}
			operators = append(operators, token)
			priorCell = cAnd
		} else if token == '|' {
			switch priorCell {
			case cValue, cVariable:
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
			case cValue, cVariable:
				return nil, fmt.Errorf("%d: unexpected operator %q", pos+1, string(token))
			}
			operators = append(operators, token)
			priorCell = cNot
		} else if 'a' <= token && token <= 'z' {
			switch priorCell {
			case cValue, cVariable:
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
			return nil, fmt.Errorf("%d: unbalanced '('", len(postfixExpression))
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
