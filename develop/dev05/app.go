package main

import (
	"flag"
)

type SearchType int

const (
	SearchDefault = SearchType(iota)
	SearchAfter
	SearchBefore
	SearchContext
	SearchInvert
)

type App struct {
	filename       string
	search         string
	contextType    SearchType
	contextIntType SearchType
	contextCount   int
	ignoreCase     bool
	invert         bool
	fixedPattern   bool
	numberPrint    bool

	formatter Formatter
}

func (app *App) init() error {
	var a, b, c, v, i, f, n bool
	var count intFlag

	flag.BoolVar(&a, "A", false, "after string")
	flag.BoolVar(&b, "B", false, "before string")
	flag.BoolVar(&c, "C", false, "context with strings")
	flag.BoolVar(&v, "v", false, "inverted result")
	flag.BoolVar(&i, "i", false, "ignore case")
	flag.Var(&count, "c", "count of strings")
	flag.BoolVar(&f, "F", false, "fixed")
	flag.BoolVar(&n, "n", false, "print count")

	flag.Parse()

	app.loadContextType(&a, &b, &c, &v)

	if err := app.loadCountOfStrings(count); err != nil {
		return err
	}

	app.loadIgnoreCase(&i)
	app.loadFixedFlag(&f)
	app.loadPrintFlag(&n)

	if err := app.getContext(); err != nil {
		return err
	}

	return nil
}
