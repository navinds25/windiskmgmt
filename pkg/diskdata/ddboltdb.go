package diskdata

import (
	"bytes"
	"encoding/gob"
	"strconv"

	"github.com/boltdb/bolt"
)

// BoltDB is the DB instance for BoltDB
type BoltDB struct {
	DB     *bolt.DB
	WTX    *bolt.Tx
	RTX    *bolt.Tx
	Bucket *bolt.Bucket
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
