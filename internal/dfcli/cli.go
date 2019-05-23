package dfcli

import "github.com/urfave/cli"

// CliFlagsStruct is the struct for the commandline flags.
type CliFlagsStruct struct {
	Dryrun          bool
	Debug           bool
	Action          string
	StartDir        string
	DelDir          string
	SkipDir         string
	SkipDirectories []string
	ListDB          bool
	DFL             string
	NoDB            bool
}

// CliFlags is the instance of the cli flags.
var CliFlags CliFlagsStruct

// deleteDuplicates is the command object for the cli
var deleteDuplicates = cli.Command{
	Name:    "dd",
	Aliases: []string{"delete_duplicates"},
	Usage:   "Move duplicate files to a delete folder.",
	Action: func(c *cli.Context) error {
		CliFlags.Action = "dd"
		return nil
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logging",
			Destination: &CliFlags.Debug,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "Disable/Enable dryrun",
			Destination: &CliFlags.Dryrun,
		},
		cli.StringFlag{
			Name:        "startdir",
			Usage:       "directory to start from",
			Destination: &CliFlags.StartDir,
		},
		cli.StringFlag{
			Name:        "deldir",
			Usage:       "directory to collect files to be deleted",
			Destination: &CliFlags.DelDir,
		},
		cli.StringFlag{
			Name:        "skipdir",
			Usage:       "directories to be skipped",
			Destination: &CliFlags.SkipDir,
		},
	},
}

var infoCommand = cli.Command{
	Name:    "info",
	Aliases: []string{"i"},
	Usage:   "info",
	Action: func(c *cli.Context) error {
		CliFlags.Action = "info"
		return nil
	},
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:        "list-db",
			Usage:       "list contents of the database",
			Destination: &CliFlags.ListDB,
		},
	},
}

var singleOpCommand = cli.Command{
	Name:    "single-op",
	Aliases: []string{"op"},
	Usage:   "run a single operation, instead of the entire process.",
	Action: func(c *cli.Context) error {
		CliFlags.Action = "single-op"
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "dfl",
			Usage:       "duplicate-files-list: input file containing list of duplicate files.",
			Destination: &CliFlags.DFL,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logging",
			Destination: &CliFlags.Debug,
		},
	},
}

var processConfCommand = cli.Command{
	Name:    "process-conf",
	Aliases: []string{"pc"},
	Usage:   "process files from conf in memory.",
	Action: func(c *cli.Context) error {
		CliFlags.Action = "process-conf"
		CliFlags.NoDB = true
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:        "dfl",
			Usage:       "duplicate-files-list: input file containing list of duplicate files.",
			Destination: &CliFlags.DFL,
		},
		cli.BoolFlag{
			Name:        "dryrun",
			Usage:       "Disable/Enable dryrun",
			Destination: &CliFlags.Dryrun,
		},
		cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug logging",
			Destination: &CliFlags.Debug,
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
		infoCommand,
		singleOpCommand,
		processConfCommand,
	}
	return app
}
