package util

// https://gist.github.com/albrow/5882501

import (
	"fmt"
	"log"
)

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func Confirm() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	yResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nResponses := []string{"n", "N", "no", "No", "NO"}

	if sliceContainsString(yResponses, response) {
		return true
	} else if sliceContainsString(nResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return Confirm()
	}
}


