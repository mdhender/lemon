// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"github.com/mdhender/lemon/internal/macros"
	"strings"
	"unicode"
)

type macroSymbolTable map[string]string

// String implements the flag.Value interface.
func (m macroSymbolTable) String() string {
	sb := strings.Builder{}
	sb.WriteByte('[')
	for k := range m {
		if sb.Len() > 1 {
			sb.Write([]byte{',', ' '})
		}
		sb.WriteString(k)
	}
	sb.WriteByte(']')
	return sb.String()
}

// Set implements the flag.Value interface
func (m macroSymbolTable) Set(name string) error {
	m[name] = "true"
	return nil
}

type macro struct {
	kind         string // if, ifdef, or ifndef
	value        bool   // value of the macro expression
	blockStart   int    // first line of macro block
	ifChildren   []*macro
	elseStart    int // optional line of macro else
	elseChildren []*macro
	blockEnd     int // last line of macro block
}

func (m *macro) evaluate(lines [][]byte, symtab map[string]string) (err error) {
	m.value, err = macros.EvalExpression(lines[m.blockStart], symtab)
	if err != nil {
		return err
	}

	var clearIf bool
	if m.kind == "if" {
		clearIf = !m.value
	} else {
		clearIf = m.value
	}
	clearElse := !clearIf

	// recursively clear out children based on the macro's value
	if clearIf {
		// clear out the "if" block
		for line := m.blockStart; line < m.elseStart; line++ {
			lines[line] = []byte{}
		}
		// evaluate "else" block children
		for _, child := range m.elseChildren {
			if err = child.evaluate(lines, symtab); err != nil {
				return err
			}
		}
	}
	if clearElse {
		// clear out the "else" block
		for line := m.elseStart; line < m.blockEnd; line++ {
			lines[line] = []byte{}
		}
		// evaluate "if" block children
		for _, child := range m.ifChildren {
			if err = child.evaluate(lines, symtab); err != nil {
				return err
			}
		}
	}

	// clear out the macro %if... line
	lines[m.blockStart] = []byte{}
	// clear out the macro %else line
	lines[m.elseStart] = []byte{}
	// clear out the macro %end line
	lines[m.blockEnd] = []byte{}

	return nil
}

// ismacro is a helper function to detect macro definitions in the input grammar.
// a line defines a macro if it has the prefix followed by a space or end of input.
func ismacro(input []byte, macro string) bool {
	if !bytes.HasPrefix(input, []byte(macro)) {
		return false
	} else if len(input) > len(macro) && !unicode.IsSpace(rune(input[len(macro)])) {
		return false
	}
	return true
}

// mdhender: moved to internal/macros/pre_process
//// Run the preprocessor over the input file text.
//// The global variables azDefine[0] through azDefine[nDefine-1] contain the
//// names of all defined macros. This routine looks for "%ifdef" and "%ifndef"
//// and "%endif" and comments them out. Text in between is also commented out
//// as appropriate.
//func preprocess_input(input []byte, azDefine []string) ([]byte, error) {
//	// azDefine is a slice containing the names of all -D macros.
//	symtab := make(map[string]string)
//	for _, macro := range azDefine {
//		symtab[macro] = "true"
//	}
//
//	// split the input into lines
//	lines := bytes.Split(input, []byte{'\n'})
//	// trim trailing spaces from all lines
//	for i, line := range lines {
//		lines[i] = bytes.TrimRightFunc(line, func(r rune) bool {
//			return unicode.IsSpace(r)
//		})
//	}
//
//	var topLevelMacros []*macro // all top level macros
//	var activeMacro *macro      // macro being processed
//	var mstk []*macro           // stack for processing child macros
//	for i, line := range lines {
//		if ismacro(line, "%if") || ismacro(line, "%ifdef") || ismacro(line, "%ifndef") {
//			var kind string
//			if bytes.HasPrefix(line, []byte("%ifndef")) {
//				kind = "if-not"
//			} else {
//				kind = "if"
//			}
//			// we need to know if we're in an active macro before we start a new macro
//			isInActiveMacro := activeMacro != nil
//			if !isInActiveMacro {
//				// not in an active macro, so start a new one and add it to the slice of defined macros
//				activeMacro = &macro{kind: kind}
//				topLevelMacros = append(topLevelMacros, activeMacro)
//			} else {
//				// create the new one as a child of the active macro.
//				// we have to add it to the correct block (the "if" or "else" block).
//				child := &macro{kind: kind}
//				if activeMacro.elseStart == 0 {
//					// we're in the if block
//					activeMacro.ifChildren = append(activeMacro.ifChildren, child)
//				} else {
//					// we're in the else block
//					activeMacro.elseChildren = append(activeMacro.elseChildren, child)
//				}
//				// push the active macro to the stack
//				mstk = append(mstk, activeMacro)
//				// and then make this child the active macro
//				activeMacro = child
//			}
//			activeMacro.blockStart = i
//		} else if ismacro(line, "%else") {
//			if activeMacro == nil {
//				return input, fmt.Errorf("%d: \"%%else\" outside of macro\n", i+1)
//			}
//			// start the "else" block on this line. note: if we don't have an actual else
//			// block in the macro, we'll set elseStart when we find the end of the block.
//			activeMacro.elseStart = i
//		} else if ismacro(line, "%endif") {
//			if activeMacro == nil {
//				return input, fmt.Errorf("%d: \"%%endif\" outside of macro\n", i+1)
//			}
//			// this is a hack to help with clearing out the blocks later.
//			if activeMacro.elseStart == 0 {
//				// there was no else block, so set it to the end of the actual block.
//				activeMacro.elseStart = i
//			}
//			// end the macro block on this line
//			activeMacro.blockEnd = i
//			// we're done with this macro
//			if len(mstk) == 0 {
//				// we don't have a parent, so there's no longer an active macro
//				activeMacro = nil
//			} else {
//				// we do have a parent, so make it active
//				activeMacro, mstk = mstk[len(mstk)-1], mstk[:len(mstk)-1]
//			}
//		}
//	}
//	if activeMacro != nil {
//		return input, fmt.Errorf("%d: macro not terminated", activeMacro.blockStart)
//	}
//
//	// evaluate all the macros.
//	// side effect of evaluating is that it erases the contents of some lines.
//	for _, m := range topLevelMacros {
//		if err := m.evaluate(lines, symtab); err != nil {
//			return input, err
//		}
//	}
//
//	return bytes.Join(lines, []byte{'\n'}), nil
//}
