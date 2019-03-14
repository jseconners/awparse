package main

import (
	"github.com/jseconners/awparse/util"
	"github.com/urfave/cli"
	"log"
	"os"
)






func main() {
	app := cli.NewApp()

	// CLI sub commands
	app.Commands = []cli.Command{
		{
			Name:    "csv",
			Aliases: []string{"ch"},
			Usage:   "Test command",
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