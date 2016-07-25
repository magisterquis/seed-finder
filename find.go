package main

/*
 * find.go
 * Find the seed for a string
 * By J. Stuart McMurray
 * Created 20170712
 * Last Modified 20170723
 */

import (
	"fmt"
	"log"
	"math"
	"math/rand"
)

/* ErrUnfindable is returned when a seed for a string does not exist. */
var ErrUnfindable = fmt.Errorf("no matching seed exists")

/* findSeed searches for a seed for s, starting at start.  It stores in the
database string -> seed pairs it finds during the search.  The search will
be limited to bytes between startc and endc, inclusive. */
func findSeed(s []byte, start int64, startc, endc byte) (int64, error) {
	/* Check the database */
	seed, unfindable, found := checkDB(s)
	if unfindable { /* We already tried */
		return 0, ErrUnfindable
	}
	if found { /* Seed was found */
		return seed, nil
	}

	/* Not found, search for seed for s */
	r := rand.New(rand.NewSource(0)) /* Random number generator */
	b := make([]byte, len(s))        /* String buffer */
SEEDS:
	for seed := start; seed <= math.MaxInt64; seed++ {
		/* Reseed */
		r.Seed(seed)
		/* Get the bytes */
		for i := range b {
			b[i] = byte(r.Intn(256))
			/* Only accept ascii characters */
			if startc > b[i] || endc < b[i] {
				continue SEEDS
			}
		}
		/* Don't bother if we already have it */
		if dbHas(b) {
			continue
		}
		/* Stick it in the database for later */
		if err := storeSeed(seed, b); nil != err {
			return 0, err
		}
		/* If it's the one we're looking for, return it */
		for i, v := range b {
			if s[i] != v {
				continue SEEDS
			}
		}
		return seed, nil
	}
	/* Seed not found for s, note as unfindable if we really tried */
	if math.MinInt64 == start {
		/* Note in the database no worky */
		if err := noteUnfindable(s); nil != err {
			return 0, err
		}
	}
	return 0, ErrUnfindable
}

/* TODO: Review below */

///* findSeed finds the seed for v using npar guessing goroutines.  If goCode is
//true, Go source code will be emitted, which can be used in other programs */
//func findSeed(v []byte, workers uint) (int64, error) {
//		/* Get a pointer to the seed for v */
//		p := calculate(v, workers)
//		/* If it's unfindable, note it and move on */
//		if nil == p {
//			logUnfindable(v)
//			return 0, ErrUnfindable
//		}
//		s = *p
//		storeSeed(s, v)
//	}
//
//	return s, nil
//
//}

/* logUnfindable makes appropriate notes and such that there is no seed for v */
func UNUSEDlogUnfindable(v []byte) {
	log.Printf("No seed found for %q", v)
	if err := noteUnfindable(v); nil != err {
		log.Printf(
			"Unable to note %q unfindable in database: %v",
			v,
			err,
		)
	}
}

///* calculate performs an exhaustive search over the entire int64 range, and
//returns the seed for v, or nil if no seed was found */
//func UNUSEDcalculate(v []byte, workers uint) *int64 {
//	var (
//		/* Number of guesses for each worker to do */
//		gap  = int64(MAX_UINT64 / uint64(workers))
//		from = MIN_INT64 /* Starting guess */
//		to   int64       /* Ending guess */
//	)
//
//	/* Initial ending guess */
//	if workers > 1 {
//		to = MIN_INT64 + int64(gap)
//	} else {
//		to = MAX_INT64
//	}
//
//	log.Printf("Starting %d workers with gap: %d", workers, gap)
//
//	/* Fire off a handful of workers */
//	done := false
//	ch := make(chan *int64, workers)
//	for i := uint(0); i < workers; i++ {
//		/* If it's the last iteration, make sure we cover enough */
//		if workers-1 == i {
//			to = MAX_INT64
//		}
//		go find(v, i, from, to, &done, ch)
//		from = to + 1
//		to += gap + 1
//	}
//	/* Wait four our workers to finish */
//	var seed *int64
//	for i := uint(0); i < workers; i++ {
//		s := <-ch
//		/* Get the worker's result */
//		if nil != s {
//			seed = s
//			done = true
//		}
//	}
//	return seed
//}
//
///* find is an individual worker looking for a seed for v.  If a seed is found,
//it will be sent back on ch.  Execution will terminate shortly after done is
//true. */
//func UNUSEDfind(v []byte, wNum uint, from, to int64, done *bool, ch chan<- *int64) {
//	var r *rand.Rand
//	var s byte
//
//	// log.Printf("Worker %d working from %d to %d: ", wNum, from, to) /* DEBUG */
//
//	/* Search through each seed */
//	for i := from; !*done && i <= to; i++ {
//		/* Seed the PRNG with the guess */
//		r = rand.New(rand.NewSource(i))
//		found := true /* Will be set to false if this fails */
//		/* Try each byte pair */
//		for _, b := range v {
//			s = byte(r.Intn(256))
//			if s != b {
//				found = false
//				break
//			}
//		}
//		/* If it wasn't set to false, we win! */
//		if found {
//			ch <- &i
//			return
//		}
//	}
//	ch <- nil
//}
