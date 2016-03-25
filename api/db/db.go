// Package db provides functions for working with the database.
package db

import (
	"github.com/boltdb/bolt"
	"log"
)

// DB is the path to the database file
const DB string = "db/db.db"

// Open opens up a connection to the database. Callers should be careful
// to close the connection after they are done with it.
func Open() *bolt.DB {
	db, err := bolt.Open(DB, 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

// Setup sets up the appropriate database schema of buckets.
func Setup() {
	db := Open()
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte("Users"))
		tx.CreateBucketIfNotExists([]byte("Meta"))
		return nil
	})
}
