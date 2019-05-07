package app

import (
	"bytes"
	"encoding/gob"
	"os"
	"strconv"
	"strings"

	"github.com/boltdb/bolt"
	"github.com/kr/fs"
	"github.com/navinds25/windiskmgmt/internal/dfcli"
)

func checkSkipDir(file string) bool {
	for _, dir := range dfcli.SkipDirectories {
		if strings.HasPrefix(file, dir) {
			return true
		}
	}
	return false
}

//// FileMatch matches the files returns true if matched + file to be deleted
//func FileMatch(file1, file2 *FileInfo) (bool, *FileInfo, error) {
//	file1Checksum, err := GetChecksum(file1)
//	if err != nil {
//		return false, nil, err
//	}
//	file1.Checksum = file1Checksum
//	file2Checksum, err := GetChecksum(file2)
//	if err != nil {
//		return false, nil, err
//	}
//	file2.Checksum = file2Checksum
//	if file1.Checksum == file2.Checksum && file1.Size == file2.Size {
//		file1mod := file1.Modified.Unix()
//		file2mod := file2.Modified.Unix()
//		if file1mod >= file2mod {
//			return true, file2, nil
//		} else {
//			return true, file1, nil
//		}
//	} else {
//		return false, nil, nil
//	}
//}

// DBTransaction is for db transactions
func DBTransaction(bucket, movedFiles *bolt.Bucket, filepath string) error {
	//var files1 []FileInfo
	files1, err := GetFileInfo(filepath)
	if err != nil {
		return err
	}
	fileInfo := files1[0]
	fileSizeStr := strconv.FormatInt(fileInfo.Size, 10)
	baseResp := bucket.Get([]byte(fileSizeStr))
	if baseResp == nil {
		var value bytes.Buffer
		enc := gob.NewEncoder(&value)
		if err = enc.Encode(files1); err != nil {
			return err
		}
		if err = bucket.Put([]byte(fileSizeStr), value.Bytes()); err != nil {
			return err
		}

	} else {
		var files2 []FileInfo
		baseRespReader := bytes.NewReader(baseResp)
		dec := gob.NewDecoder(baseRespReader)
		if err = dec.Decode(&files2); err != nil {
			return err
		}
		fileInfo2 := files2[0]
		res, delfile, err := FileMatch(&fileInfo2, &fileInfo)
		if err != nil {
			return err
		}
		if res {
			if !Dryrun {
				if err := DeleteFile(delfile); err != nil {
					return err
				}
				if err = movedFiles.Put([]byte(delfile.File), []byte(fileSizeStr)); err != nil {
					return err
				}
			} else {
				log.Info("Moving file: ", delfile.File)
				if err = movedFiles.Put([]byte(delfile.File), []byte(fileSizeStr)); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// FileFind Finds the files
func FileFind(bucketName, delBucketName []byte, db *bolt.DB) error {
	var count int64
	count = 0
	tx, err := db.Begin(true)
	if err != nil {
		return err
	}
	bucket := tx.Bucket(bucketName)
	movedFiles := tx.Bucket(delBucketName)
	walker := fs.Walk(dfcli.StartDir)
	for walker.Step() {
		if err := walker.Err(); err != nil {
			log.Error(err)
			continue
		}
		stat := walker.Stat()
		filename := walker.Path()
		filemode := stat.Mode()
		if stat.IsDir() {
			continue
		} else if checkSkipDir(filename) {
			log.Println("skipping file", filename)
			continue
		} else if filemode&os.ModeSymlink == os.ModeSymlink {
			log.Println("found symlink", filename)
			continue
		} else if filemode&os.ModeDevice == os.ModeDevice {
			log.Println("found device", filename)
			continue
		} else if filemode&os.ModeSocket == os.ModeSocket {
			log.Println("found socket", filename)
			continue
		} else {
			count, tx, bucket, movedFiles, err = GetDBTX(count, db, tx, bucket, movedFiles, bucketName, delBucketName)
			if err != nil {
				return err
			}
			filepath := walker.Path()
			if err := DBTransaction(bucket, movedFiles, filepath); err != nil {
				return err
			}
		}
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

// GetDBTX commits and returns new transactions for the db
func GetDBTX(count int64, db *bolt.DB, tx *bolt.Tx, bucket, movedFiles *bolt.Bucket, bucketName, delBucketName []byte) (int64, *bolt.Tx, *bolt.Bucket, *bolt.Bucket, error) {
	if count == 500 {
		if err := tx.Commit(); err != nil {
			return 0, nil, nil, nil, err
		}
		newtx, err := db.Begin(true)
		if err != nil {
			return 0, nil, nil, nil, err
		}
		newBucket := newtx.Bucket(bucketName)
		newMovedFiles := newtx.Bucket(delBucketName)
		count = 0
		return count, newtx, newBucket, newMovedFiles, nil
	} else {
		count++
		return count, tx, bucket, movedFiles, nil
	}
}

//// Run main function for the application
//func Run() error {
//	SkipDirectories = GetSkipDirectories()
//	bucketName := []byte("files")
//	delBucketName := []byte("moved_files")
//	database, err := bolt.Open("data.db", 0600, nil)
//	if err != nil {
//		return err
//	}
//	tx, err := database.Begin(true)
//	if err != nil {
//		return err
//	}
//	_, err = tx.CreateBucketIfNotExists(bucketName)
//	if err != nil {
//		return err
//	}
//	_, err = tx.CreateBucketIfNotExists(delBucketName)
//	tx.Commit()
//	err = FileFind(bucketName, delBucketName, database)
//	if err != nil {
//		return err
//	}
//	database.Close()
//	return nil
//}
