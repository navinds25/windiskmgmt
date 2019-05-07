package diskdata

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"github.com/boltdb/bolt"
	"github.com/dgraph-io/badger"
)

// DataStore is for the Database instance
var DataStore Store

// InitDB Initializes the Database
func InitDB(s Store) {
	DataStore = s
}

// Store is the main interface for the backend
type Store interface {
	AddFile(*FileInfo) error
	CheckFileExists(*FileInfo) ([]FileInfo, bool, error)
}

// BadgerDB is the DB instance for BadgerDB
type BadgerDB struct {
	DB  *badger.DB
	WTX *badger.Txn
	RTX *badger.Txn
}

// BoltDB is the DB instance for BoltDB
type BoltDB struct {
	DB     *bolt.DB
	WTX    *bolt.Tx
	RTX    *bolt.Tx
	Bucket *bolt.Bucket
}

// AddFile is the method for adding a file entry
func (badgerdb BadgerDB) AddFile(file *FileInfo) error {
	var value bytes.Buffer
	fileSize := strconv.FormatInt(file.Size, 10)
	if err := gob.NewEncoder(&value).Encode(file); err != nil {
		return err
	}
	badgerdb.WTX = badgerdb.DB.NewTransaction(true)
	if err := badgerdb.WTX.Set([]byte(fileSize), value.Bytes()); err != nil {
		return err
	}
	if err := badgerdb.WTX.Commit(); err != nil {
		return err
	}
	return nil
}

// AddFile is the method for adding a file entry
func (boltdb BoltDB) AddFile(file *FileInfo) error {
	var value bytes.Buffer
	fileSize := strconv.FormatInt(file.Size, 10)
	if err := gob.NewEncoder(&value).Encode(file); err != nil {
		return err
	}
	if err := boltdb.Bucket.Put([]byte(fileSize), value.Bytes()); err != nil {
		return err
	}
	return nil
}
