package diskdata

import (
	"bufio"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus" // need to remove the logging or move to internal
)

// DiskFile is the interface for FileInfo
type DiskFile interface {
	GetInfo(string) error
	DeleteFile(string, bool) error
	GetChecksum() error
	ToDB()
	CheckDB()
}

// FileInfo struct for file info
type FileInfo struct {
	File     string
	Basename string
	Checksum uint32
	Size     int64
	Modified time.Time
}

// GetInfo returns FileInfo struct with info on files
func (file *FileInfo) GetInfo(filename string) error {
	theFile, err := os.Open(filename)
	if err != nil {
		return errors.New("unable to open file")
	}
	defer theFile.Close()
	fileDetails, err := theFile.Stat()
	if err != nil {
		return err
	}
	fileSize := fileDetails.Size()
	err = file.GetChecksum()
	if err != nil {
		return err
	}
	file = &FileInfo{
		File:     filename,
		Basename: path.Base(filename),
		Checksum: file.Checksum,
		Size:     fileSize,
		Modified: fileDetails.ModTime(),
	}
	return nil
}

// DeleteFile moves the file to a delete directory.
func (file *FileInfo) DeleteFile(DelDir string, Dryrun bool) error {
	_, err := os.Stat(DelDir)
	if err != nil {
		if err := os.Mkdir(DelDir, 0755); err != nil {
			return err
		}
	}
	newlocation := path.Join(DelDir, file.Basename)
	if !Dryrun {
		if err := os.Rename(file.File, newlocation); err != nil {
			return err
		}
	}
	log.Printf("Moving file %s to %s \n", file.File, newlocation)
	return nil
}

// GetChecksum returns the checksum of entire file
func (file *FileInfo) GetChecksum() error {
	fd, err := os.Open(file.File)
	if err != nil {
		return err
	}
	defer fd.Close()
	crc32q := crc32.MakeTable(0xD5828281)
	data, err := ioutil.ReadAll(fd)
	if err != nil {
		return err
	}
	file.Checksum = crc32.Checksum(data, crc32q)
	return nil
}

// ToDB adds the file to the DB
func (file *FileInfo) ToDB() error {
	return nil
}

// CheckDB checks the DB for the file
func (file *FileInfo) CheckDB() error {
	return nil
}

// Unused func...

// Get5mChecksum returns the checksum of files, max 5MB
func Get5mChecksum(file io.Reader, filesize int64) (uint32, error) {
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
