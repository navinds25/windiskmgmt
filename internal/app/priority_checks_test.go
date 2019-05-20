package app

import (
	"testing"

	"github.com/navinds25/windiskmgmt/pkg/diskdata"
)

// TestCheckHighPriorityFolders checks for regex matching
func TestCheckHighPriorityFolders(t *testing.T) {
	inboundFromFTPServerFile := &diskdata.FileInfo{
		File: "X:\\InboundFromFTPServer\\test\\test_file_inboundfromftp.txt",
	}
	if err := CheckHighPriorityFolders(inboundFromFTPServerFile); err != nil {
		t.Error(err)
	}
	if inboundFromFTPServerFile.Priority != 1 {
		t.Fail()
	}

	inboundScriptPullFile := &diskdata.FileInfo{
		File: "X:\\InboundScriptPull\\test\\test_file_inboundscriptpull.txt",
	}
	if err := CheckHighPriorityFolders(inboundScriptPullFile); err != nil {
		t.Error(err)
	}
	if inboundScriptPullFile.Priority != 1 {
		t.Fail()
	}
}

// TestCheckLowPriorityFiles checks for regex for evive_backups
func TestCheckLowPriorityFiles(t *testing.T) {
	eviveBackupsFile := &diskdata.FileInfo{
		File:     "X:\\InboundFromFTPServer\\cae_tri_ad\\test\\JustTesting.txt__evivebackup_20161118_131002",
		Basename: "JustTesting.txt__evivebackup_20161118_131002",
	}
	if err := CheckLowPriorityFiles(eviveBackupsFile); err != nil {
		t.Error(err)
	}
	if eviveBackupsFile.Priority != 0 {
		t.Fail()
	}
	notEviveBackupsFile := &diskdata.FileInfo{
		Basename: "The_should_increment.txt",
	}
	if err := CheckLowPriorityFiles(notEviveBackupsFile); err != nil {
		t.Error(err)
	}
	if notEviveBackupsFile.Priority != 1 {
		t.Fail()
	}
}
