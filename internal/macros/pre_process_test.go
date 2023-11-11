// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros_test

import (
	"github.com/mdhender/lemon/internal/macros"
	"testing"
)

func TestPreProcess(t *testing.T) {
	symtab := make(map[string]string)
	symtab["a"] = "true"
	symtab["b"] = "true"

	type test_case struct {
		id     int
		input  string
		expect string
		err    error
	}

	for _, tc := range []test_case{
		{id: 1,
			input:  "bof\n%if a\na is true\n%else\na is false\n%endif\neof",
			expect: "bof\n\na is true\n\n\n\neof",
		},
		{id: 2,
			input:  "bof\n%if !b\n!b is true\n%else\n!b is false\n%endif\neof",
			expect: "bof\n\n\n\n!b is false\n\neof",
		},
		{id: 3,
			input:  "bof\n%ifdef c\nc is defined\n%else\nc is not defined\n%endif\neof",
			expect: "bof\n\n\n\nc is not defined\n\neof",
		},
		{id: 4,
			input:  "bof\n%ifndef d\nd is not defined\n%else\nd is defined\n%endif\neof",
			expect: "bof\n\nd is not defined\n\n\n\neof",
		},
	} {
		got, err := macros.PreProcess([]byte(tc.input), symtab)
		if tc.err != nil {
			if err != nil {
				want, got := tc.err.Error(), err.Error()
				if want != got {
					t.Errorf("%2d: want %v: got %+v\n", tc.id, tc.err, err)
				}
			} else {
				t.Errorf("%2d: want %v: got success\n", tc.id, tc.err)
			}
		} else if err != nil {
			t.Errorf("%2d: want success: got %+v\n", tc.id, err)
		} else if tc.expect != string(got) {
			t.Errorf("%2d: want\nvvvvv\n%s\n^^^^\n%2d: got\nvvvv\n%s\n^^^^\n", tc.id, tc.expect, tc.id, got)
		}
	}
}
