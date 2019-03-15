package util

import (
	"github.com/araddon/dateparse"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// regex for string spaces
var spaceRegex = regexp.MustCompile(`\s+`)

// ordered month full names
var months = []string {
	"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

// csv header check score from being compared
// to all other headers in archive
type HeaderScore struct {
	FilePath string
	Score int
}

// get all headers from list of file paths
func fileHeaders(dataFiles *[]string) [][]string {
	var headers [][]string
	for _, df := range *dataFiles {
		dr := NewDataReader(df, '\t', true, true)
		header, err := dr.GetHeader(5)
		dr.Close()

		if err != nil {
			log.Fatal(err)
		}
		headers = append(headers, header)
	}
	return headers
}

// Get a sorted slice of all data file paths for a given data
// archive directory and file glob pattern
func fileList(dirName string, glob string) []string {
	files := make([]string, 0, 100)
	yearDirs, _ := filepath.Glob(filepath.Join(dirName, "/[1-2][0-9][0-9][0-9]"))
	for _, yearDir := range yearDirs {
		for _, month := range months {
			monthDir := filepath.Join(yearDir, month)
			if fileExists(monthDir) {
				files = append(files, monthFiles(monthDir, glob)...)
			}
		}
	}
	return files
}

// Get a sorted list of data files from a specific month's directory
func monthFiles(dirName string, glob string) []string {
	files, _ := filepath.Glob(filepath.Join(dirName, glob))
	fileMap := make(map[int]string)
	dayInts := make([]int, 0, 31)

	for _, f := range files {
		day := dataFileDay(filepath.Base(f))
		dayInts = append(dayInts, day)
		fileMap[day] = f
	}
	sort.Ints(dayInts)

	sortedFiles := make([]string, 0, 31)
	for _, day := range dayInts {
		sortedFiles = append(sortedFiles, fileMap[day])
	}
	return sortedFiles
}

// check data file headers and return list of scores
func checkHeaders(dataFiles *[]string) []HeaderScore {
	headers := fileHeaders(dataFiles)
	var allScores []HeaderScore
	var badScores []HeaderScore
	var intScores []int

	for i := 0; i < len(headers); i++ {
		score := 0
		for j := 0; j < len(headers); j++ {
			if  i==j {
				continue
			}
			if !headersEqual(headers[i], headers[j]) {
				score += 1
			}
		}
		intScores = append(intScores, score)
		allScores = append(allScores, HeaderScore{(*dataFiles)[i], score})
	}
	// get and return bad scoring headers
	sort.Ints(intScores)
	minScore := intScores[0]
	for _, s := range allScores {
		if s.Score > minScore {
			badScores = append(badScores, s)
		}
	}
	return badScores
}


// Check if file or directory at pathStr exists
func fileExists(pathStr string) bool {
	_, err := os.Stat(pathStr);
	return !os.IsNotExist(err)
}

// Get day of the month from the data file name
// Note: this should probably made more flexible
func dataFileDay(fileName string) int {
	dayInt, _ := strconv.Atoi(fileName[6:8])
	return dayInt
}

// Check if a row of values contains a date
func containsDate(row []string) bool {
	for _, val := range row {
		_, e := dateparse.ParseAny(val)
		if e == nil {
			return true
		}
	}
	return false
}

// check if two csv headers are equal
func headersEqual(h1, h2 []string) bool {
	if len(h1) != len(h2) {
		return false
	}
	for i := 0; i < len(h1); i++ {
		if h1[i] != h2[i] {
			return false
		}
	}
	return true
}

// string in string slice check
func sliceContainsString(items []string, item string) bool {
	for _, s := range items {
		if item == s {
			return true
		}
	}
	return false
}

// trim and replace whitespace with underscore
func replaceWhiteSpace(item string, repl string) string {
	return spaceRegex.ReplaceAllLiteralString(strings.TrimSpace(item), repl)
}

// trim whitespace from string items in place
func trimFields(items *[]string) {
	for i, item := range *items {
		(*items)[i] = strings.TrimSpace(item)
	}
}