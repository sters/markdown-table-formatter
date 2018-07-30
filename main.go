package main

import (
	"fmt"
	"io/ioutil"
	"os"

	markdowntableformatter "github.com/sters/markdown-table-formatter/formatter"
)

func main() {
	body, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	result := markdowntableformatter.Execute(string(body))
	fmt.Print(result)
}
