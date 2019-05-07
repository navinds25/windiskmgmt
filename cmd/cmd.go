package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/navinds25/windiskmgmt/internal/app"
	"github.com/navinds25/windiskmgmt/internal/dfcli"
	"github.com/navinds25/windiskmgmt/internal/dfconfig"
)

// GetSkipDirectories returns the directories to be skipped
func GetSkipDirectories() []string {
	defaultDirs := []string{"/proc", "/run", "/sys", "/tmp"}
	if strings.Contains(dfcli.SkipDir, ",") {
		skipdirs := strings.Split(dfcli.SkipDir, ",")
		dirs := defaultDirs
		for _, dir := range skipdirs {
			dirs = append(dirs, dir)
		}
		return dirs
	} else {
		dirs := append(defaultDirs, dfcli.SkipDir)
		return dirs
	}
}

// Run main function for the application
func Run() error {
	//	SkipDirectories = GetSkipDirectories()
	duplicateFilesConf := "duplicate_files_2019-04-15.txt"
	dfconf, err := dfconfig.GetConfig(duplicateFilesConf)
	if err != nil {
		return err
	}
	for key, files := range dfconf.Files {
		log.Info("Processing for key:", key)
		//filesInfo, err := app.GetInfoConfFiles(files)
		_, err := app.GetInfoConfFiles(files)
		if err != nil {
			return err
		}
		//log.Println(filesInfo)
	}
	//app.GetInfoConfFiles(dfconf.Files)
	return nil
}
