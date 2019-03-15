package util

import (
	"fmt"
	"github.com/schollz/progressbar"
	"io"
	"log"
)

// a csv file from the data files in the data archive directory matching glob pattern
func MakeCSV(dataDir string, glob string, outputFile string) {
	files := fileList(dataDir, glob)

	if len(files) == 0 {
		fmt.Println(fmt.Sprintf("Didn't find any files in dir: %s, using pattern: %s", dataDir, glob))
		return
	}
	fmt.Println(len(files), " data files found")
	fmt.Println("Checking headers ...")
	headerScores := checkHeaders(&files)
	var skipFiles []string

	if len(headerScores) > 0 {
		fmt.Println("The following files will be ignored because of malformed headers:")
		for _, h := range headerScores {
			skipFiles = append(skipFiles, h.FilePath)
			fmt.Println("\t ", h.FilePath)
		}
	} else {
		fmt.Println("All headers are valid")
	}
	fmt.Print("Continue? (yes or no): ")
	if !Confirm() {
		return
	}

	fmt.Println("\nGenerating CSV file")
	bar := progressbar.New(len(files) - len(headerScores))
	dw := NewDataWriter(outputFile, ',')

	headerWritten := false
	for _, f := range files {

		// skip bad header files
		if len(skipFiles) > 0 && sliceContainsString(skipFiles, f) {
			continue
		}

		// open data file for reading and get header and set offset
		dr := NewDataReader(f, '\t', true, true)
		h, _ := dr.GetHeader(5)

		// replace whitespace in header fields
		for i, field := range h {
			h[i] = replaceWhiteSpace(field, "_")
		}

		// write the output header if not already
		if !headerWritten {
			dw.Writer.Write(h)
			headerWritten = true
		}

		// read data file and write output
		for {
			line, e := dr.Reader.Read()
			if e == io.EOF {
				break
			} else if e != nil {
				log.Fatal(e)
			}
			// trim whitespace from row items and write to output
			trimFields(&line)
			dw.Writer.Write(line)
		}
		dr.Close()
		bar.Add(1)
	}
	dw.Close()
}
