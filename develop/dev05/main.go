package main

import (
	"fmt"
)

func main() {
	app := App{}

	if err := app.init(); err != nil {
		fmt.Printf("error on parse flags: %v", err)
		return
	}

	result, err := getFileStrings(app.filename)
	if err != nil {
		fmt.Printf("error on reading file: %v", err)
		return
	}

	rw := app.formatResult(result)

	app.printResults(rw)
}

func (a *App) formatResult(s []string) []StringResult {
	return a.formatter.Format(s, a.contextCount)
}

func (a *App) printResults(s []StringResult) {
	n := a.numberPrint
	fmt.Printf("Result (count: %v):\n", len(s))
	for _, v := range s {
		if n {
			fmt.Printf("%v: %v\n", v.id, v.str)
		} else {
			fmt.Println(v.str)
		}
	}
}
