//
// structures for reading and writing weather data files
// wrapping the encoding/csv package
//

package util

import (
	"bufio"
	"encoding/csv"
	"errors"
	"io"
	"os"
)

// csv data file write wrapper
type DataWriter struct {
	File *os.File
	Writer *csv.Writer
}

// DataWriter constructor
func NewDataWriter(filePath string, comma rune) *DataWriter {
	dw := new(DataWriter)
	fh, _ := os.Create(filePath)
	dw.File = fh
	dw.Writer = csv.NewWriter(dw.File)
	dw.Writer.Comma = comma
	return dw
}

// close DataWriter, flushing buffer
func (dw *DataWriter) Close() {
	dw.Writer.Flush()
	dw.File.Close()
}

// csv data file read wrapper
type DataReader struct {
	FilePath string
	File *os.File
	Reader *csv.Reader
	Comma rune
	VarLen bool
	Trim bool
}

// DataReader constructor
func NewDataReader(filePath string, comma rune, varLen bool, trim bool) *DataReader {
	dr := new(DataReader)
	dr.FilePath = filePath
	dr.File, _ = os.Open(filePath)
	dr.Comma = comma
	dr.VarLen = varLen
	dr.Trim = trim
	dr.SetReader()
	return dr
}

// initialize the csv.Reader property for a DataReader
func (dr *DataReader) SetReader() {
	dr.Reader = csv.NewReader(bufio.NewReader(dr.File))
	dr.Reader.Comma = dr.Comma
	if dr.VarLen {
		dr.Reader.FieldsPerRecord = -1
	}
	if dr.Trim {
		dr.Reader.TrimLeadingSpace = true
	}
}

// close a DataReader
func (dr *DataReader) Close() {
	dr.File.Close()
}

// set the record offset for a DataReader
func (dr *DataReader) SetRecordOffset(offset int) {
	dr.File.Seek(0, 0)
	dr.SetReader()
	for i := 0; i < offset; i++ {
		dr.Reader.Read()
	}
}

// detect and return header row for this DataReader and set the
// record offset so next read is from the first data row
func (dr *DataReader) GetHeader(maxDepth int) ([]string, error) {
	// reset reader to beginning of file
	dr.SetRecordOffset(0)

	rowIndex := 0
	var rows [][]string
	for {
		line, e := dr.Reader.Read()
		if e != nil || e == io.EOF || rowIndex == maxDepth {
			return make([]string, 0), errors.New("could not detect header")
		}
		if containsDate(line) && len(rows) > 0 {
			// Rewind and advance to next read is from first data record
			dr.SetRecordOffset(rowIndex)
			return rows[rowIndex - 1], nil
		} else {
			rows = append(rows, line)
		}
		rowIndex += 1
	}
}


