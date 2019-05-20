package diskdata

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/dgraph-io/badger"
	log "github.com/sirupsen/logrus"
)

// DataStore is for the Database instance
var DataStore Store

// InitDB Initializes the Database
func InitDB(s Store) {
	DataStore = s
}

// Store is the main interface for the backend
type Store interface {
	CloseDB() error
	AddFile(*FileInfo) error
	AddFileMap(map[string][]FileInfo) error
	//CheckFileExists(*FileInfo) ([]FileInfo, bool, error)
	//ReadAllFiles() (<-chan *[]FileInfo, error)
	ReadAllFiles() ([]*FileList, error)
}

// BadgerDB is the DB instance for BadgerDB
type BadgerDB struct {
	DB  *badger.DB
	WTX *badger.Txn
	RTX *badger.Txn
}

type FileList struct {
	Files []FileInfo
}

// AddFile is the method for adding a file entry
func (badgerdb BadgerDB) AddFile(file *FileInfo) error {
	var value bytes.Buffer
	if err := gob.NewEncoder(&value).Encode(file); err != nil {
		return err
	}
	badgerdb.WTX = badgerdb.DB.NewTransaction(true)
	return nil
}

// AddFileMap takes a map of files by  checksum and adds it to the database.
func (badgerdb BadgerDB) AddFileMap(fileMap map[string][]FileInfo) error {
	badgerdb.WTX = badgerdb.DB.NewTransaction(true)
	for checksum, files := range fileMap {
		var value bytes.Buffer
		fileList := FileList{
			Files: files,
		}
		if err := gob.NewEncoder(&value).Encode(fileList); err != nil {
			return err
		}
		log.Infof("checksum: %s, files: %v", checksum, files)
		if err := badgerdb.WTX.Set([]byte(checksum), value.Bytes()); err != nil {
			return err
		}
	}
	if err := badgerdb.WTX.Commit(); err != nil {
		return err
	}
	return nil
}

// CloseDB Closes the Database, should be called only from main() with defer.
func (badgerdb BadgerDB) CloseDB() error {
	if err := badgerdb.DB.Close(); err != nil {
		return err
	}
	return nil
}

// ReadAllFilesChan is the injest function for the file processing pipeline.
//func (badgerdb BadgerDB) ReadAllFilesChan() (<-chan *]FileInfo, error) {
//	out := make(chan *[]FileInfo)
//	badgerdb.DB.View(func(txn *badger.Txn) error {
//		opts := badger.DefaultIteratorOptions
//		opts.PrefetchSize = 10
//		it := txn.NewIterator(opts)
//		defer it.Close()
//		for it.Rewind(); it.Valid(); it.Next() {
//			item := it.Item()
//			k := item.Key()
//			valCopy, err := item.ValueCopy(nil)
//			if err != nil {
//				return err
//			}
//			//err := item.Value(func(v []byte) error {
//			fValue := &FileList{}
//			valReader := bytes.NewReader(valCopy)
//			if err := gob.NewDecoder(valReader).Decode(fValue); err != nil {
//				return err
//			}
//			fmt.Printf("key=%s, value=%v\n", k, fValue)
//			out <- fValue.Files
//			defer close(out)
//		}
//		return nil
//	})
//	return out, nil
//}

// ReadAllFilesArray is the injest function for the file processing pipeline.
func (badgerdb BadgerDB) ReadAllFiles() (out []*FileList, err error) {
	//out := [][]FileInfo{}
	badgerdb.DB.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			//k := item.Key()
			valCopy, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			fValue := &FileList{}
			valReader := bytes.NewReader(valCopy)
			if err := gob.NewDecoder(valReader).Decode(fValue); err != nil {
				return err
			}
			fmt.Printf(" value=%v\n", fValue)
			out = append(out, fValue)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return out, nil
}
