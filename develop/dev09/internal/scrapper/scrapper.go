package scrapper

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"dev09_01/internal/commons"

	"golang.org/x/net/html"
)

var (
	extensions = []string{".png", ".jpg", ".jpeg", ".json", ".js", ".tiff", ".pdf", ".txt", ".gif", ".psd", ".ai", "dwg", ".bmp", ".zip", ".tar", ".gzip", ".svg", ".avi", ".mov", ".json", ".xml", ".mp3", ".wav", ".mid", ".ogg", ".acc", ".ac3", "mp4", ".ogm", ".cda", ".mpeg", ".avi", ".swf", ".acg", ".bat", ".ttf", ".msi", ".lnk", ".dll", ".db", ".css"}
	falseURLs  = []string{"mailto:", "javascript:", "tel:", "whatsapp:", "callto:", "wtai:", "sms:", "market:", "geopoint:", "ymsgr:", "msnim:", "gtalk:", "skype:"}
	validURL   = regexp.MustCompile(`\(([^()]*)\)`)
	validCSS   = regexp.MustCompile(`\{(\s*?.*?)*?\}`)
)

type Scrapper struct{}

func New() Scrapper {
	return Scrapper{}
}

func (s *Scrapper) getLinks(domain string) (page commons.Page, attachments []string, err error) {
	resp, err := http.Get(domain)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	page.URL = domain
	foundMeta := false

	var f func(*html.Node)
	f = func(n *html.Node) {
		for _, a := range n.Attr {
			if a.Key == "style" {
				if strings.Contains(a.Val, "url(") {
					found := s.getURLEmbeeded(a.Val)
					if found != "" {
						link, err := resp.Request.URL.Parse(found)
						if err == nil {
							foundLink := s.sanitizeURL(link.String())
							if s.isValidAttachment(foundLink) {
								attachments = append(attachments, foundLink)
							}
						}
					}

				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "meta" {
			for _, a := range n.Attr {
				if a.Key == "name" && a.Val == "robots" {
					foundMeta = true
				}
				if foundMeta {
					if a.Key == "content" && strings.Contains(a.Val, "noindex") {
						page.NoIndex = true
					}
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "link" {
			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						} else if s.isValidLink(foundLink) {
							page.Links = append(page.Links, commons.Links{Href: foundLink})
						}
					}
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "script" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "img" {
			for _, a := range n.Attr {
				if a.Key == "src" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
				if a.Key == "srcset" {
					links := strings.Split(a.Val, " ")
					for _, val := range links {
						link, err := resp.Request.URL.Parse(val)
						if err == nil {
							foundLink := s.sanitizeURL(link.String())
							if s.isValidAttachment(foundLink) {
								attachments = append(attachments, foundLink)
							}
						}
					}
				}
			}
		}

		if n.Type == html.ElementNode && n.Data == "a" {
			ok := false
			newLink := commons.Links{}

			for _, a := range n.Attr {
				if a.Key == "href" {
					link, err := resp.Request.URL.Parse(a.Val)
					if err == nil {
						foundLink := s.sanitizeURL(link.String())
						if s.isValidLink(foundLink) {
							ok = true
							newLink.Href = foundLink
						} else if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}

			}

			if ok && !s.doesLinkExist(newLink, page.Links) {
				page.Links = append(page.Links, newLink)
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return
}

func (s *Scrapper) TakeLinks(toScan string, started chan int, finished chan int, scanning chan int, newLinks chan []commons.Links, pages chan commons.Page, attachments chan []string) {
	started <- 1
	scanning <- 1
	defer func() {
		<-scanning
		finished <- 1
		fmt.Printf("\rStarted: %6d - Finished %6d", len(started), len(finished))
	}()

	page, attached, err := s.getLinks(toScan)
	if err != nil {
		return
	}

	pages <- page

	attachments <- attached

	newLinks <- page.Links
}

func (s *Scrapper) isInternLink(link string) bool {
	if strings.Index(link, commons.Root) == 0 {
		return true
	}
	return false
}

func (s *Scrapper) removeQuery(link string) string {
	return strings.Split(link, "?")[0]
}

func (s *Scrapper) isStart(link string) bool {
	if strings.Compare(link, commons.Root) == 0 {
		return true
	}
	return false
}

func (s *Scrapper) sanitizeURL(link string) string {
	for _, fal := range falseURLs {
		if strings.Contains(link, fal) {
			return ""
		}
	}

	link = strings.TrimSpace(link)

	if string(link[len(link)-1]) != "/" {
		link = link + "/"
	}

	tram := strings.Split(link, "#")[0]

	return tram
}

func (s *Scrapper) IsValidExtension(link string) bool {
	for _, extension := range extensions {
		if strings.Contains(strings.ToLower(link), extension) {
			return false
		}
	}
	return true
}

func (s *Scrapper) isValidLink(link string) bool {
	if s.isInternLink(link) && !s.isStart(link) && s.IsValidExtension(link) {
		return true
	}

	return false
}

func (s *Scrapper) isValidAttachment(link string) bool {
	if s.isInternLink(link) && !s.isStart(link) && !s.IsValidExtension(link) {
		return true
	}

	return false
}

func (s *Scrapper) doesLinkExist(newLink commons.Links, existingLinks []commons.Links) (exists bool) {
	for _, val := range existingLinks {
		if strings.Compare(newLink.Href, val.Href) == 0 {
			exists = true
		}
	}

	return
}

func (s *Scrapper) IsURLInSlice(search string, array []string) bool {
	withSlash := search[:len(search)-1]
	withoutSlash := search

	if string(search[len(search)-1]) == "/" {
		withSlash = search
		withoutSlash = search[:len(search)-1]
	}

	for _, val := range array {
		if val == withSlash || val == withoutSlash {
			return true
		}
	}

	return false
}

func (s *Scrapper) IsLinkScanned(link string, scanned []string) (exists bool) {
	for _, val := range scanned {
		if strings.Compare(link, val) == 0 {
			exists = true
		}
	}

	return
}

func (s *Scrapper) getURLEmbeeded(body string) (url string) {
	valid := validURL.Find([]byte(body))
	if valid == nil {
		return
	}

	url = string(valid)

	if string(url[0]) == `(` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `)` {
		url = url[:len(url)-1]
	}

	if string(url[0]) == `"` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `"` {
		url = url[:len(url)-1]
	}

	if string(url[0]) == `'` {
		url = url[1:]
	}
	if string(url[len(url)-1]) == `'` {
		url = url[:len(url)-1]
	}

	return url
}

func (s *Scrapper) GetInsideAttachments(url string) (attachments []string) {
	if commons.IsFinal(url) {
		url = commons.RemoveLastSlash(url)
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	body := buf.String()

	if strings.Contains(body, "url(") {
		blocks := validCSS.FindAll([]byte(body), -1)
		for _, b := range blocks {
			rules := strings.Split(string(b), ";")
			for _, r := range rules {
				found := s.getURLEmbeeded(r)
				if found != "" {
					link, err := resp.Request.URL.Parse(found)
					if err == nil {
						foundLink := s.sanitizeURL(link.String())
						if s.isValidAttachment(foundLink) {
							attachments = append(attachments, foundLink)
						}
					}
				}
			}
		}
	}

	return
}
