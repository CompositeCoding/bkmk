package main

import (
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/google/uuid"
)

var index bleve.Index
var err error

type Domain struct {
	ID    string
	Value string
	Alias string
}

func queryDomains(query string) ([]Domain, error) {

	toQuery := bleve.NewWildcardQuery("*" + query + "*")
	searchRequest := bleve.NewSearchRequest(toQuery)
	searchRequest.Fields = []string{"Value"}

	var returnArray []Domain
	result, err := index.Search(searchRequest)

	if err != nil {
		return nil, err
	}

	for _, hit := range result.Hits {
		val := hit.Fields["Value"].(string)
		alias := hit.Fields["Alias"].(string)
		domainInstance := Domain{Value: val, ID: hit.ID, Alias: alias}
		returnArray = append(returnArray, domainInstance)

	}
	return returnArray, nil
}

func addDomain(domain string, alias string) error {

	tempDomain := Domain{ID: uuid.New().String(), Value: domain, Alias: alias}

	err := index.Index(tempDomain.ID, tempDomain)
	if err != nil {
		log.Panic(err)
	}

	return nil
}

func deleteDomain(id string) error {
	err = index.Delete(id)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	index, err = bleve.Open("index.bleve")
	if err != nil {
		index, err = bleve.New("index.bleve", bleve.NewIndexMapping())
		if err != nil {
			log.Fatalf("Fatal memory error %v", err)
		}
	}
}
