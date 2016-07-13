package main

/*
 * main.go
 * Generated PRNG seeds to represent strings
 * By Henri Koski, github.com/heppu/seed-finder
 * Modified by J. Stuart McMurray
 * Created 20160712
 * Last Modified 20160712
 */

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

const (
	MAX_UINT64 uint64 = ^uint64(0)
	MAX_INT64  int64  = int64(MAX_UINT64 >> 1)
	MIN_INT64  int64  = -MAX_INT64 - 1
)

var EXPECTED = flag.String("w", "", "word")

func main() {
	var (
		dbFile = flag.String(
			"db",
			"/var/db/seed-finder.db",
			"Database `file`",
		)
		wordFile = flag.String(
			"f",
			"",
			"Name of `file` containing strings to convert, or "+
				"\"-\" to read from stdin",
		)
		zero = flag.Bool(
			"0",
			false,
			"Split on null bytes (\\0), not newlines",
		)
		goOut = flag.Bool(
			"go",
			false,
			"Output Go source code",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options] [string [string...]]

Finds a random number seed which can be used to recrate the original string,
optionally first checking a database (which will be updated with found seeds).

Options:
`,
			os.Args[0],
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* TODO: Update this */
	fmt.Printf("DO SOMETHING WITH:\n*dbFile: %v\n*wordFile: %v\n*zero: %v\n*goOut: %v\n", *dbFile, *wordFile, *zero, *goOut)

	if *EXPECTED == "" {
		flag.PrintDefaults()
		return
	}

	workers := runtime.NumCPU()
	gap := int64(MAX_UINT64 / uint64(workers))
	from := MIN_INT64

	var to int64
	if workers > 1 {
		to = MIN_INT64 + int64(gap)
	} else {
		to = MAX_INT64
	}

	log.Printf("Starting %d workers with gap: %d", workers, gap)

	wg := &sync.WaitGroup{}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		tmp := i
		go find(tmp, from, to, wg)
		from = to + 1
		to += gap + 1
	}

	wg.Wait()
}

func find(wNum int, from, to int64, wg *sync.WaitGroup) {
	defer wg.Done()
	var r *rand.Rand
	var s byte
	var full string

	log.Printf("Worker %d working from %d to %d: ", wNum, from, to)

	for i := from; i <= to; i++ {
		r = rand.New(rand.NewSource(i))
		full = ""
		for j := 0; j < len(*EXPECTED); j++ {
			s = byte(r.Intn(26) + 97)
			if s != (*EXPECTED)[j] {
				break
			}
			full += string(s)
		}

		if full == *EXPECTED {
			log.Fatalf("Result: '%s' found with seed: '%d'", *EXPECTED, i)
		}
	}
}
