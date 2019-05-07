package app

import (
	"log"

	"github.com/navinds25/windiskmgmt/internal/dfconfig"
	"github.com/navinds25/windiskmgmt/pkg/diskdata"
)

// GetInfoConfFiles gets the file Info
// when reading from duplicate files configuration.
func GetInfoConfFiles(files []string) ([]diskdata.FileInfo, error) {
	filesInfo := []diskdata.FileInfo{}
	for _, file := range files {
		fInfo, err := dfconfig.GetFileInfo(file)
		if err != nil {
			if err.Error() == "unable to open file" {
				continue
			} else {
				return nil, err
			}
		}
		filesInfo = append(filesInfo, fInfo)
		log.Println(filesInfo)
	}
	return filesInfo, nil
}

// CompareFiles compares files of the same size.
func CompareFiles(input []diskdata.FileInfo) error {

	return nil
}
