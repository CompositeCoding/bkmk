package main

import (
	"fmt"

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

		var domain Domain

		val := hit.Fields["Value"].(string)
		if hit.Fields["Alias"] != nil {
			domain = Domain{Value: val, ID: hit.ID, Alias: hit.Fields["Alias"].(string)}
		} else {
			domain = Domain{Value: val, ID: hit.ID}
		}

		returnArray = append(returnArray, domain)

	}
	return returnArray, nil
}

func addDomain(domain string, alias string) error {

	var tempDomain Domain

	if alias != "" {
		tempDomain = Domain{ID: uuid.New().String(), Value: domain, Alias: alias}
	} else {
		tempDomain = Domain{ID: uuid.New().String(), Value: domain}
	}

	err := index.Index(tempDomain.ID, tempDomain)
	if err != nil {
		return err
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

	config := ReadOrCreateConfig()
	var profile string = fmt.Sprintf("%v.bleve", config.Profile)

	index, err = bleve.Open(profile)
	if err != nil {
		index, err = bleve.New(profile, bleve.NewIndexMapping())
		if err != nil {
			log_error(err, 2)
		}
	}
}
