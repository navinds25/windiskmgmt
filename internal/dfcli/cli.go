package dfcli

import "github.com/urfave/cli"

// Dryrun for controlling dryrun flag & operation
var Dryrun bool

// Action for command type
var Action string

// StartDir root dir for searching
var StartDir string

// DelDir destination directory for files to be deleted.
var DelDir string

// SkipDir contains the commandline directories to be skipped.
var SkipDir string

// SkipDirectories is a string slice containing all the directories to be skipped.
var SkipDirectories []string

// Debug for controlling debug flag & operation
var Debug bool

// deleteDuplicates is the command object for the cli
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
		cli.StringFlag{
			Name:        "skipdir",
			Usage:       "directories to be skipped",
			Destination: &SkipDir,
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
