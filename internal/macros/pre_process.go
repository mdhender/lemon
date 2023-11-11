// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"bytes"
	"fmt"
	"unicode"
)

// PreProcess runs the input through the macro preprocessor.
// It returns the processed text. As a side effect of the processing,
// trailing spaces are stripped from all lines of the input.
//
// This routine looks for "%if," "%ifdef," and "%ifndef" macros in the input.
// These blocks end with "%endif" and may have an optional "%else" section.
//
// The macros are evaluated using the symbol table passed in. If the macro
// name is defined in the table, its value is true. If not, its value
// is false.
//
// If the macro evaluates to true, text in the "if" block is kept and text
// in the "else" block is removed. If the macro evaluates to false, then
// the opposite happens.
//
// The processed text will have the same number of lines; text is removed
// by deleting the contents of each line in the appropriate block.
func PreProcess(input []byte, symtab map[string]string) ([]byte, error) {
	// split the input into lines
	lines := bytes.Split(input, []byte{'\n'})

	// trim trailing spaces from all lines
	for i, line := range lines {
		lines[i] = bytes.TrimRightFunc(line, func(r rune) bool {
			return unicode.IsSpace(r)
		})
	}

	var topLevelMacros []*macro // all top level macros
	var activeMacro *macro      // macro being processed
	var mstk []*macro           // stack for processing child macros
	for i, line := range lines {
		if ismacro(line, "%if") || ismacro(line, "%ifdef") || ismacro(line, "%ifndef") {
			var kind string
			if bytes.HasPrefix(line, []byte("%ifndef")) {
				kind = "if-not"
			} else {
				kind = "if"
			}
			// we need to know if we're in an active macro before we start a new macro
			isInActiveMacro := activeMacro != nil
			if !isInActiveMacro {
				// not in an active macro, so start a new one and add it to the slice of defined macros
				activeMacro = &macro{kind: kind}
				topLevelMacros = append(topLevelMacros, activeMacro)
			} else {
				// create the new one as a child of the active macro.
				// we have to add it to the correct block (the "if" or "else" block).
				child := &macro{kind: kind}
				if activeMacro.elseStart == 0 {
					// we're in the if block
					activeMacro.ifChildren = append(activeMacro.ifChildren, child)
				} else {
					// we're in the else block
					activeMacro.elseChildren = append(activeMacro.elseChildren, child)
				}
				// push the active macro to the stack
				mstk = append(mstk, activeMacro)
				// and then make this child the active macro
				activeMacro = child
			}
			activeMacro.blockStart = i
		} else if ismacro(line, "%else") {
			if activeMacro == nil {
				return input, fmt.Errorf("%d: \"%%else\" outside of macro\n", i+1)
			}
			// start the "else" block on this line. note: if we don't have an actual else
			// block in the macro, we'll set elseStart when we find the end of the block.
			activeMacro.elseStart = i
		} else if ismacro(line, "%endif") {
			if activeMacro == nil {
				return input, fmt.Errorf("%d: \"%%endif\" outside of macro\n", i+1)
			}
			// this is a hack to help with clearing out the blocks later.
			if activeMacro.elseStart == 0 {
				// there was no else block, so set it to the end of the actual block.
				activeMacro.elseStart = i
			}
			// end the macro block on this line
			activeMacro.blockEnd = i
			// we're done with this macro
			if len(mstk) == 0 {
				// we don't have a parent, so there's no longer an active macro
				activeMacro = nil
			} else {
				// we do have a parent, so make it active
				activeMacro, mstk = mstk[len(mstk)-1], mstk[:len(mstk)-1]
			}
		}
	}
	if activeMacro != nil {
		return input, fmt.Errorf("%d: macro not terminated", activeMacro.blockStart)
	}

	// evaluate all the macros.
	// side effect of evaluating is that it erases the contents of some lines.
	for _, m := range topLevelMacros {
		if err := m.evaluate(lines, symtab); err != nil {
			return input, err
		}
	}

	// join the lines back together and return them
	return bytes.Join(lines, []byte{'\n'}), nil
}
