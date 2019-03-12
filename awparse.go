package main

import (
	"bufio"
	"encoding/csv"
	"github.com/schollz/progressbar"
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
 * Get a sorted slice of data file paths from a specified
 * month's data directory
 */
func getSortedMonthFiles(dirName string, glob string) []string {
	dataFiles, _ := filepath.Glob(filepath.Join(dirName, glob))
	fileMap := make(map[int]string)
	dayInts := make([]int, 0, 31)

	for _, df := range dataFiles {
		di := dataFileDay(filepath.Base(df))
		dayInts = append(dayInts, di)
		fileMap[di] = df
	}
	sort.Ints(dayInts)

	sortedDataFiles := make([]string, 0, 31)
	for _, di := range dayInts {
		sortedDataFiles = append(sortedDataFiles, fileMap[di])
	}
	return sortedDataFiles
}

/**
 * Get a sorted slice of all data file paths for a given data
 * archive directory and file glob pattern
 */
func getAllSortedDataFiles(dirName string, glob string) []string {
	allDataFiles := make([]string, 0, 100)
	yearDirs, _ := filepath.Glob(filepath.Join(dirName, "/[1-2][0-9][0-9][0-9]"))
	for _, yearDir := range yearDirs {
		for _, month := range months {
			monthDir := filepath.Join(yearDir, month)
			if fileExists(monthDir) {
				allDataFiles = append(allDataFiles, getSortedMonthFiles(monthDir, glob)...)
			}
		}
	}
	return allDataFiles
}


/**
 * Parse data files and write generated CSV to stdout
 */
func makeCSV(dataDir string, glob string) {
	allDataFiles := getAllSortedDataFiles(dataDir, glob)
	bar := progressbar.New(len(allDataFiles))

	outFile, _ := os.Create("data.csv")
	writer := csv.NewWriter(outFile)
	writer.Comma = ','

	for _, df := range allDataFiles {
		bar.Add(1)

		inFile, _ := os.Open(df)
		reader := csv.NewReader(bufio.NewReader(inFile))
		reader.Comma = '\t'
		reader.FieldsPerRecord = -1

		for {
			line, e := reader.Read()
			if e == io.EOF {
				break
			} else if e != nil {
				log.Fatal(e)
			}
			writer.Write(line)
		}
		inFile.Close()
	}
	writer.Flush()
	outFile.Close()
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
				makeCSV(c.Args().Get(0), c.Args().Get(1))
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