package main

import (
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/net/context"
	elastic "gopkg.in/olivere/elastic.v5"
)

type Searcher struct {
	elasticsearchClient *elastic.Client
}

type SearchResult struct {
	numberOfHits int64
	articles     []Article
	aggregations []Aggregation
}

type Aggregation struct {
	key   string
	count int64
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

func (s *Searcher) ArticleGroupSalesStartHistogram(articleGroup string, startDate time.Time) (*SearchResult, error) {

	query := elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("ArticleGroup", articleGroup)).
		Filter(elastic.NewRangeQuery("SalesStart").Gte(startDate.Format("2006-01-02")))

	aggregation := elastic.NewDateHistogramAggregation().Field("SalesStart").Interval("day").MinDocCount(1)

	searchResult, err := s.elasticsearchClient.Search().
		Index("articles").
		Query(query).
		Aggregation("aggs", aggregation).
		From(0).Size(0).
		Pretty(true). // pretty print request and response JSON
		Do(context.TODO())

	if err != nil {
		return nil, err
	}

	agg, found := searchResult.Aggregations.DateHistogram("aggs")
	parseAggregations(agg.Buckets)

	var aggregations []Aggregation
	if found {
		aggregations = parseAggregations(agg.Buckets)
	}

	return &SearchResult{
		aggregations: aggregations,
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

func parseAggregations(items []*elastic.AggregationBucketHistogramItem) []Aggregation {
	aggregations := make([]Aggregation, len(items))

	for i, dateBucket := range items {
		aggregations[i] = Aggregation{key: *dateBucket.KeyAsString, count: dateBucket.DocCount}
	}

	return aggregations
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
