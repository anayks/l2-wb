package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"dev09_01/internal/commons"
	"dev09_01/internal/download"
	"dev09_01/internal/scrapper"
	"dev09_01/internal/sitemap"
)

func main() {
	numberOfConnections := 3

	var domain string
	var newDomain string
	flag.StringVar(&domain, "u", "", "old url")
	flag.StringVar(&newDomain, "new", "", "new url")
	flag.Parse()

	if domain == "" {
		fmt.Printf("url is empty!")
		return
	}

	s := scrapper.New()
	d := download.New(domain, newDomain)

	fmt.Println("Domain:", d.Conf.OldDomain)
	if newDomain != "" {
		fmt.Println("New Domain: ", d.Conf.NewDomain)
	}
	fmt.Println("numberOfConnections:", numberOfConnections)

	if numberOfConnections < 1 {
		fmt.Println("There can't be less than 1 numberOfConnections")
		os.Exit(1)
	}

	scanning := make(chan int, numberOfConnections)
	newLinks := make(chan []commons.Links, 100000)
	pages := make(chan commons.Page, 100000)
	attachments := make(chan []string, 100000)
	started := make(chan int, 100000)
	finished := make(chan int, 100000)

	var indexed, forSitemap, files []string

	seen := make(map[string]bool)

	start := time.Now()

	defer func() {
		close(newLinks)
		close(pages)
		close(started)
		close(finished)
		close(scanning)

		fmt.Printf("\nDuration: %s\n", time.Since(start))
		fmt.Printf("Number of pages: %6d\n", len(indexed))
	}()

	resp, err := http.Get(domain)
	if err != nil {
		fmt.Printf("Error on getting http: %v", err)
		return
	}

	defer resp.Body.Close()

	commons.Root = resp.Request.URL.String()

	s.TakeLinks(domain, started, finished, scanning, newLinks, pages, attachments)
	seen[domain] = true

	for {
		select {
		case links := <-newLinks:
			for _, link := range links {
				if !seen[link.Href] {
					seen[link.Href] = true
					go s.TakeLinks(link.Href, started, finished, scanning, newLinks, pages, attachments)
				}
			}
		case page := <-pages:
			if !s.IsURLInSlice(page.URL, indexed) {
				indexed = append(indexed, page.URL)
			}

			if !page.NoIndex {
				if !s.IsURLInSlice(page.URL, forSitemap) {
					forSitemap = append(forSitemap, page.URL)
				}
			}
		case attachment := <-attachments:
			for _, link := range attachment {
				if !s.IsURLInSlice(link, files) {
					files = append(files, link)
				}
			}
		}

		if len(started) > 1 && len(scanning) == 0 && len(started) == len(finished) {
			break
		}
	}

	for _, attachedFile := range files {
		if strings.Contains(attachedFile, ".css") {
			moreAttachments := s.GetInsideAttachments(attachedFile)
			for _, link := range moreAttachments {
				if !s.IsURLInSlice(link, files) {
					fmt.Println("Appended: ", link)
					files = append(files, link)
				}
			}
		}
	}

	os.MkdirAll("website", 0755)

	sitemap.CreateSitemap(forSitemap, "website/sitemap.xml")

	d.All(indexed)

	d.Attachments(files)
}
