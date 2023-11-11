// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"os"
)

func fprintf(fp *os.File, format string, args ...any) {
	_, _ = fmt.Fprintf(fp, format, args...)
}
