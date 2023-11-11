// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package scanner

// Errors used by the package.

const (
	FileLimitExceeded = constError("file limit exceeded")
)

// declarations to support constant errors
type constError string

func (ce constError) Error() string {
	return string(ce)
}
