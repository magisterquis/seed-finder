package main

/*
 * main.go
 * Generated PRNG seeds to represent strings
 * By Henri Koski, github.com/heppu/seed-finder
 * Modified by J. Stuart McMurray
 * Created 20160712
 * Last Modified 20160723
 */

import (
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/boltdb/bolt"
)

func main() {
	var (
		dbFile = flag.String(
			"db",
			"/var/db/seed-finder.db",
			"Database `file`",
		)
		/* Input control */
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
		/* Output control */
		oFileName = flag.String(
			"o",
			"",
			"Write output to this file instead of stdout",
		)
		/* Code-generation control */
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
		pkgName = flag.String(
			"gopkg",
			"main",
			"Package name to use when outputting Go source",
		)
		/* Database building */
		buildDB = flag.Uint(
			"b",
			0,
			"Pre-build a database with ascii strings of this "+
				"`length`",
		)
		buildStart = flag.Int64(
			"bstart",
			math.MinInt64,
			"Starting `seed` for -b",
		)
		buildNonPrint = flag.Bool(
			"bnonprint",
			false,
			"Include ASCII characters outside the range "+
				"0x20-0x7E (space - tilde) when building the "+
				"database",
		)
	)
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: %v [options] [string [string...]]

Finds the random seed needed to generate the given string(s), mainly for use
in removing readable strings from Go binaries.  A database is used to speed
up lookups.  To further decrease runtime, the database can be pre-populated.

Strings must not be longer than %v bytes.

Options:
`,
			os.Args[0],
			bolt.MaxKeySize,
		)
		flag.PrintDefaults()
	}
	flag.Parse()

	/* Open the database */
	db, err := dbOpen(*dbFile)
	if nil != err {
		log.Fatalf("Database error with %v: %v", *dbFile, err)
	}
	if nil != db {
		log.Printf("Opened database %v", *dbFile)
		defer db.Close()
	}

	/* If we're just to build a database, do that */
	if 0 != *buildDB {
		buildDatabase(*buildStart, *buildDB, *buildNonPrint)
		return
	}

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

	/* Open the output file */
	ofile := os.Stdout
	if "" != *oFileName {
		ofile, err = os.OpenFile(
			*oFileName,
			os.O_TRUNC|os.O_CREATE|os.O_APPEND|os.O_WRONLY,
			0644,
		)
		if nil != err {
			log.Fatalf(
				"Unable to open output file %v: %v",
				*oFileName,
				err,
			)
		}
		defer ofile.Close()
	}

	/* If we're printing out Go code, print the boilerplate */
	if *goCode {
		gcBoilerplate(ofile, *pkgName)
	}

	/* Process the words to stringify */
	for in := range ins {
		/* Go's got it's own functions.  Fancy. */
		if *goCode {
			gcVar(in, *goLen, ofile)
			continue
		}
		/* Find the seeds for the strings if we're not making Go. */
		seed, err := findSeed(in, math.MinInt64, 0, 0xFF)
		if nil != err {
			log.Printf("%q ERROR: %v", in, err)
		} else {
			fmt.Fprintf(ofile, "%s -> %v\n", in, seed)
		}
		continue
	}
	/* TODO: Check from here */
	//
	//
	//	/* Find each seed */
	//	for in := range ins {
	//		seed, err := findSeed(in, *nParallel)
	//		if nil != err {
	//			log.Printf("Unable to find seed for %v: %v", in, err)
	//			continue
	//		}
	//		fmt.Fprintf(ofile, "%q %v", in, seed)
	//	}
}
