// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package macros_test

import (
	"fmt"
	"github.com/mdhender/lemon/internal/macros"
	"testing"
)

func TestParseMacroExpression(t *testing.T) {
	variables := make(map[string]string)
	variables["a"] = "true"
	variables["b"] = "true"

	type test_case struct {
		id     int
		input  string
		expect bool
		err    error
	}
	for _, tc := range []test_case{
		{id: 1, input: "a", expect: true},
		{id: 2, input: "b", expect: true},
		{id: 3, input: "c", expect: false},
		{id: 4, input: "a && b", expect: true},
		{id: 5, input: "a && c", expect: false},
		{id: 6, input: "!a", expect: false},
		{id: 7, input: "!a || !c", expect: true},
		{id: 8, input: "( a ) && ! ( b )", expect: false},
		{id: 9, input: "a || b || (((c)))", expect: true},
		{id: 10, input: "!(c)", expect: true},
		{id: 11, input: "!(!c)", expect: false},
		{id: 12, input: "!(!(c))", expect: false},
		{id: 13, input: ")", err: fmt.Errorf("1: unbalanced ')'")},
	} {
		//fmt.Printf("%2d: %q\n", tc.id, tc.input)
		got, err := macros.EvalExpression([]byte(tc.input), variables)
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
		} else if tc.expect != got {
			t.Errorf("%2d: want %v: got %v\n", tc.id, tc.expect, got)
		}
	}
}
