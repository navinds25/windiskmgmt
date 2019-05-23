package app

import (
	"regexp"

	"github.com/navinds25/windiskmgmt/pkg/diskdata"
)

// CheckHighPriorityFolders increments priority for files
// in InboundScriptPull and InboundFromFTPServer
func CheckHighPriorityFolders(input *diskdata.FileInfo) error {
	r1, err := regexp.Compile(".*InboundScriptPull.*")
	if err != nil {
		return err
	}
	r2, err := regexp.Compile(".*InboundFromFTPServer.*")
	if err != nil {
		return err
	}
	reg1match := r1.FindString(input.File)
	if reg1match != "" {
		input.Priority = input.Priority + 1
	}
	reg2match := r2.FindString(input.File)
	if reg2match != "" {
		input.Priority = input.Priority + 1
		input.DoNotDelete = true
	}
	return nil
}

// CheckLowPriorityFiles increments priority for files
// that do not contain "evive_backup" in their name.
func CheckLowPriorityFiles(input *diskdata.FileInfo) error {
	r, err := regexp.Compile(".*_evivebackup_.*")
	if err != nil {
		return err
	}
	regmatch := r.FindString(input.Basename)
	// Please note " == ", as decrementing int 0 will be a problem.
	if regmatch == "" {
		input.DoNotDelete = false
	} else {
		input.Priority = input.Priority + 1
	}
	return nil
}
