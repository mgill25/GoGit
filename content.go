package main

import (
	"io/ioutil"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readContentFile(contentPath string) string {
	data, err := ioutil.ReadFile(contentPath)
	check(err)
	return BytesToStr(data)
}

func prependContentHeaders(contentType string, content string) string {
	header := contentType + " " + strconv.Itoa(len(content)) + "\000"
	newContent := header + content
	return newContent
}
