package main

import (
	"strings"
)

type Formatter interface {
	Format([]string, int) []StringResult
}

type FormatDefault struct {
	app *App
}

type StringResult struct {
	id  int
	str string
}

func (f FormatDefault) Format(s []string, count int) []StringResult {
	var result = []StringResult{}
	for id, v := range s {
		if len(result) >= count {
			break
		}
		if SearchEntry(v, f.app) {
			result = append(result, StringResult{
				id:  id,
				str: v,
			})
		}
	}
	return result
}

type FormatBefore struct {
	app *App
}

func (f FormatBefore) Format(s []string, count int) []StringResult {
	var result = []StringResult{}
	var lastEnter int
	for id, v := range s {
		if SearchEntry(v, f.app) {
			lastEnter = id
		}
	}
	start := lastEnter - count
	var strResult []string
	if start < 0 {
		start = 0
	}
	strResult = s[start:lastEnter]
	for id, v := range strResult {
		result = append(result, StringResult{
			id:  id + start,
			str: v,
		})
	}
	return result
}

type FormatAfter struct {
	app *App
}

func (f FormatAfter) Format(s []string, count int) []StringResult {
	var result = []StringResult{}
	var firstEnter int = -1
	for id, v := range s {
		if SearchEntry(v, f.app) {
			firstEnter = id
		}
	}
	if firstEnter == -1 {
		return []StringResult{}
	}
	start := firstEnter + 1
	if start >= len(s) {
		return []StringResult{}
	}
	end := firstEnter + count + 1
	if end > len(s) {
		end = len(s)
	}

	strResult := s[start:end]
	for id, v := range strResult {
		result = append(result, StringResult{
			id:  firstEnter + id,
			str: v,
		})
	}
	return result
}

type FormatContext struct {
	app *App
}

func (f FormatContext) Format(s []string, count int) []StringResult {
	var result = []StringResult{}
	var firstEnter int = -1
	for id, v := range s {
		if SearchEntry(v, f.app) {
			firstEnter = id
			continue
		}
	}
	if firstEnter == -1 {
		return []StringResult{}
	}

	start := firstEnter - count
	if start < 0 {
		start = 0
	}
	end := firstEnter + count + 1
	if end > len(s) {
		end = len(s)
	}

	strResult := s[start:end]
	for id, v := range strResult {
		result = append(result, StringResult{
			id:  start + id,
			str: v,
		})
	}
	return result
}

type FormatInverted struct {
	app *App
}

func (f FormatInverted) Format(s []string, count int) []StringResult {
	var result = []StringResult{}
	for id, v := range s {
		if len(result) >= count {
			break
		}
		if SearchEntry(v, f.app) {
			continue
		}
		result = append(result, StringResult{
			id:  id,
			str: v,
		})
	}
	return result
}

func SearchEntry(s string, a *App) bool {
	fixed := a.fixedPattern
	search := a.search
	ignore := a.ignoreCase

	if ignore {
		search = strings.ToLower(search)
		s = strings.ToLower(s)
	}

	if fixed {
		return search == s
	} else {
		return strings.Contains(s, search)
	}
}
