package app

import (
	"fmt"
	"strings"

	"github.com/navinds25/windiskmgmt/internal/dfcli"
)

/*
This is for processing the logic for the cli.
*/

// ProcessCli is for processing the CLI arguments.
func ProcessCli() error {
	switch dfcli.CliFlags.Action {
	case "info":
		if dfcli.CliFlags.ListDB {
			if err := ReadDBFiles(); err != nil {
				return err
			}
		}
	case "single-op":
		if err := LoadToDB(dfcli.CliFlags.DFL); err != nil {
			return err
		}
	case "dd":
		if err := ProcessDBFiles(); err != nil {
			return err
		}
	default:
		fmt.Println("select one of the options or -h for help.")
	}
	return nil
}

// GetSkipDirectories returns the directories to be skipped
func GetSkipDirectories() []string {
	defaultDirs := []string{"/proc", "/run", "/sys", "/tmp"}
	if strings.Contains(dfcli.CliFlags.SkipDir, ",") {
		skipdirs := strings.Split(dfcli.CliFlags.SkipDir, ",")
		dirs := defaultDirs
		for _, dir := range skipdirs {
			dirs = append(dirs, dir)
		}
		return dirs
	}
	dirs := append(defaultDirs, dfcli.CliFlags.SkipDir)
	return dirs
}
