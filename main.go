package main

import (
	"io"
	"os"

	"github.com/navinds25/windiskmgmt/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	// start commandline
	app := cmd.App()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	// setup logging
	logfile, err := os.OpenFile("diskmgmt.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Error(err)
	}
	defer logfile.Close()
	logmw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(logmw)

	//run the app
	if cmd.Action == "dd" && cmd.StartDir != "" {
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

	}
}
