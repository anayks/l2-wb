package main

import (
	"flag"
	"fmt"
	"math"
	"strconv"
)

type intFlag struct {
	IsSet bool
	Val   int
}

func (i *intFlag) String() string {
	if !i.IsSet {
		return "<not set>"
	}
	return strconv.Itoa(i.Val)
}

func (i *intFlag) Set(value string) error {
	v, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	i.IsSet = true
	i.Val = v
	return nil
}

type stringFlag struct {
	set   bool
	value string
}

func (sf *stringFlag) Set(x string) error {
	sf.value = x
	sf.set = true
	return nil
}

func (sf *stringFlag) String() string {
	return sf.value
}

func (app *App) loadPrintFlag(n *bool) {
	if *n {
		app.numberPrint = true
	}
}

func (app *App) loadContextType(a *bool, b *bool, c *bool, v *bool) {
	if *a == true {
		app.contextIntType = SearchAfter
		app.formatter = FormatAfter{
			app: app,
		}
	} else if *b == true {
		app.contextIntType = SearchBefore
		app.formatter = FormatBefore{
			app: app,
		}
	} else if *c == true {
		app.contextIntType = SearchContext
		app.formatter = FormatContext{
			app: app,
		}
	} else if *v == true {
		app.contextIntType = SearchInvert
		app.formatter = FormatInverted{
			app: app,
		}
	} else {
		app.contextIntType = SearchDefault
		app.formatter = FormatDefault{
			app: app,
		}
	}
}

func (app *App) loadCountOfStrings(c intFlag) error {
	if !c.IsSet {
		if app.contextIntType != SearchContext {
			app.contextCount = math.MaxInt
			return nil
		}
		return fmt.Errorf("Not entered count: -c")
	} else if c.Val < 1 {
		return fmt.Errorf("Unexpected count: %d", c.Val)
	}

	app.contextCount = c.Val

	return nil
}

func (app *App) loadIgnoreCase(i *bool) {
	if *i {
		app.ignoreCase = true
	}
}

func (app *App) loadFixedFlag(F *bool) {
	if *F {
		app.fixedPattern = true
	}
}

func (app *App) getContext() error {
	flag.Parse()
	args := flag.Args()
	if len(args) < 2 {
		return fmt.Errorf("filename or search not writed at flag")
	}
	app.filename = args[0]
	app.search = args[1]
	return nil
}
