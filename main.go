package main

import (
	"flag"
	"io/ioutil"
	"os"

	markdowntableformatter "github.com/sters/markdown-table-formatter/formatter"
)

type config struct {
	input  *os.File
	output *os.File
}

func parseArgs() (*config, error) {
	config := &config{
		input:  os.Stdin,
		output: os.Stdout,
	}

	var (
		input  = flag.String("i", "", "input file, default = stdin")
		output = flag.String("o", "", "output file, default = stdout")
		err    error
	)
	flag.Parse()

	if *input != "" {
		config.input, err = os.OpenFile(*input, os.O_RDWR, 0666)
		if err != nil {
			return nil, err
		}
	}

	if *output != "" {
		config.output, err = os.OpenFile(*output, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func main() {
	config, err := parseArgs()
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(config.input)
	if err != nil {
		panic(err)
	}

	result := markdowntableformatter.Execute(string(body))
	if _, err := config.output.WriteString(result); err != nil {
		panic(err)
	}

	config.input.Close()
	config.output.Close()
}
