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
	"log"
	"os"
	"strconv"
	"strings"
	"time"
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
code output.  p gives the package name to use. */
func gcBoilerplate(o *os.File, p string) {
	fmt.Fprintf(o, "package %v\n%v\n", p, `import "math/rand"
func seedToString(seed int64, l int) string {
	r := rand.New(rand.NewSource(seed))
	b := make([]byte, l)
	for i := range b {
		b[i] = byte(r.Intn(256))
	}
	return string(b)
}`)
}

/* gcVar tries to make a variable for the input slice.  Slices longer than
sublen bytes will be chopped into sublen-byte chunks and reassembled in the
output code.  Output will be written to o. */
func gcVar(v []byte, sublen uint, nParallel uint, o *os.File) {
	/* If we've already seen this one, don't bother */
	if _, ok := seen[string(v)]; ok {
		return
	}
	seen[string(v)] = struct{}{}

	/* Split into chunks */
	var ss []int64
	var ls []int
	var start, end int
	started := time.Now()
	for start = 0; start < len(v); start += int(sublen) {
		/* Work out subslice to find */
		end = start + int(sublen)
		if end > len(v) {
			end = len(v)
		}
		/* Get the seed for the subslice */
		log.Printf(
			"Working on bytes %v - %v of %q: %q",
			start, end-1,
			v,
			v[start:end],
		)
		s, err := findSeed(v[start:end], nParallel)
		if nil != err {
			fmt.Fprintf(o, "/* %q not found: %v */\n", v, err)
		}
		/* Save it */
		ss = append(ss, s)
		ls = append(ls, end-start)
	}
	/* Time the whole thing took */
	d := time.Now().Sub(started)
	log.Printf(
		"Found seeds for %q in %v (%v/byte)",
		v,
		d,
		d/time.Duration(len(v)),
	)

	/* Once we've got all the seeds, print a nice line of Go */
	gcFound(v, ss, ls, o)
}

/* gcFound emits code to be used as a string replacement.  The seeds needed
to make the substrings go in ss, and the lengths of the corresponding
substrings in ls.  Output will be written to o. */
func gcFound(v []byte, ss []int64, ls []int, o *os.File) {
	/* Variable name form of v */
	vn := strconv.Quote(string(v))
	vn = strings.Map(func(r rune) rune {
		if (r >= '0' && r <= '9') ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') {
			return r
		}
		return '_'
	}, vn[1:len(vn)-1])
	/* Append a number if we need */
	if n, ok := varnames[vn]; ok {
		vn = fmt.Sprintf("%v_%v", vn, n)
		varnames[vn] = n + 1
	}
	/* Append a _randstr to note it's one of these */
	vn += "_rs"
	/* Print go code for variable */
	fmt.Fprintf(o, "/* %q */\n", v)
	fmt.Fprintf(o, "var %v = \"\"", vn)
	for i, s := range ss {
		fmt.Fprintf(o, " + seedToString(%v, %v)", s, ls[i])
	}
	fmt.Fprintf(o, "\n")
}
