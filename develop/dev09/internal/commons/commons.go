package commons

type Page struct {
	URL       string
	Canonical string
	Links     []Links
	NoIndex   bool
}

type Links struct {
	Href string
}

var Root string

const PATH = "website/"

func IsInSlice(search string, array []string) bool {
	for _, val := range array {
		if val == search {
			return true
		}
	}

	return false
}

func IsFinal(url string) bool {
	if string(url[len(url)-1]) == "/" {
		return true
	}

	return false
}

func RemoveLastSlash(url string) string {
	if string(url[len(url)-1]) == "/" {
		return url[:len(url)-1]
	}
	return url
}
