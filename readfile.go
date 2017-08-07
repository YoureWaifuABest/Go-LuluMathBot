package main

import (
	"os"
)

/* Eventually have this not require an "amount" argument */
func readFileToString(open string, amount int) string {
	file, err := os.Open(open)
	checkErrorPanic(err)
	defer file.Close()

	bytes := make([]byte, amount)
	_, err = file.Read(bytes)
	checkErrorPanic(err)

	return string(bytes)
}
