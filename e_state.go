// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

/* The state of the parser */
type e_state int

const (
	INITIALIZE e_state = iota
	WAITING_FOR_DECL_OR_RULE
	WAITING_FOR_DECL_KEYWORD
	WAITING_FOR_DECL_ARG
	WAITING_FOR_PRECEDENCE_SYMBOL
	WAITING_FOR_ARROW
	IN_RHS
	LHS_ALIAS_1
	LHS_ALIAS_2
	LHS_ALIAS_3
	RHS_ALIAS_1
	RHS_ALIAS_2
	PRECEDENCE_MARK_1
	PRECEDENCE_MARK_2
	RESYNC_AFTER_RULE_ERROR
	RESYNC_AFTER_DECL_ERROR
	WAITING_FOR_DESTRUCTOR_SYMBOL
	WAITING_FOR_DATATYPE_SYMBOL
	WAITING_FOR_FALLBACK_ID
	WAITING_FOR_WILDCARD_ID
	WAITING_FOR_CLASS_ID
	WAITING_FOR_CLASS_TOKEN
	WAITING_FOR_TOKEN_NAME
)

var e_state_names = [...]string{
	INITIALIZE:                    "INITIALIZE",
	WAITING_FOR_DECL_OR_RULE:      "WAITING_FOR_DECL_OR_RULE",
	WAITING_FOR_DECL_KEYWORD:      "WAITING_FOR_DECL_KEYWORD",
	WAITING_FOR_DECL_ARG:          "WAITING_FOR_DECL_ARG",
	WAITING_FOR_PRECEDENCE_SYMBOL: "WAITING_FOR_PRECEDENCE_SYMBOL",
	WAITING_FOR_ARROW:             "WAITING_FOR_ARROW",
	IN_RHS:                        "IN_RHS",
	LHS_ALIAS_1:                   "LHS_ALIAS_1",
	LHS_ALIAS_2:                   "LHS_ALIAS_2",
	LHS_ALIAS_3:                   "LHS_ALIAS_3",
	RHS_ALIAS_1:                   "RHS_ALIAS_1",
	RHS_ALIAS_2:                   "RHS_ALIAS_2",
	PRECEDENCE_MARK_1:             "PRECEDENCE_MARK_1",
	PRECEDENCE_MARK_2:             "PRECEDENCE_MARK_2",
	RESYNC_AFTER_RULE_ERROR:       "RESYNC_AFTER_RULE_ERROR",
	RESYNC_AFTER_DECL_ERROR:       "RESYNC_AFTER_DECL_ERROR",
	WAITING_FOR_DESTRUCTOR_SYMBOL: "WAITING_FOR_DESTRUCTOR_SYMBOL",
	WAITING_FOR_DATATYPE_SYMBOL:   "WAITING_FOR_DATATYPE_SYMBOL",
	WAITING_FOR_FALLBACK_ID:       "WAITING_FOR_FALLBACK_ID",
	WAITING_FOR_WILDCARD_ID:       "WAITING_FOR_WILDCARD_ID",
	WAITING_FOR_CLASS_ID:          "WAITING_FOR_CLASS_ID",
	WAITING_FOR_CLASS_TOKEN:       "WAITING_FOR_CLASS_TOKEN",
	WAITING_FOR_TOKEN_NAME:        "WAITING_FOR_TOKEN_NAME",
}

func (e e_state) String() string {
	return e_state_names[e]
}
