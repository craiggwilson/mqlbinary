package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/craiggwilson/mqlbinary/internal"
)

func main() {
	input, err := ioutil.ReadFile("mql.g4.tmpl")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output, err := internal.Generate(string(input), internal.JavaLanguage{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = ioutil.WriteFile("mql.g4", []byte(output), 0666)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
