// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package sets

type Set struct {
	elements []bool
}

// New returns a new set with the given number of elements
func New(size int) *Set {
	return &Set{
		elements: make([]bool, size+1),
	}
}

// Add will add a new element to the set.
// Returns TRUE if the element was added, and FALSE if it was already there.
func (s *Set) Add(e int) bool {
	if !(0 <= e && e < len(s.elements)) {
		panic("assert(0 <= e && e < size)")
	}
	if s.elements[e] {
		// already set to true
		return false
	}
	s.elements[e] = true
	return true
}

// Union adds every element of s2 to set.  Return TRUE if set is updated.
func (s *Set) Union(s2 *Set) bool {
	if !(len(s.elements) == len(s2.elements)) {
		panic("assert(s1.size == s2.size)")
	}
	updatedCount := 0
	for i := 0; i < len(s.elements); i++ {
		if s2.elements[i] && !s.elements[i] {
			s.elements[i], updatedCount = true, updatedCount+1
		}
	}
	return updatedCount != 0
}
