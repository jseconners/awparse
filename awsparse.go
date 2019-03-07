package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {

	// available commands
	commands := map[string]int{
		"check-header": 1,
	}
	command := os.Args[1]
	commandVal, exists := commands[command]

	if commandVal!=1 || !exists {
		fmt.Printf("Command %s is not available!\n", command)
		os.Exit(1)
	}

	fmt.Printf("Executing command: %s ...\n", command)
	os.Exit(0)

	dataFile := "something"
	csvFile, _ := os.Open(dataFile)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = '\t'

	for {
		line, e := reader.Read()
		if e == io.EOF {
			break
		} else if e != nil {
			log.Fatal(e)
		}
		fmt.Println(line)
	}
}