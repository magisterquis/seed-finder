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
	"os"
	"runtime"

	"github.com/boltdb/bolt"
)

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
			"Split on null bytes (\\0), not newlines with -f",
		)
		goCode = flag.Bool(
			"go",
			false,
			"Output Go source code",
		)
		goLen = flag.Uint(
			"golen",
			3,
			"Split strings into chunks of this `length` if Go "+
				"source code is to be output",
		)
		nParallel = flag.Uint(
			"n",
			uint(runtime.NumCPU()),
			"Split seed-finding into `count` parallel attempts",
		)
		oFileName = flag.String(
			"o",
			"",
			"Write output to this file instead of stdout",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options] [string [string...]]

Finds a random number seed which can be used to recrate the original string,
optionally first checking a database (which will be updated with found seeds).

If Go source is to be output, strings will be split into smaller chunks, to
considerably speed up searching.

Strings must not be longer than %v bytes.

Options:
`,
			os.Args[0],
			bolt.MaxKeySize,
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Get the strings or chunks for which to find seeds */
	ins, err := getInput(flag.Args(), *wordFile, *zero)
	if nil != err {
		log.Fatalf("Unable to get input: %v", err)
	}
	if nil == ins {
		fmt.Fprintf(os.Stderr, "No strings to convert.\n\n")
		flag.Usage()
		os.Exit(1)
	}

	/* Open the database */
	db, err := dbOpen(*dbFile)
	if nil != err {
		log.Fatalf("Database error with %v: %v", *dbFile, err)
	}
	if nil != db {
		log.Printf("Opened database %v", *dbFile)
		defer db.Close()
	}

	/* Open the output file */
	ofile := os.Stdout
	if "" != *oFileName {
		ofile, err = os.OpenFile(
			*oFileName,
			os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if nil != err {
			log.Fatalf(
				"Unable to open output file %v: %v",
				*oFileName,
				err,
			)
		}
	}

	/* If we're printing out Go code, print the boilerplate */
	if *goCode {
		gcBoilerplate(ofile)
		for in := range ins {
			gcVar(in, *goLen, *nParallel, ofile)
		}
		return
	}

	/* Find each seed */
	for in := range ins {
		seed, err := findSeed(in, *nParallel)
		if nil != err {
			log.Printf("Unable to find seed for %v: %v", in, err)
			continue
		}
		fmt.Fprintf(ofile, "%q %v", in, seed)
	}
}
