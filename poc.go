package main

import (
	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

func Run() {
	db, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	tx, err := db.Begin(true)
	if err != nil {
		log.Fatal(err)
	}
	bucket, err := tx.CreateBucket([]byte("dev"))
	if err != nil {
		log.Fatal(err)
	}
	if err := bucket.Put([]byte("hello"), []byte("42\n")); err != nil {
		log.Fatal(err)
	}
	resp := bucket.Get([]byte("hello"))
	if err := bucket.Put([]byte("hello"), append(resp, []byte("47\n")...)); err != nil {
		log.Fatal(err)
	}
	resp = bucket.Get([]byte("hello"))
	log.Println(string(resp))
	if err := tx.Commit(); err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
