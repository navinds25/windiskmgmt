package cmd

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/boltdb/bolt"
	"github.com/kr/fs"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

// Dryrun for controlling dryrun flag & operation
var Dryrun bool

// Action for command type
var Action string

// StartDir root dir for searching
var StartDir string

// DelDir destination directory for files to be deleted.
var DelDir string

// Debug for controlling debug flag & operation
var Debug bool

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

// DBTransaction is for db transactions
func DBTransaction(bucket *bolt.Bucket, filepath string) error {
	//var files1 []FileInfo
	files1, err := GetFileInfo(filepath)
	if err != nil {
		return err
	}
	fileInfo := files1[0]
	baseResp := bucket.Get([]byte(fileInfo.Basename))
	if baseResp == nil {
		var value bytes.Buffer
		enc := gob.NewEncoder(&value)
		if err = enc.Encode(files1); err != nil {
			return err
		}
		if err = bucket.Put([]byte(fileInfo.Basename), value.Bytes()); err != nil {
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
		res, delfile := FileMatch(&fileInfo2, &fileInfo)
		if res {
			if err := DeleteFile(delfile); err != nil {
				return err
			}
		}
		//files2 = append(files2, fileInfo)

	}
	return nil
}

// DeleteFile moves the file to a delete directory.
func DeleteFile(file *FileInfo) error {
	_, err := os.Stat(DelDir)
	if err != nil {
		if err := os.Mkdir(DelDir, 0755); err != nil {
			return err
		}
	}
	newlocation := path.Join(DelDir, file.Basename)
	if !Dryrun {
		if err := os.Rename(file.File, newlocation); err != nil {
			log.Error(err)
		}
	}
	log.Printf("Moving file %s to %s \n", file.File, newlocation)
	return nil
}

// FileMatch matches the files
func FileMatch(file1, file2 *FileInfo) (bool, *FileInfo) {
	if file1.Checksum == file2.Checksum && file1.Size == file2.Size {
		file1mod := file1.Modified.Unix()
		file2mod := file2.Modified.Unix()
		if file1mod >= file2mod {
			return true, file2
		} else {
			return true, file1
		}
	} else {
		return false, nil
	}
}

// FileFind Finds the files
func FileFind(bucket *bolt.Bucket) error {
	directory := StartDir
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
			if err := DBTransaction(bucket, filepath); err != nil {
				return err
			}
		}
	}
	return nil
}

// GetChecksum returns the checksum of entire file
func GetChecksum(file io.Reader) (uint32, error) {
	crc32q := crc32.MakeTable(0xD5828281)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, err
	}
	return crc32.Checksum(data, crc32q), nil
}

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

// GetFileInfo returns FileInfo struct with info on files
func GetFileInfo(filename string) ([]FileInfo, error) {
	var fileInfo FileInfo
	var files []FileInfo
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
	checksum, err := Get5mChecksum(file, fileSize)
	if err != nil {
		return nil, err
	}
	fileInfo = FileInfo{
		File:     filename,
		Basename: path.Base(filename),
		Checksum: checksum,
		Size:     fileSize,
		Modified: fileDetails.ModTime(),
	}
	retfiles := append(files, fileInfo)
	return retfiles, nil
}

var deleteDuplicates = cli.Command{
	Name:    "dd",
	Aliases: []string{"delete_duplicates"},
	Usage:   "Move duplicate files to a delete folder.",
	Action: func(c *cli.Context) error {
		Action = "dd"
		return nil
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logging",
			Destination: &Debug,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "Disable/Enable dryrun",
			Destination: &Dryrun,
		},
		cli.StringFlag{
			Name:        "startdir",
			Usage:       "directory to start from",
			Destination: &StartDir,
		},
		cli.StringFlag{
			Name:        "deldir",
			Usage:       "directory to collect files to be deleted",
			Destination: &DelDir,
		},
	},
}

// App for all commandline arguments
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "windiskmgmt"
	app.Usage = "For finding duplicate files & deleting them"
	app.Commands = []cli.Command{
		deleteDuplicates,
	}
	return app
}

// Run main function for the application
func Run() error {
	database, err := bolt.Open("data.db", 0600, nil)
	if err != nil {
		return err
	}
	defer database.Close()
	tx, err := database.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Commit()
	bucket, err := tx.CreateBucketIfNotExists([]byte("files"))
	if err != nil {
		return err
	}
	err = FileFind(bucket)
	if err != nil {
		return err
	}
	return nil
}
