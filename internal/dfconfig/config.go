package dfconfig

import (
	"encoding/json"
	"errors"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/ghodss/yaml"
	"github.com/navinds25/windiskmgmt/pkg/diskdata"
)

// Config is for the config file from duplicate files finder
type Config struct {
	Files map[string][]string
}

// GetConfig returns the configuration
func GetConfig(configFile string) (*Config, error) {
	var config Config
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func getCheckSum(file io.Reader) (uint32, error) {
	crc32q := crc32.MakeTable(0xD5828281)
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return 0, nil
	}
	checksum := crc32.Checksum(data, crc32q)
	return checksum, nil
}

// GetFileInfo returns FileInfo struct with info on files
func GetFileInfo(filename string) (diskdata.FileInfo, error) {
	theFile, err := os.Open(filename)
	if err != nil {
		return diskdata.FileInfo{}, errors.New("unable to open file")
	}
	defer theFile.Close()
	fileDetails, err := theFile.Stat()
	if err != nil {
		return diskdata.FileInfo{}, err
	}
	fileSize := fileDetails.Size()
	checksum, err := getCheckSum(theFile)
	if err != nil {
		return diskdata.FileInfo{}, err
	}
	file := diskdata.FileInfo{
		File:     filename,
		Basename: path.Base(filename),
		Checksum: checksum,
		Size:     fileSize,
		Modified: fileDetails.ModTime(),
	}
	return file, err
}
