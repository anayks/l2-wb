package download

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"

	"dev09_01/internal/commons"
)

type Settings struct {
	OldDomain string
	NewDomain string
}

type Downloader struct {
	Conf Settings
}

func New(oldDomain, newDomain string) Downloader {
	return Downloader{
		Conf: Settings{
			OldDomain: oldDomain,
			NewDomain: newDomain,
		},
	}
}

func (d *Downloader) download(url string, filename string, changeDomain bool) (err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	defer f.Close()

	if changeDomain && d.Conf.NewDomain != "" {
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		newStr := buf.String()

		newStr = strings.ReplaceAll(newStr, d.Conf.OldDomain, d.Conf.NewDomain)

		newContent := bytes.NewBufferString(newStr)
		_, err = io.Copy(f, newContent)

	} else {
		_, err = io.Copy(f, resp.Body)
	}

	return
}

func (d *Downloader) All(indexed []string) {
	for _, url := range indexed {
		filepath := d.GetPath(url)
		if filepath == "" {
			filepath = "/index.html"
		}

		if d.hasPaths(filepath) {
			if commons.IsFinal(filepath) {
				if !d.exists(commons.PATH + filepath) {
					os.MkdirAll(commons.PATH+filepath, 0755)
					filepath = filepath + "index.html"
				}
			} else {
				path := d.getOnlyPath(filepath)
				if !d.exists(commons.PATH + path) {
					os.MkdirAll(commons.PATH+path, 0755)
				}
			}

		}

		d.download(url, commons.PATH+filepath, true)
	}
}

func (d *Downloader) Attachments(attachments []string) {
	for _, url := range attachments {
		filepath := d.GetPath(url)
		if filepath == "" {
			continue
		}

		if d.hasPaths(filepath) {
			if commons.IsFinal(filepath) {
				filepath = commons.RemoveLastSlash(filepath)
				url = commons.RemoveLastSlash(url)
			}

			path := d.getOnlyPath(filepath)
			if !d.exists(commons.PATH + path) {
				os.MkdirAll(commons.PATH+path, 0755)
			}
		}

		d.download(url, commons.PATH+filepath, false)
	}
}

func (d *Downloader) hasPaths(url string) bool {
	paths := strings.Split(url, "/")
	if len(paths) > 1 {
		return true
	}
	return false
}

func (d *Downloader) getOnlyPath(url string) (path string) {

	paths := strings.Split(url, "/")
	if len(paths) <= 1 {
		return url
	}

	total := paths[:len(paths)-1]

	return strings.Join(total[:], "/")
}

func (d *Downloader) GetPath(link string) string {
	return strings.Replace(link, commons.Root, "", 1)
}

func (d *Downloader) exists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
