package main

import (
	"bufio"
	"database/sql"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/kr/fs"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
)

// FileInfo struct for file info
type FileInfo struct {
	File     string
	Basename string
	Checksum uint32
	Size     int64
	Modified time.Time
}

func skipDir(file string) bool {
	skipDirs := []string{"/proc", "/run", "/sys", "/tmp"}
	for _, dir := range skipDirs {
		if strings.HasPrefix(file, dir) {
			return true
		}
	}
	return false
}

// FileFind Finds the files
func FileFind() error {
	directory := "/"
	walker := fs.Walk(directory)
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
		} else if skipDir(filename) {
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
			filepath := walker.Path()
			log.Println(filepath)
			fileInfo, err := GetFileInfo(filepath)
			if err != nil {
				return err
			}
			log.Println(fileInfo)
		}
	}
	return nil
}

// SQLITEFileFind Finds the files
func SQLITEFileFind(insertStatement, selectStatement *sql.Stmt) error {
	directory := "/home/navin"
	walker := fs.Walk(directory)
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
		} else if skipDir(filename) {
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
			filepath := walker.Path()
			log.Println(filepath)
			fileInfo, err := GetFileInfo(filepath)
			if err != nil {
				return err
			}
			log.Println(fileInfo)
			insertStatement.Exec(fileInfo.File, fileInfo.Basename, fileInfo.Checksum, fileInfo.Size, fileInfo.Modified)
		}
	}
	return nil
}

// GetChecksum returns the checksum of files, max 5m
func GetChecksum(file io.Reader, filesize int64) (uint32, error) {
	crc32q := crc32.MakeTable(0xD5828281)
	if filesize < 5000000 {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return 0, err
		}
		return crc32.Checksum(data, crc32q), nil
	} else {
		reader := bufio.NewReaderSize(file, 5000000)
		data, err := reader.Peek(5000000)
		if err != nil {
			return 0, err
		}
		checksum := crc32.Checksum(data[:5000000], crc32q)
		reader.Reset(file)
		return checksum, nil
	}
}

// GetFileInfo returns FileInfo struct with info on files
func GetFileInfo(filename string) (*FileInfo, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileDetails, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileDetails.Size()
	checksum, err := GetChecksum(file, fileSize)
	if err != nil {
		return nil, err
	}
	fileInfo := FileInfo{
		File:     filename,
		Basename: path.Base(filename),
		Checksum: checksum,
		Size:     fileSize,
		Modified: fileDetails.ModTime(),
	}
	return &fileInfo, nil
}

func main() {
	//err := FileFind()
	//if err != nil {
	//	log.Error(err)
	//}
	database, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer database.Close()
	tablestatement, err := database.Prepare("CREATE TABLE IF NOT EXISTS files (id INTEGER PRIMARY KEY, file TEXT, basename TEXT, checksum TEXT, size INTEGER, modified INTEGER);")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	result, err := tablestatement.Exec()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	insertStatement, err := database.Prepare("INSERT INTO files (file, basename, checksum, size, modified) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	selectStatement, err := database.Prepare("SELECT * FROM files WHERE file = ?")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	err = SQLITEFileFind(insertStatement, selectStatement)
	if err != nil {
		log.Error(err)
	}
	log.Info(result.RowsAffected())
	log.Info(result.LastInsertId())
}
