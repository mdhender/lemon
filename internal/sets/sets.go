// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package sets

type Set struct {
	slice []bool
}

var (
	size = 0
)

func SetSize(n int) {
	size = n
}
