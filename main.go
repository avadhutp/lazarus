package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	loc := "/tmp/lazarus"
	files, _ := ioutil.ReadDir(loc)

	for _, f := range files {
		if strings.Split(f.Name(), ".")[0] == "3vui82" {
			fmt.Println(f.Name())
		}
	}
}
