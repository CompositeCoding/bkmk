package main

import (
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// Function that takes a bookmark file path and imports into index
func importer(path string) {

	file, err := os.Open(path)

	if err != nil {
		log.Print(err)
	}

	doc, err := html.Parse(file)

	if err != nil {
		log.Print(err)
	}

	// All bookmark files have this string
	if doc.FirstChild.Data != "netscape-bookmark-file-1" {
		log.Println("Please provide a valid bookmark file")
		return
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
}
