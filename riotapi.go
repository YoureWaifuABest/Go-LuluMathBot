package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type item struct {
	Plaintext   string
	Description string
	Name        string
	Id          int
}

var rtoken = readFileToString("riotkey", 42)
var client = &http.Client{}

func stringToID(find string, data []byte) string {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	/*
	 * This feels terribly inefficient, but I don't know of a better way
	 * To do this. Don't understand interfaces well enough
	 */
	/* Defines m as a map[string]interface{} of f using a type assertion */
	m := f.(map[string]interface{})
	/* Defines m1 similarly; this is to, essentially, remove the "data" from the start of the JSON */
	m1 := m["data"].(map[string]interface{})
	/*
	 * Iterates over every element of m1,
	 * returns the key where name is equal to the string we're searching for
	 * This ignores coloq values, currently.
	 * Not hard to have it work with them
	 *
	 * Defines m2, and immediately discards it if it does not equal find.
	 * Would probably be possibly to do this without a type assertion, or defining
	 * m2, but I don't know how.
	 */
	for k, _ := range m1 {
		m2 := m1[k].(map[string]interface{})
		if m2["colloq"] != nil {
			if strings.EqualFold(m2["colloq"].(string), find) {
				return k
			}
		}

		if m2["name"] != nil {
			if strings.EqualFold(m2["name"].(string), find) {
				/* Feels bad storing what is, essentially, an int as a string */
				return k
			}
		}
	}
	return ""
}

func parseData(itemID string, data []byte) (item item) {
	var f interface{}
	err := json.Unmarshal(data, &f)
	if err != nil {
		fmt.Println(err)
	}

	m := f.(map[string]interface{})
	m1 := m["data"].(map[string]interface{})
	itemb, err := json.Marshal(m1[itemID])
	json.Unmarshal(itemb, &item)
	return
}

func getValue(url string) ([]byte, error) {
	var modify bool
	/*
	 * Very basic caching, probably not very efficient
	 * Hard coded file for now, since I'm not sure what I'm going to do with directories
	 * This also assumes the "data" directory is already created,
	 * rather than creating the directory if it's not
	 *
	 * Obviously that's an issue which needs to be fixed at some point
	 */
	info, err := os.Stat("data/items")
	/*
	 * If os.Stat claims the file does not exist, set modify to true
	 * Meaning the file needs to be modified (in this case, created)
	 */
	if os.IsNotExist(err) {
		modify = true
	} else if checkErrorPrint(err) {
		return nil, err
	} else if (time.Now().Sub(info.ModTime())) > (30 * time.Minute) {
		/*
		 * time.Now() returns the current time
		 * info.ModTime() returns when the file from os.Stat was last modified
		 *
		 * time.Now() - info.ModTime() = Time elapsed since file has been modified
		 * If this value is greater than 30 * time.Minute (30 minutes),
		 * set modify to true
		 */
		modify = true
	}

	if modify == true {
		r, err := http.NewRequest("GET", url, nil)
		if checkErrorPrint(err) {
			return nil, err
		}

		r.Header.Add("X-Riot-Token", rtoken)
		resp, err := client.Do(r)
		if checkErrorPrint(err) {
			return nil, err
		}

		if resp.StatusCode != 200 {
			fmt.Println(resp.StatusCode)
			return nil, fmt.Errorf("response code != 200, code: %d", resp.StatusCode)
		}

		body, err := ioutil.ReadAll(resp.Body)
		if checkErrorPrint(err) {
			return nil, err
		}

		/* Dumps body into file */
		err = ioutil.WriteFile("data/items", body, 0644)
		if checkErrorPrint(err) {
			return nil, err
		}

		return body, nil
	}
	/* Spit the whole file out */
	fileContents, err := ioutil.ReadFile("data/items")
	if checkErrorPrint(err) {
		return nil, err
	}
	return fileContents, nil
}
