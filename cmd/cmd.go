package main

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/dgraph-io/badger/options"

	"github.com/navinds25/windiskmgmt/pkg/diskdata"

	"github.com/dgraph-io/badger"

	"github.com/navinds25/windiskmgmt/internal/app"
	"github.com/navinds25/windiskmgmt/internal/dfcli"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func logSetup() error {
	logfileName := "diskmgmt.log"
	logfile, err := os.OpenFile(logfileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	logwriter := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(logwriter)
	log.SetReportCaller(true)
	customLogFormat := new(logrus.JSONFormatter)
	customLogFormat.PrettyPrint = true
	customLogFormat.TimestampFormat = "2006-01-02 15:04:05"
	log.SetFormatter(customLogFormat)
	//if app.CliVal.Debug {
	//	log.SetLevel(log.DebugLevel)
	//	log.Debug("Debug logs enabled!")
	//}
	return nil
}

func main() {
	// start commandline
	cliapp := dfcli.App()
	err := cliapp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// Read config
	if err := logSetup(); err != nil {
		log.Fatal(err)
	}

	// DB Setup
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	datadir := path.Join(currentDir, "data")
	filesDBOpts := badger.DefaultOptions
	filesDBOpts.Dir = datadir
	filesDBOpts.ValueDir = path.Join(datadir, "value")
	filesDBOpts.Logger = log.New()
	filesDBOpts.Truncate = true
	filesDBOpts.ValueLogLoadingMode = options.FileIO
	filesDB, err := badger.Open(filesDBOpts)
	if err != nil {
		log.Fatal(err)
	}
	diskdata.InitDB(diskdata.BadgerDB{
		DB: filesDB,
	})
	defer diskdata.DataStore.CloseDB()
	log.Debug("Database has been setup")

	// Process Cli
	if err := app.ProcessCli(); err != nil {
		log.Fatal(err)
	}
}
