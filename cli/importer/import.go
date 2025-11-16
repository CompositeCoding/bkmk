package importer

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"bkmk/lib"
	"bkmk/store"

	"golang.org/x/net/html"
)

// Importer takes a bookmarks HTML file path and imports the links into the index.
func Importer(path string) error {
	file, err := os.Open(path)
	if err != nil {
		lib.LogError(fmt.Errorf("error: could not open file - %v", err), 0)
		return err
	}
	defer file.Close()

	doc, err := html.Parse(file)
	if err != nil {
		lib.LogError(fmt.Errorf("error: could not parse file - %v", err), 0)
		return err
	}

	// All bookmark files have this string
	if doc.FirstChild == nil || doc.FirstChild.Data != "netscape-bookmark-file-1" {
		lib.LogError(errors.New("please provide a valid bookmark file"), 0)
		return errors.New("invalid bookmark file")
	}

	// recursive function to traverse html tree
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attribute := range node.Attr {
				if attribute.Key == "href" {
					splitOnQueryParam := strings.Split(attribute.Val, `?`)
					_ = store.AddDomain(splitOnQueryParam[0], "")
				}
			}
		}
		for c := node.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return nil
}
