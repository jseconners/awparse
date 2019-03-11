package main

import (
	"bufio"
	"encoding/csv"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)


var months = []string {
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

func fileExists(pathStr string) bool {
	_, err := os.Stat(pathStr);
	return !os.IsNotExist(err)
}

func dataFileDay(fileName string) int {
	dayInt, _ := strconv.Atoi(fileName[6:8])
	return dayInt
}

func makeCSV(dataDir string) {
	yearDirs, _ := filepath.Glob(filepath.Join(dataDir, "/[1-2][0-9][0-9][0-9]"))
	for _, yearDir := range yearDirs {
		for _, month := range months {
			monthDir := filepath.Join(yearDir, month)
			if fileExists(monthDir) {
				dataFiles, _ := filepath.Glob(filepath.Join(monthDir, "*.100"))

				// Ensure we're parsing these files in order
				fileMap := make(map[int]string)
				dayInts := make([]int, 0, 30)
				for _, df := range dataFiles {
					di := dataFileDay(filepath.Base(df))
					dayInts = append(dayInts, di)
					fileMap[di] = df
				}
				sort.Ints(dayInts)

				for _, di := range dayInts {
					dataFile := fileMap[di]
					fileHandle, _ := os.Open(dataFile)
					reader := csv.NewReader(bufio.NewReader(fileHandle))
					reader.Comma = '\t'
					reader.FieldsPerRecord = -1

					writer := csv.NewWriter(os.Stdout)
					writer.Comma = ','

					for {
						line, error := reader.Read()
						if error == io.EOF {
							break
						} else if error != nil {
							log.Fatal(error)
						}
						writer.Write(line)
					}
					fileHandle.Close()
					writer.Flush()
				}
			}
		}
	}

}



func main() {
	app := cli.NewApp()
	app.Commands = []cli.Command{
		{
			Name:    "csv",
			Aliases: []string{"ch"},
			Usage:   "Test command",
			Action: func(c *cli.Context) error {
				makeCSV(c.Args().First())
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}