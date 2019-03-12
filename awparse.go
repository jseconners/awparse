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


// Ordered month full name slice
var months = []string {
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

/*
 * Check if file or directory at pathStr exists
 */
func fileExists(pathStr string) bool {
	_, err := os.Stat(pathStr);
	return !os.IsNotExist(err)
}

/*
 * Get day of the month from the data file name
 * Note: this should probably made more flexible
 */
func dataFileDay(fileName string) int {
	dayInt, _ := strconv.Atoi(fileName[6:8])
	return dayInt
}

/**
 * Get slice of data files sorted by dataFileDay() return
 */
func getSortedDataFiles(dirName string, glob string) []string {
	dataFiles, _ := filepath.Glob(filepath.Join(dirName, glob))
	fileMap := make(map[int]string)
	dayInts := make([]int, 0, 30)

	for _, df := range dataFiles {
		di := dataFileDay(filepath.Base(df))
		dayInts = append(dayInts, di)
		fileMap[di] = df
	}
	sort.Ints(dayInts)

	sortedDataFiles := make([]string, 0, 30)
	for _, di := range dayInts {
		sortedDataFiles = append(sortedDataFiles, fileMap[di])
	}
	return sortedDataFiles
}


/**
 * Parse data files and write generated CSV to stdout
 */
func makeCSV(dataDir string) {
	yearDirs, _ := filepath.Glob(filepath.Join(dataDir, "/[1-2][0-9][0-9][0-9]"))
	for _, yearDir := range yearDirs {
		for _, month := range months {
			monthDir := filepath.Join(yearDir, month)
			if fileExists(monthDir) {

				sortedDataFiles := getSortedDataFiles(monthDir, "*100")

				for _, df := range sortedDataFiles {
					fileHandle, _ := os.Open(df)
					reader := csv.NewReader(bufio.NewReader(fileHandle))
					reader.Comma = '\t'
					reader.FieldsPerRecord = -1

					writer := csv.NewWriter(os.Stdout)
					writer.Comma = ','

					for {
						line, e := reader.Read()
						if e == io.EOF {
							break
						} else if e != nil {
							log.Fatal(e)
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

	// CLI sub commands
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

	// Run CLI app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}