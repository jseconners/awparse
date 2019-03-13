package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"github.com/araddon/dateparse"
	"github.com/schollz/progressbar"
	"github.com/urfave/cli"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)


type CSVW struct {
	File *os.File
	Writer *csv.Writer
}

func NewCSVW(filePath string) *CSVW {
	csvw := new(CSVW)
	fh, _ := os.Create(filePath)
	csvw.File = fh
	csvw.Writer = csv.NewWriter(csvw.File)
	csvw.Writer.Comma = ','
	return csvw
}

func (csvw *CSVW) Close() {
	csvw.Writer.Flush()
	csvw.File.Close()
}


type CSVR struct {
	FilePath string
	File *os.File
	Reader *csv.Reader
}

func NewCSVR(filePath string) *CSVR {
	csvr := new(CSVR)
	csvr.FilePath = filePath
	csvr.File, _ = os.Open(filePath)
	csvr.SetReader()
	return csvr
}

func (c *CSVR) SetReader() {
	c.Reader = csv.NewReader(bufio.NewReader(c.File))
	c.Reader.Comma = '\t'
	c.Reader.FieldsPerRecord = -1
}

func (c *CSVR) Close() {
	c.File.Close()
}

func (c *CSVR) Rewind() {
	c.File.Seek(0, 0)
	c.SetReader()
}


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
 * Check if a row of values contains a date
 */
func containsDate(row []string) bool {
	for _, val := range row {
		_, e := dateparse.ParseAny(val)
		if e == nil {
			return true
		}
	}
	return false
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
 * Detect and return the header row for a CSVR and set its read offset
 * so the next Read() returns the first data record. Or return an empty
 * slice and error if header row couldn't be detected
 */
func getHeader(csvr *CSVR) ([]string, error) {
	max_depth := 5
	row_index := 0

	var rows [][]string
	for {
		line, e := csvr.Reader.Read()
		if e != nil || e == io.EOF || row_index == max_depth {
			return make([]string, 0), errors.New("could not detect header")
		}
		if containsDate(line) && len(rows) > 0 {
			// Rewind and advance to next read is from first data record
			csvr.Rewind()
			for i := 0; i < row_index; i++ {
				csvr.Reader.Read()
			}
			return rows[row_index - 1], nil
		} else {
			rows = append(rows, line)
		}
		row_index += 1
	}
}

func getHeaders(dataFiles *[]string) [][]string {
	var headers [][]string
	for _, df := range *dataFiles {
		csvr := NewCSVR(df)
		header, err := getHeader(csvr)
		csvr.Close()

		if err != nil {
			log.Fatal(err)
		}
		headers = append(headers, header)
	}
	return headers
}


func checkHeaders(dataFiles *[]string) bool {
	headers := getHeaders(dataFiles)
	scores := make([][]int, len(headers))

	for i := 0; i < len(headers); i++ {
		score := 0
		for j := 0; j < len(headers); j++ {
			if (i==j) {
				continue
			}
			if (headers[i] != headers[j]) {
				score += 1
			}
		}
		scores[i] = score
	}
	return false
}


/**
 * Parse data files and write generated CSV to stdout
 */
func makeCSV(dataDir string, glob string) {
	allDataFiles := getAllSortedDataFiles(dataDir, glob)

	headerCheck := checkHeaders(&allDataFiles)
	return


	bar := progressbar.New(len(allDataFiles))


	csvw := NewCSVW("data.csv")
	for _, df := range allDataFiles {
		bar.Add(1)

		csvr := NewCSVR(df)
		for {
			line, e := csvr.Reader.Read()
			if e == io.EOF {
				break
			} else if e != nil {
				log.Fatal(e)
			}
			csvw.Writer.Write(line)
		}
		csvr.Close()
	}
	csvw.Close()
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