package main

import (
	"io/ioutil"
	"log"
)

func readfile(filename string) string {

	content, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	return string(content)
}
