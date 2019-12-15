package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/craiggwilson/mqlbinary/internal/bsongen"
	"github.com/craiggwilson/mqlbinary/internal/lang"
)

func main() {
	input, err := ioutil.ReadFile("mql.g4.tmpl")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	generator := bsongen.NewBSONGenerator(lang.Java{})

	output, err := generator.Generate(string(input))
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
