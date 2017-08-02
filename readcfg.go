package main

import (
	"os"
)

/* 
 * For now this just reads the token.
 * A dedicated config is an idea that may or may not happen
 */
func readCFG() string {
	file, err := os.Open("config")
	checkErrorPanic(err)
	defer file.Close()

	bytes := make([]byte, 59)
	_, err = file.Read(bytes)
	checkErrorPanic(err)

	return string(bytes)
}
