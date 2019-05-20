package app

import (
	"os"
	"sort"

	"github.com/navinds25/windiskmgmt/internal/dfcli"

	"github.com/navinds25/windiskmgmt/internal/dfconfig"
	"github.com/navinds25/windiskmgmt/pkg/diskdata"
	log "github.com/sirupsen/logrus"
)

// GetInfoConfFiles gets the file Info
// when reading from duplicate files configuration.
func GetInfoConfFiles(files []string) error {
	log.Info("files", files)
	fileMap := make(map[string][]diskdata.FileInfo)
	for _, file := range files {
		fInfo, err := dfconfig.GetFileInfo(file)
		if err != nil {
			if err.Error() == "unable to open file" {
				log.Error(err)
				continue
			} else {
				log.Error(err)
				return err
			}
		}
		fileMap[fInfo.Checksum] = append(fileMap[fInfo.Checksum], fInfo)
	}
	if err := diskdata.DataStore.AddFileMap(fileMap); err != nil {
		return err
	}
	return nil
}

// ReadDBFiles shows the list of files in the DB.
func ReadDBFiles() error {
	dbFiles, err := diskdata.DataStore.ReadAllFiles()
	if err != nil {
		return err
	}
	log.Info(dbFiles)
	return nil
}

// ProcessDBFiles reads from the db and returns a map based on checksum
func ProcessDBFiles() error {
	dbFiles, err := diskdata.DataStore.ReadAllFiles()
	if err != nil {
		return err
	}
	for _, files := range dbFiles {
		log.Println("this thing ?:", files)
		filesWStatus, err := compareFiles(files.Files)
		if err != nil {
			return err
		}
		log.Println("delFiles: ", filesWStatus)
		if err := deleteFiles(files.Files); err != nil {
			return err
		}
	}
	return nil
}

// CompareFiles compares files of the same size.
func compareFiles(input []diskdata.FileInfo) (output []diskdata.FileInfo, err error) {
	for _, file := range input {
		if err := CheckHighPriorityFolders(&file); err != nil {
			return nil, err
		}
		if err := CheckLowPriorityFiles(&file); err != nil {
			return nil, err
		}
	}
	sort.Slice(input, func(i, j int) bool {
		return input[i].Priority < input[j].Priority
	})
	input[0].DelStatus = "keep"
	for i := 1; i < len(input); i++ {
		input[i].DelStatus = "delete"
	}
	return input, nil
}

func deleteFiles(input []diskdata.FileInfo) error {
	for _, file := range input {
		if file.DelStatus == "delete" {
			if !dfcli.CliFlags.Dryrun {
				if err := os.Remove(file.File); err != nil {
					return err
				}
			} else {
				log.WithField("delFile", file.File)
			}
			log.Info("deleted file: ", file.File)
		} else {
			log.Info("keeping file: ", file.File)
		}
	}
	return nil
}

// LoadToDB loads a config file to DB.
func LoadToDB(duplicateFilesConf string) error {
	dfconf, err := dfconfig.GetConfig(duplicateFilesConf)
	if err != nil {
		return err
	}
	for key, files := range dfconf.Files {
		log.Info("Processing for key:", key)
		err := GetInfoConfFiles(files)
		if err != nil {
			return err
		}
	}
	return nil
}
