// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"os"
)

func ErrorMsg(filename string, lineno int, format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %s\n", filename, lineno, fmt.Sprintf(format, args...))
}
