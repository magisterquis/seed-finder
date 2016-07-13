package main

/*
 * gocode.go
 * Emit go source code
 * By J. Stuart McMurray
 * Created 20170712
 * Last Modified 20170712
 */

import (
	"fmt"
	"strconv"
	"strings"
)

/* Set of already-used string names */
var varnames map[string]int
var seen map[string]struct{}

func init() {
	varnames = make(map[string]int)
	for _, s := range []string{"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct", "chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type", "continue", "for", "import", "return", "var", "bool", "byte", "complex64", "complex128", "error", "float32", "float64", "int", "int8", "int16", "int32", "int64", "rune", "string", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "true", "false", "iota", "nil", "append", "cap", "close", "complex", "copy", "delete", "imag", "len", "make", "new", "panic", "print", "println", "real", "recover"} {
		varnames[s] = 1
	}
	seen = make(map[string]struct{})
}

/* gcBoilerplate prints the necessary boilerplate to use the rest of the Go
code output. */
func gcBoilerplate() {
	fmt.Printf("%v\n", `
import "math/rand"
func seedToString(seed int64, len int) string {
	r = rand.New(rand.NewSource(i))
	r.Seed(seed)
	b := make([]byte, len)
	for i := range b {
		b[i] = rand.Intn(256)
	}
	return string(b)
}`)
}

/* gcFound emits code to be used as a string replacement */
func gcFound(v []byte, s int64) {
	/* If we've already seen this one, don't bother */
	if _, ok := seen[string(v)]; ok {
		return
	}
	seen[string(v)] = struct{}{}

	/* Variable name form of v */
	vn := strconv.Quote(string(v))
	vn = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return r
		}
		return '_'
	}, vn[1:len(vn)-1])
	/* Append a number if we need */
	if n, ok := varnames[vn]; ok {
		vn = fmt.Sprintf("%v_%v", vn, n)
		varnames[vn] = n + 1
	}
	fmt.Printf("/* %q */\n", v)
	fmt.Printf("%v = seedToString(%v, %v)\n", vn, s, len(v))
}

/* gcNotFond notes if a string's not been found */
func gcNotFound(v []byte) {
	fmt.Printf("/* %q not found */\n", v)
}
