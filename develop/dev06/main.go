package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"
)

type StringFlag struct {
	isSet bool
	val   string
}

func (s *StringFlag) String() string {
	return s.val
}

func (s *StringFlag) Set(v string) error {
	s.val = v
	s.isSet = true
	return nil
}

type App struct {
	fields    []int
	delimeter string
	separated bool
}

func main() {
	app := App{}

	if err := app.LoadFlags(); err != nil {
		fmt.Printf("Error on load flags: %v", err)
		return
	}

	result := app.ReadCLI()
	table := app.SplitToTable(result)

	fmt.Printf("Result (%v rows):\n", len(table))

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 1, ' ', 0)

	for _, v := range table {
		result := strings.Join(v, "\t")
		fmt.Fprintln(w, result)
	}

	w.Flush()

}

func (a *App) SplitToTable(s []string) (result [][]string) {
	d := a.delimeter
	sep := a.separated

	for _, v := range s {
		r := strings.Split(v, d)
		if sep == true && len(r) == 1 {
			continue
		}
		r = a.FilterFields(r)
		result = append(result, r)
	}
	return
}

func (a *App) FilterFields(s []string) (result []string) {
	if len(a.fields) == 0 {
		result = s
		return
	}

	for _, v := range a.fields {
		if v >= len(s) {
			continue
		}
		result = append(result, s[v])
	}

	return
}

func (a *App) ReadCLI() (result []string) {
	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			if err == io.EOF {
				return result
			}
		}

		str := strings.TrimLeft(string(line), " ")
		str = strings.TrimRight(str, " ")
		if len(str) == 0 {
			continue
		}

		result = append(result, str)
	}
}

func (a *App) LoadFlags() error {
	var f StringFlag
	var d StringFlag
	var s bool

	flag.Var(&f, "f", "выбрать поля (колонки)")
	flag.Var(&d, "d", "использовать другой разделитель")
	flag.BoolVar(&s, "s", false, "только строки с разделителем")
	flag.Parse()

	if err := a.LoadFields(f); err != nil {
		return err
	}
	if err := a.LoadDelimeter(d); err != nil {
		return err
	}
	a.LoadSeparator(s)

	return nil
}

func (a *App) LoadFields(s StringFlag) error {
	if !s.isSet {
		return nil
	}

	value := s.val

	if len(value) == 0 {
		return fmt.Errorf("fields are empty")
	}

	fields := strings.Split(value, ",")
	a.fields = make([]int, 0)

	for _, v := range fields {
		field, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("error while parsing fields: %v", err)
		}
		if field < 0 {
			return fmt.Errorf("field %v is not correct", field)
		}

		a.fields = append(a.fields, field)
	}

	return nil
}

func (a *App) LoadDelimeter(s StringFlag) error {
	if !s.isSet {
		a.delimeter = "\t"
		return nil
	}

	value := s.val

	if len(value) == 0 {
		return fmt.Errorf("delimeter is empty")
	}

	a.delimeter = value
	return nil
}

func (a *App) LoadSeparator(s bool) {
	a.separated = s
}
