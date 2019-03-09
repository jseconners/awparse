package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
)


func readHeader(fileName string) {
	fh, _ := os.Open(fileName)

	s := bufio.NewScanner(fh)
	for s.Scan() {
		fmt.Println(s.Text())
	}

	/**
	reader := csv.NewReader(bufio.NewReader(fh))
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1
	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}
		fmt.Println(line)
	}
	**/
}


func main() {

	app := cli.NewApp()

	app.Commands = []cli.Command{
		{
			Name:    "checkheader",
			Aliases: []string{"ch"},
			Usage:   "Check header definition file",
			Action: func(c *cli.Context) error {
				readHeader(c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}