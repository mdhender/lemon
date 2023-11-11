// lemon - a parser generator
// Copyright (c) 2023 Michael D Henderson. All rights reserved.

package main

import (
	"bytes"
	"fmt"
	"github.com/mdhender/lemon/internal/macros"
	"os"
	"strings"
)

// Parse (in spite of its name) scans the entire input file.
// It reads the input and tokenizes it.
// Each token is passed to the function "parseSingleToken" which builds all
// the appropriate data structures in the global state vector "gp".
//
// symtab is a table of the macro names defined on the command line with -D.
func Parse(gp *lemon, symtab map[string]string) {
	ps := pstate{
		//debug:    true,
		gp:       gp,
		filename: gp.filename,
		state:    INITIALIZE,
	}

	/* Begin by reading the input file */
	input, err := os.ReadFile(ps.filename)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: can't open this file for reading.\n", ps.filename)
		gp.errorcnt++
		return
	} else if len(input) == 0 {
		_, _ = fmt.Fprintf(os.Stderr, "%s: can't read in all %d bytes of this file.\n", ps.filename, len(input))
		gp.errorcnt++
		return
	} else if len(input) > 100_000_000 {
		_, _ = fmt.Fprintf(os.Stderr, "%s: input file too large.\n", ps.filename)
		gp.errorcnt++
		return
	}

	// pre-process the input. this evaluates the macros to include and exclude text blocks.
	input, err = macros.PreProcess(input, symtab)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s: %v\n", ps.filename, err)
		gp.errorcnt++
		return
	} else if gp.printPreprocessed {
		_, _ = fmt.Printf("%s\n", string(input))
		return
	}

	/* Now scan the text of the input file */
	pos, lineno, startline := 0, 1, 0
	for pos < len(input) {
		if input[pos] == '\n' {
			lineno++ /* Keep track of the line number */
		}
		if isspace(input[pos]) { /* Skip all white space */
			pos++
			continue
		} else if comments := scanCPPComment(input[pos:]); len(comments) != 0 { // skip c++ style comments
			pos += len(comments)
			continue
		} else if comments := scanCComment(input[pos:]); len(comments) != 0 { // skip c style comments
			lineno = bytes.Count(comments, []byte{'\n'})
			pos += len(comments)
			continue
		}

		tokenStart := pos                                              /* Mark the beginning of the token */
		ps.tokenlineno = lineno                                        /* Line number on which token begins */
		if literal := scanStringLiteral(input[pos:]); literal != nil { /* String literals */
			lineno += bytes.Count(literal, []byte{'\n'})
			pos += len(literal)
			if len(literal) == 1 || literal[len(literal)-1] != '"' {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: string starting on this line is not terminated before the end of the file.\n", ps.filename, startline)
				ps.errorcnt++
			}
			ps.tokenstart = literal
		} else if codeBlock, err := scanCodeBlock(input[pos:]); codeBlock != nil { /* A block of C code */
			lineno += bytes.Count(codeBlock, []byte{'\n'})
			pos += len(codeBlock)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: C code starting on this line: %v.\n", ps.filename, ps.tokenlineno, err)
				ps.errorcnt++
			} else if len(codeBlock) == 1 || codeBlock[len(codeBlock)-1] != '}' {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: C code starting on this line is not terminated before the end of the file.\n", ps.filename, ps.tokenlineno)
				ps.errorcnt++
			}
			ps.tokenstart = codeBlock
		} else if isalnum(input[pos]) { /* Identifiers */
			pos += 1
			for pos < len(input) && (input[pos] == '_' || isalnum(input[pos])) {
				pos++
			}
			ps.tokenstart = input[tokenStart:pos]
		} else if bytes.HasPrefix(input[pos:], []byte{':', ':', '='}) { /* The operator "::=" */
			pos += 3
			ps.tokenstart = input[tokenStart:pos]
		} else if len(input[pos:]) > 1 && input[pos] == '/' && isalpha(input[pos+1]) {
			pos += 1
			for pos < len(input) && (input[pos] == '_' || isalnum(input[pos])) {
				pos++
			}
			ps.tokenstart = input[tokenStart:pos]
		} else if len(input[pos:]) > 1 && input[pos] == '|' && isalpha(input[pos+1]) {
			pos += 1
			for pos < len(input) && (input[pos] == '_' || isalnum(input[pos])) {
				pos++
			}
			ps.tokenstart = input[tokenStart:pos]
		} else { // all other (one character) operators */
			pos++
			ps.tokenstart = input[tokenStart:pos]
		}
		// and parse the token
		parseSingleToken(&ps)
	}
	gp.rule = ps.firstrule
	gp.errorcnt = ps.errorcnt
}

// parse a single token
func parseSingleToken(psp *pstate) {
	x := string(psp.tokenstart)
	if psp.debug {
		fmt.Printf("%s:%d: Token=[%s] state=%d\n", psp.filename, psp.tokenlineno, x, psp.state)
	}
	switch psp.state {
	case INITIALIZE:
		psp.prevrule = nil
		psp.preccounter = 0
		psp.firstrule = nil
		psp.lastrule = nil
		psp.gp.nrule = 0
		fallthrough
	case WAITING_FOR_DECL_OR_RULE:
		if x[0] == '%' {
			psp.state = WAITING_FOR_DECL_KEYWORD
		} else if isNonTerminalName(x) {
			psp.lhs = Symbol_new(x)
			psp.nrhs = 0
			psp.lhsalias = ""
			psp.state = WAITING_FOR_ARROW
		} else if x[0] == '{' {
			if psp.prevrule == nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: there is no prior rule upon which to attach the code fragment which begins on this line.\n", psp.filename, psp.tokenlineno)
				psp.errorcnt++
			} else if len(psp.prevrule.code) != 0 {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: code fragment beginning on this line is not the first to follow the previous rule.\n", psp.filename, psp.tokenlineno)
				psp.errorcnt++
			} else if x == "{NEVER-REDUCE" {
				psp.prevrule.neverReduce = true
			} else {
				psp.prevrule.line = psp.tokenlineno
				psp.prevrule.code = x[1:]
				psp.prevrule.noCode = false
			}
		} else if x[0] == '[' {
			psp.state = PRECEDENCE_MARK_1
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: token %q should be either \"%%\" or a non-terminal name.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
		}
		break
	case PRECEDENCE_MARK_1:
		if !isTerminalName(x) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: the precedence symbol must be a terminal.\n", psp.filename, psp.tokenlineno)
			psp.errorcnt++
		} else if psp.prevrule == nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: there is no prior rule to assign precedence \"[%s]\".\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
		} else if psp.prevrule.precsym != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: precedence mark on this line is not the first to follow the previous rule.\n", psp.filename, psp.tokenlineno)
			psp.errorcnt++
		} else {
			psp.prevrule.precsym = Symbol_new(x)
		}
		psp.state = PRECEDENCE_MARK_2
		break
	case PRECEDENCE_MARK_2:
		if x[0] != ']' {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: missing \"]\" on precedence mark.\n", psp.filename, psp.tokenlineno)
			psp.errorcnt++
		}
		psp.state = WAITING_FOR_DECL_OR_RULE
		break
	case WAITING_FOR_ARROW:
		if x[0] == ':' && x[1] == ':' && x[2] == '=' {
			psp.state = IN_RHS
		} else if x[0] == '(' {
			psp.state = LHS_ALIAS_1
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: expected to see a \":\" following the LHS symbol %q.\n", psp.filename, psp.tokenlineno, psp.lhs.name)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case LHS_ALIAS_1:
		if isalpha(x[0]) {
			psp.lhsalias = x
			psp.state = LHS_ALIAS_2
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %q is not a valid alias for the LHS %q.\n", psp.filename, psp.tokenlineno, x, psp.lhs.name)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case LHS_ALIAS_2:
		if x[0] == ')' {
			psp.state = LHS_ALIAS_3
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: missing \")\" following LHS alias name %q.\n", psp.filename, psp.tokenlineno, psp.lhsalias)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case LHS_ALIAS_3:
		if x[0] == ':' && x[1] == ':' && x[2] == '=' {
			psp.state = IN_RHS
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: missing \".\" following: \"%s(%s)\".\n", psp.filename, psp.tokenlineno, psp.lhs.name, psp.lhsalias)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case IN_RHS:
		if x[0] == '.' {
			rp := &rule{}
			if rp == nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: can't allocate enough memory for this rule.\n", psp.filename, psp.tokenlineno)
				psp.errorcnt++
				psp.prevrule = nil
			} else {
				rp.ruleline = psp.tokenlineno
				for i := 0; i < psp.nrhs; i++ {
					rhs, alias := psp.rhs[i], psp.alias[i]
					rp.rhs = append(rp.rhs, rhs)
					rp.rhsalias = append(rp.rhsalias, alias)
					rp.rhs[i].bContent = alias != ""
				}
				rp.lhs = psp.lhs
				rp.lhsalias = psp.lhsalias
				rp.nrhs = psp.nrhs
				rp.code = ""
				rp.noCode = true
				rp.precsym = nil
				rp.index = psp.gp.nrule
				psp.gp.nrule++
				rp.nextlhs = rp.lhs.rule
				rp.lhs.rule = rp
				rp.next = nil
				if psp.firstrule == nil {
					psp.firstrule = rp
					psp.lastrule = rp
				} else {
					psp.lastrule.next = rp
					psp.lastrule = rp
				}
				psp.prevrule = rp
			}
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if isalpha(x[0]) {
			if len(psp.rhs) >= MAXRHS {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: too many symbols on RHS of rule beginning at %q.\n", psp.filename, psp.tokenlineno, x)
				psp.errorcnt++
				psp.state = RESYNC_AFTER_RULE_ERROR
			} else {
				psp.rhs = append(psp.rhs, Symbol_new(x))
				psp.alias = append(psp.alias, "")
				psp.nrhs = len(psp.rhs)
			}
		} else if (x[0] == '|' || x[0] == '/') && len(psp.rhs) != 0 && isTerminalName(x[1:]) {
			psp.nrhs = len(psp.rhs)
			// make the top symbol a multi-terminal symbol if it isn't already.
			msp := psp.rhs[psp.nrhs-1]
			if msp.type_ != MULTITERMINAL {
				origsp := msp
				msp = &symbol{
					name:    origsp.name,
					type_:   MULTITERMINAL,
					nsubsym: 1,
					subsym:  []*symbol{origsp},
				}
				psp.rhs[psp.nrhs-1] = msp
			}
			msp.subsym = append(msp.subsym, Symbol_new(string(x[1:])))
			msp.nsubsym = len(msp.subsym)
			if isNonTerminalName(x[1:]) || isNonTerminalName(msp.subsym[0].name) {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: can't form a compound containing a non-terminal.\n", psp.filename, psp.tokenlineno)
				psp.errorcnt++
			}
		} else if x[0] == '(' && psp.nrhs > 0 {
			psp.state = RHS_ALIAS_1
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: illegal character on RHS of rule: %q.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case RHS_ALIAS_1:
		if isalpha(x[0]) {
			psp.alias[psp.nrhs-1] = x
			psp.state = RHS_ALIAS_2
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %q is not a valid alias for the RHS symbol %q\n", psp.filename, psp.tokenlineno, x, psp.rhs[psp.nrhs-1].name)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case RHS_ALIAS_2:
		if x[0] == ')' {
			psp.state = IN_RHS
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: missing \")\" following LHS alias name %q.\n", psp.filename, psp.tokenlineno, psp.lhsalias)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_RULE_ERROR
		}
		break
	case WAITING_FOR_DECL_KEYWORD:
		if isalpha(x[0]) {
			psp.declkeyword = x
			psp.declargslot = nil
			psp.decllinenoslot = nil
			psp.insertLineMacro = true
			psp.state = WAITING_FOR_DECL_ARG
			switch psp.declkeyword {
			case "name":
				psp.declargslot = &(psp.gp.name)
				psp.insertLineMacro = false
			case "include":
				psp.declargslot = &(psp.gp.include)
			case "code":
				psp.declargslot = &(psp.gp.extracode)
			case "token_destructor":
				psp.declargslot = &psp.gp.tokendest
			case "default_destructor":
				psp.declargslot = &psp.gp.vardest
			case "token_prefix":
				psp.declargslot = &psp.gp.tokenprefix
				psp.insertLineMacro = false
			case "syntax_error":
				psp.declargslot = &(psp.gp.error)
			case "parse_accept":
				psp.declargslot = &(psp.gp.accept)
			case "parse_failure":
				psp.declargslot = &(psp.gp.failure)
			case "stack_overflow":
				psp.declargslot = &(psp.gp.overflow)
			case "extra_argument":
				psp.declargslot = &(psp.gp.arg)
				psp.insertLineMacro = false
			case "extra_context":
				psp.declargslot = &(psp.gp.ctx)
				psp.insertLineMacro = false
			case "token_type":
				psp.declargslot = &(psp.gp.tokentype)
				psp.insertLineMacro = false
			case "default_type":
				psp.declargslot = &(psp.gp.vartype)
				psp.insertLineMacro = false
			case "stack_size":
				psp.declargslot = &(psp.gp.stacksize)
				psp.insertLineMacro = false
			case "start_symbol":
				psp.declargslot = &(psp.gp.start)
				psp.insertLineMacro = false
			case "left":
				psp.preccounter++
				psp.declassoc = LEFT
				psp.state = WAITING_FOR_PRECEDENCE_SYMBOL
			case "right":
				psp.preccounter++
				psp.declassoc = RIGHT
				psp.state = WAITING_FOR_PRECEDENCE_SYMBOL
			case "nonassoc":
				psp.preccounter++
				psp.declassoc = NONE
				psp.state = WAITING_FOR_PRECEDENCE_SYMBOL
			case "destructor":
				psp.state = WAITING_FOR_DESTRUCTOR_SYMBOL
			case "type":
				psp.state = WAITING_FOR_DATATYPE_SYMBOL
			case "fallback":
				psp.fallback = nil
				psp.state = WAITING_FOR_FALLBACK_ID
			case "token":
				psp.state = WAITING_FOR_TOKEN_NAME
			case "wildcard":
				psp.state = WAITING_FOR_WILDCARD_ID
			case "token_class":
				psp.state = WAITING_FOR_CLASS_ID
			default:
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: unknown declaration keyword: \"%%%s\".\n", psp.filename, psp.tokenlineno, psp.declkeyword)
				psp.errorcnt++
				psp.state = RESYNC_AFTER_DECL_ERROR
			}
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: illegal declaration keyword: %q.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		}
		break
	case WAITING_FOR_DESTRUCTOR_SYMBOL:
		if !isalpha(x[0]) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: symbol name missing after %%destructor keyword.\n", psp.filename, psp.tokenlineno)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		} else {
			sp := Symbol_new(x)
			psp.declargslot = &sp.destructor
			psp.decllinenoslot = &sp.destLineno
			psp.insertLineMacro = true
			psp.state = WAITING_FOR_DECL_ARG
		}
		break
	case WAITING_FOR_DATATYPE_SYMBOL:
		if !isalpha(x[0]) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: symbol name missing after %%type keyword.\n", psp.filename, psp.tokenlineno)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		} else {
			sp := Symbol_find(x)
			if sp != nil && sp.datatype != "" {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: symbol %%type %q already defined.\n", psp.filename, psp.tokenlineno, sp.name)
				psp.errorcnt++
				psp.state = RESYNC_AFTER_DECL_ERROR
			} else {
				if sp == nil {
					sp = Symbol_new(x)
				}
				psp.declargslot = &sp.datatype
				psp.insertLineMacro = false
				psp.state = WAITING_FOR_DECL_ARG
			}
		}
		break
	case WAITING_FOR_PRECEDENCE_SYMBOL:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if isupper(x[0]) {
			sp := Symbol_new(x)
			if sp.prec >= 0 {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: symbol %q has already be given a precedence.\n", psp.filename, psp.tokenlineno, sp.name)
				psp.errorcnt++
			} else {
				sp.prec = psp.preccounter
				sp.assoc = psp.declassoc
			}
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: can't assign a precedence to %q.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
		}
		break
	case WAITING_FOR_DECL_ARG:
		if x[0] == '{' || x[0] == '"' || isalnum(x[0]) {
			if psp.declargslot == nil {
				panic("assert(psp.declargslot != nil)")
			}
			buffer := *psp.declargslot
			addLineMacro := !psp.gp.nolinenosflag && psp.insertLineMacro && psp.tokenlineno > 1 && (psp.decllinenoslot == nil || *psp.decllinenoslot != 0)
			if addLineMacro {
				if len(buffer) > 0 && !strings.HasSuffix(buffer, "\n") {
					buffer = buffer + "\n"
				}
				buffer = buffer + fmt.Sprintf("#line %d %s\n", psp.tokenlineno, psp.filename)
			}
			buffer = buffer + x
			*psp.declargslot = buffer
			if psp.decllinenoslot != nil && *psp.decllinenoslot == 0 {
				*psp.decllinenoslot = psp.tokenlineno
			}
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: illegal argument to %%%s: %q.\n", psp.filename, psp.tokenlineno, psp.declkeyword, x)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		}
		break
	case WAITING_FOR_FALLBACK_ID:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if !isupper(x[0]) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %%fallback argument %q should be a token.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
		} else {
			sp := Symbol_new(x)
			if psp.fallback == nil {
				psp.fallback = sp
			} else if sp.fallback != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: more than one fallback assigned to token %q.\n", psp.filename, psp.tokenlineno, x)
				psp.errorcnt++
			} else {
				sp.fallback = psp.fallback
				psp.gp.has_fallback = true
			}
		}
		break
	case WAITING_FOR_TOKEN_NAME:
		// Tokens do not have to be declared before use.  But they can be
		// in order to control their assigned integer number.  The number for
		// each token is assigned when it is first seen.  So by including
		//
		//     %token ONE TWO THREE.
		//
		// early in the grammar file, that assigns small consecutive values
		// to each of the tokens ONE TWO and THREE.
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if !isupper(x[0]) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %%token argument %q should be a token.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
		} else {
			Symbol_new(x)
		}
		break
	case WAITING_FOR_WILDCARD_ID:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if !isupper(x[0]) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %%wildcard argument %q should be a token.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
		} else {
			sp := Symbol_new(x)
			if psp.gp.wildcard == nil {
				psp.gp.wildcard = sp
			} else {
				_, _ = fmt.Fprintf(os.Stderr, "%s:%d: extra wildcard to token: %q.\n", psp.filename, psp.tokenlineno, x)
				psp.errorcnt++
			}
		}
		break
	case WAITING_FOR_CLASS_ID:
		if !ISLOWER(x[0]) {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %%token_class must be followed by an identifier: %q.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		} else if Symbol_find(x) != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: symbol %q already used.\n", psp.filename, psp.tokenlineno, x)
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: symbol %q already used.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		} else {
			psp.tkclass = Symbol_new(x)
			psp.tkclass.type_ = MULTITERMINAL
			psp.state = WAITING_FOR_CLASS_TOKEN
		}
		break
	case WAITING_FOR_CLASS_TOKEN:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		} else if isupper(x[0]) {
			msp := psp.tkclass
			msp.subsym = append(msp.subsym, Symbol_new(x))
			msp.nsubsym = len(msp.subsym)
		} else if (x[0] == '|' || x[0] == '/') && isupper(x[1]) {
			msp := psp.tkclass
			msp.subsym = append(msp.subsym, Symbol_new(string(x[1:])))
			msp.nsubsym = len(msp.subsym)
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "%s:%d: %%token_class argument %q should be a token.\n", psp.filename, psp.tokenlineno, x)
			psp.errorcnt++
			psp.state = RESYNC_AFTER_DECL_ERROR
		}
		break
	case RESYNC_AFTER_RULE_ERROR, RESYNC_AFTER_DECL_ERROR:
		if x[0] == '.' {
			psp.state = WAITING_FOR_DECL_OR_RULE
		}
		if x[0] == '%' {
			psp.state = WAITING_FOR_DECL_KEYWORD
		}
		break
	}
}
