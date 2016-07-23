package main

/*
 * db.go
 * Handle the database
 * By J. Stuart McMurray
 * Created 20170712
 * Last Modified 20170712
 */

import (
	"crypto/sha512"
	"encoding/binary"

	"github.com/boltdb/bolt"
)

const BUCKETNAME = "seeds"

var DB *bolt.DB

/* dbOpen opens the database in the file named fn, and returns the database
as well as the bucket containing seeds. */
/* dbOpen opens the database in the file named fn.  Regardless of whether
the database has been opened, checkDB and updateDB will work */
func dbOpen(fn string) (*bolt.DB, error) {
	/* Create/open the database */
	db, err := bolt.Open(fn, 0644, nil)
	if nil != err {
		return nil, err
	}
	DB = db

	/* Make sure the table exists */
	if err := DB.Update(func(tx *bolt.Tx) error {
		/* Get the bucket */
		_, err := tx.CreateBucketIfNotExists([]byte(BUCKETNAME))
		return err
	}); nil != err {
		return nil, err
	}

	return db, nil
}

/* checkDB returns the seed for b, whether it's known to be unfindable, and
whether it was in the database */
func checkDB(b []byte) (seed int64, unfindable, found bool, err error) {
	/* Don't bother if we haven't a database */
	if nil == DB {
		return 0, false, false, nil
	}

	/* Get the seed, if there is one */
	var v []byte
	DB.View(func(tx *bolt.Tx) error {
		/* Get the bucket */
		bucket := tx.Bucket([]byte(BUCKETNAME))
		if nil == bucket {
			panic("No bucket")
		}
		/* Try to get the value stored for the key */
		v = bucket.Get(hashBuf(b))
		return nil
	})

	if nil == v {
		return 0, false, false, nil
	}
	if 0 == len(v) {
		return 0, true, true, nil
	}
	/* Convert to an int64 */
	seed, n := binary.Varint(v)
	if 0 >= n {
		return 0, false, false, nil
	}
	return seed, false, true, nil
}

/* storeSeed puts the seed for b in the database */
func storeSeed(seed int64, b []byte) error {
	/* Convert to a storable varint */
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, seed)
	buf = buf[:n]
	/* Stick it in the database */
	return storeDB(b, buf)
}

/* noteUnfindable notes that there is no seed for b in the database */
func noteUnfindable(b []byte) error {
	return storeDB(b, []byte{})
}

/* storeDB sticks the key and value in the database */
func storeDB(key, value []byte) error {
	if nil == DB {
		return nil
	}
	/* Stick it in the database */
	return DB.Update(func(tx *bolt.Tx) error {
		/* Get the bucket */
		bucket := tx.Bucket([]byte(BUCKETNAME))
		/* Try to get the value stored for the key */
		return bucket.Put(hashBuf(key), value)
	})
}

/* hashBuf returns a hash of b if b is longer than the maximum supported key
length, otherwise it returns b. */
func hashBuf(b []byte) []byte {
	/* All set if it's small enough */
	if bolt.MaxKeySize >= len(b) {
		return b
	}
	c := sha512.Sum512(b)
	b = make([]byte, len(c))
	for i, v := range c {
		b[i] = v
	}
	return append([]byte("seed-finder-key-"), b...)
}

/* nSeedStrings returns the number of strings in the database */
func nSeedStrings() int {
	if nil == DB {
		panic("Nil database")
	}
	var n int /* Number of key value pairs */
	DB.View(func(tx *bolt.Tx) error {
		n = tx.Bucket([]byte(BUCKETNAME)).Stats().KeyN
		return nil
	})
	return n
}
