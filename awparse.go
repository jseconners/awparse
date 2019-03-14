package main

import (
	"github.com/jseconners/awparse/util"
	"github.com/urfave/cli"
	"log"
	"os"
)

// main functions runs a CLI app with sub commands
// for parsing weather data archives
func main() {
	app := cli.NewApp()
	app.Name = "awparse"
	app.Usage = "Parse AWS weather data from AMRC"

	// CLI sub commands
	app.Commands = []cli.Command{
		{
			Name:    "build-csv",
			Aliases: []string{"csv"},
			Usage:
				"Specify data directory, glob pattern and output file name to " +
				"generate a compiled CSV data file.",
			Action: func(c *cli.Context) error {
				util.MakeCSV(c.Args().Get(0), c.Args().Get(1), c.Args().Get(2))
				return nil
			},
		},
	}
	// Run CLI app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}