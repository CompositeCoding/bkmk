package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/blevesearch/bleve/v2"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var index, _ = bleve.Open("index.bleve")

type Domain struct {
	ID    string
	Value string
}

func queryDomains(query string) ([]string, error) {

	toQuery := bleve.NewWildcardQuery("*" + query + "*")
	searchRequest := bleve.NewSearchRequest(toQuery)
	searchRequest.Fields = []string{"Value"}

	var returnArray []string
	result, err := index.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	for _, hit := range result.Hits {
		if val, ok := hit.Fields["Value"].(string); ok {
			returnArray = append(returnArray, val)
		}
	}
	return returnArray, nil
}

func addDomain(domain string) error {

	tempDomain := Domain{ID: uuid.New().String(), Value: domain}

	err := index.Index(tempDomain.ID, tempDomain)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func init() {

	index, err := bleve.New("index.bleve", bleve.NewIndexMapping())

	if err == nil {
		log.Print("Creating index")

		index.Close()
	} else if strings.Contains(err.Error(), "cannot create new index, path already exists") {
		log.Print("index exists")
		return
	} else {
		log.Print(err)
	}
}

func main() {

	var rootCmd = &cobra.Command{Use: "bkmk"}

	var cmdAdd = &cobra.Command{
		Use:   "add [string to add]",
		Short: "Add a new url to your bookmarks",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("Added item:", args[0])
		},
	}

	var cmdOpen = &cobra.Command{
		Use:   "open [path]",
		Short: "Open a bookmark",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			log.Print("Opening item at path:", args[0])
		},
	}

	rootCmd.AddCommand(cmdAdd, cmdOpen)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
