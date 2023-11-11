// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import "bytes"

// implements an associative array for strings.
// uses a hash table to store the values.

// create a global string internment table.
var x1a = make(map[string][]byte)

// Strsafe works like strdup, sort of.
// A copy of the data is interned in a global table.
// If the data is already in the table, it is not inserted again.
func Strsafe(data []byte) []byte {
	if data == nil {
		return nil
	}
	dst, ok := x1a[string(data)]
	if !ok {
		dst = append([]byte{}, data...)
		x1a[string(data)] = dst
	}
	return dst
}

func lemonStrlen(b []byte) int {
	n := bytes.IndexByte(b, 0)
	if n == -1 {
		n = len(b)
	}
	return n
}
