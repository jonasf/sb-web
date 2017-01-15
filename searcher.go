package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Searcher struct {
	elasticsearchClient *elastic.Client
}

type SearchResult struct {
	numberOfHits int64
	articles     []Article
}

func (s *Searcher) SearchArticleGroup(articleGroup string, from int, size int) (*SearchResult, error) {
	termQuery := elastic.NewMatchQuery("ArticleGroup", articleGroup)
	searchResult, err := s.search(termQuery, from, size)

	if err != nil {
		return nil, err
	}

	return &SearchResult{
		numberOfHits: searchResult.TotalHits(),
		articles:     parseArticles(searchResult.Hits.Hits),
	}, nil
}

func (s *Searcher) search(query elastic.Query, from int, size int) (*elastic.SearchResult, error) {

	searchResult, err := s.elasticsearchClient.Search().
		Index("articles").
		Query(query).
		//Sort("Name", true). // sort by "Name" field, ascending
		From(from).Size(size).
		Pretty(true). // pretty print request and response JSON
		Do(context.TODO())
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return searchResult, nil
}

func parseArticles(hits []*elastic.SearchHit) []Article {

	articles := make([]Article, len(hits))

	for i, hit := range hits {
		var article Article
		err := json.Unmarshal(*hit.Source, &article)

		if err != nil {
			fmt.Println(err)
		}

		articles[i] = article
	}

	return articles
}

func NewSearcher() Searcher {

	client, err := elastic.NewClient()
	if err != nil {
		panic(err)
	}

	return Searcher{
		elasticsearchClient: client,
	}
}
