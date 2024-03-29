package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Function that takes a bookmark file path and imports into index
func importer(path string) error {

	file, err := os.Open(path)

	if err != nil {
		log_error(fmt.Errorf("error: could not open file - %d", err), 0)
		return err
	}

	doc, err := html.Parse(file)

	if err != nil {
		log_error(fmt.Errorf("error: could not parse file - %d", err), 0)
		return err
	}

	// All bookmark files have this string
	if doc.FirstChild.Data != "netscape-bookmark-file-1" {
		log_error(errors.New("please provide a valid bookmark file"), 0)
		return err
	}

	// recursive function to traverse html tree
	var f func(*html.Node)
	f = func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "a" {
			for _, attribute := range node.Attr {
				if attribute.Key == "href" {
					split_on_query_param := strings.Split(attribute.Val, `?`)
					addDomain(split_on_query_param[0], "")
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
