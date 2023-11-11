// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros

import (
	"bytes"
	"unicode"
)

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
	line := lines[m.blockStart]
	if bytes.HasPrefix(line, []byte("%ifndef")) {
		line = line[7:]
	} else if bytes.HasPrefix(line, []byte("%ifdef")) {
		line = line[6:]
	} else if bytes.HasPrefix(line, []byte("%if")) {
		line = line[3:]
	}
	for len(line) != 0 && unicode.IsSpace(rune(line[0])) {
		line = line[1:]
	}
	// fmt.Printf("macro evaluate(%q)\n", string(line))

	m.value, err = EvalExpression(line, symtab)
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
