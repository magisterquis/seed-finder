package main

/*
 * builddb.go
 * Build the database the easy way
 * By J. Stuart McMurray
 * Created 20160723
 * Last Modified 20160723
 */

import (
	"log"
	"math"
	"math/rand"
	"time"
)

/* buildDatabase builds the database up with strings of the given length by
iterating over seeds and sticking found strings in the database */
func buildDatabase(l uint, start int64, tgt string) {
	r := rand.New(rand.NewSource(0)) /* Random number generator */
	b := make([]byte, l)             /* String buffer */
	var err error
	startt := time.Now() /* Start time */
	lastt := time.Now()  /* Last update time */
	lasti := start       /* Last update seed */
	found := false       /* True when found the target */
	/* Iterate through all the seeds */
SEEDS:
	for seed := start; seed <= math.MaxInt64; seed++ {
		/* Reseed */
		r.Seed(seed)
		/* Get the bytes */
		for i := range b {
			b[i] = byte(r.Intn(256))
			/* Only accept ascii characters */
			if ' ' > b[i] || '~' < b[i] {
				continue SEEDS
			}
		}
		if err = storeSeed(seed, b); nil != err {
			log.Fatalf("Unable to store %q->%v: %v", b, seed, err)
		}

		/* If we've found our target, give up */
		if "" != tgt {
			if string(b) == tgt {
				found = true
			}
		}

		/* Print a summary every 15 seconds */
		if time.Now().Sub(lastt) > (15*time.Second) || found {
			now := time.Now()
			log.Printf(
				"Last: %s -> %v.  Finished %v in %v (%0.2f/s).  %v in database.",
				b, seed,
				seed-lasti,
				now.Sub(startt),
				float64(seed-lasti)/
					time.Now().Sub(lastt).Seconds(),
				nSeedStrings(),
			)
			if found {
				break
			}
			lastt = now
			lasti = seed
		}
	}
}
