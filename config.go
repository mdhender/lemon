// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"fmt"
	"github.com/mdhender/lemon/internal/sets"
)

type config struct {
	rp     *rule     // The rule upon which the configuration is based
	dot    int       // The parse point
	fws    *sets.Set // Follow-set for this configuration o
	fplp   *plink    // Follow-set forward propagation links
	bplp   *plink    // Follow-set backwards propagation links
	stp    *state    // Pointer to state which contains this
	status cfgstatus // used during followset and shift computations
	next   *config   // Next configuration in the state
	bp     *config   // The next basis configuration
}

func (c *config) String() string {
	return fmt.Sprintf("(config (hash %v) (dot %d))", confighash(c), c.dot)
}
