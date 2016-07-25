package main

/*
 * builddb.go
 * Build the database the easy way
 * By J. Stuart McMurray
 * Created 20160723
 * Last Modified 20160723
 */

import (
	"fmt"
	"log"
)

/* buildDatabase ensures all the strings of length l are in the database,
starting with seed start.  If nonprint is true, all strings, even those with
nonprintable characetrs, are included. */
func buildDatabase(start int64, l uint, nonprint bool) {
	/* Starting and ending characters */
	cstart := byte(' ')
	cend := byte('~')
	if nonprint {
		cstart = 0
		cend = 0xFF
	}
	/* String to find */
	s := make([]byte, l)
	for i := range s {
		s[i] = cstart
	}
	/* Starting seed */
	seed := start

	/* Loop over each string */
	for done := false; !done; done = incStr(s, cstart, cend) {

		/* Skip string if we have it already */
		if dbHas(s) {
			//log.Printf("Found %q", s) /* DEBUG */
			continue
		}
		log.Printf("Finding %q (% X)", s, s) /* DEBUG */

		/* Find it (and stick it in) if we don't */
		fs, err := findSeed(s, seed, cstart, cend)
		if nil != err {
			e := fmt.Sprintf(
				"Unable to find seed for %q: %v",
				s,
				err,
			)
			if ErrUnfindable != err {
				log.Fatalf("%v", e)
			}
			log.Printf("%v", e)
			continue
		}
		log.Printf("Found %q -> %v", s, fs) /* DEBUG */

	}
	/* TODO: Timing info */
}

/* incStr increments the string in s, using start and end as the bounds for
each byte in s.  It returns true on overflow */
func incStr(s []byte, start, end byte) bool {
	for i := len(s) - 1; 0 <= i; i-- {
		/* Check for overflow on this bytes */
		if end == s[i] {
			/* Check for string overflow */
			if 0 == i {
				return true
			}
			s[i] = start
			continue
		}
		s[i]++
		return false
	}
	return true
}

/* TODO: Fix or remove below */
//func delme() {
//	r := rand.New(rand.NewSource(0)) /* Random number generator */
//	b := make([]byte, l)             /* String buffer */
//	var err error
//	startt := time.Now() /* Start time */
//	lastt := time.Now()  /* Last update time */
//	lasti := start       /* Last update seed */
//	found := false       /* True when found the target */
//	/* Iterate through all the seeds */
//SEEDS:
//	for seed := start; seed <= math.MaxInt64; seed++ {
//		/* Reseed */
//		r.Seed(seed)
//		/* Get the bytes */
//		for i := range b {
//			b[i] = byte(r.Intn(256))
//			/* Only accept ascii characters */
//			if ' ' > b[i] || '~' < b[i] {
//				continue SEEDS
//			}
//		}
//		if err = storeSeed(seed, b); nil != err {
//			log.Fatalf("Unable to store %q->%v: %v", b, seed, err)
//		}
//
//		/* If we've found our target, give up */
//		if "" != tgt {
//			if string(b) == tgt {
//				found = true
//			}
//		}
//
//		/* Print a summary every 15 seconds */
//		if time.Now().Sub(lastt) > (15*time.Second) || found {
//			now := time.Now()
//			log.Printf(
//				"Last: %s -> %v.  Finished %v in %v (%0.2f/s).  %v in database.",
//				b, seed,
//				seed-lasti,
//				now.Sub(startt),
//				float64(seed-lasti)/
//					time.Now().Sub(lastt).Seconds(),
//				nSeedStrings(),
//			)
//			if found {
//				break
//			}
//			lastt = now
//			lasti = seed
//		}
//	}
//}
