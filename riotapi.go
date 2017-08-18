package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
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
	r, err := http.NewRequest("GET", url, nil)
	if checkErrorPrint(err) {
		return
	}

	r.Header.Add("X-Riot-Token", rtoken)
	resp, err := client.Do(r)
	if checkErrorPrint(err) {
		return
	}

	if resp.StatusCode != 200 {
		fmt.Println(resp.StatusCode)
		return nil, fmt.Errorf("response code != 200, code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body, nil
}
