package search

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	elastic "github.com/olivere/elastic"
	"golang.org/x/net/context"
)

type Searcher struct {
	elasticsearchClient *elastic.Client
	elasticsearchindex  string
}

type SearchResult struct {
	NumberOfHits int64         `json:"hits"`
	Articles     []Article     `json:"articles"`
	Aggregations []Aggregation `json:"aggregations"`
}

type Aggregation struct {
	Key   string `json:"key"`
	Count int64  `json:"count"`
}

func (s *Searcher) SearchArticleGroup(articleGroup string, from int, size int) (*SearchResult, error) {
	termQuery := elastic.NewMatchQuery("ArticleGroup", articleGroup)
	searchResult, err := s.search(termQuery, from, size)

	if err != nil {
		return nil, err
	}

	return &SearchResult{
		NumberOfHits: searchResult.TotalHits(),
		Articles:     parseArticles(searchResult.Hits.Hits),
	}, nil
}

func (s *Searcher) SearchArticleGroupSalesStart(articleGroup string, startDate time.Time, from int, size int) (*SearchResult, error) {
	query := elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("ArticleGroup", articleGroup)).
		Filter(elastic.NewRangeQuery("SalesStart").
			Gte(time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 0, 0, 0, 0, time.UTC)).
			Lte(time.Date(startDate.Year(), startDate.Month(), startDate.Day(), 23, 59, 59, 0, time.UTC)))
	searchResult, err := s.search(query, from, size)

	if err != nil {
		return nil, err
	}

	return &SearchResult{
		NumberOfHits: searchResult.TotalHits(),
		Articles:     parseArticles(searchResult.Hits.Hits),
	}, nil
}

func (s *Searcher) ArticleGroupSalesStartHistogram(articleGroup string, startDate time.Time) (*SearchResult, error) {

	query := elastic.NewBoolQuery().
		Must(elastic.NewMatchQuery("ArticleGroup", articleGroup)).
		Filter(elastic.NewRangeQuery("SalesStart").Gte(startDate.Format("2006-01-02")))

	aggregation := elastic.NewDateHistogramAggregation().Field("SalesStart").Interval("day").MinDocCount(1).Format("yyyy-MM-dd")

	searchResult, err := s.elasticsearchClient.Search().
		Index(s.elasticsearchindex).
		Query(query).
		Aggregation("aggs", aggregation).
		From(0).Size(0).
		Pretty(true). // pretty print request and response JSON
		Do(context.TODO())

	if err != nil {
		return nil, err
	}

	agg, found := searchResult.Aggregations.DateHistogram("aggs")

	var aggregations []Aggregation
	if found {
		aggregations = parseAggregations(agg.Buckets)
	}

	return &SearchResult{
		Aggregations: aggregations,
	}, nil
}

func (s *Searcher) search(query elastic.Query, from int, size int) (*elastic.SearchResult, error) {

	searchResult, err := s.elasticsearchClient.Search().
		Index(s.elasticsearchindex).
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
		aggregations[i] = Aggregation{Key: *dateBucket.KeyAsString, Count: dateBucket.DocCount}
	}

	return aggregations
}

func NewSearcher(serverURL string, indexName string) *Searcher {

	client, err := retryConnect(15, 5*time.Second, func() (*elastic.Client, error) {
		return elastic.NewClient(elastic.SetURL(serverURL))
	})
	if err != nil {
		panic(err)
	}

	return &Searcher{
		elasticsearchClient: client,
		elasticsearchindex:  indexName,
	}
}

func retryConnect(attempts int, sleep time.Duration, callback func() (*elastic.Client, error)) (client *elastic.Client, err error) {
	for i := 0; ; i++ {
		client, err := callback()
		if err == nil {
			return client, nil
		}

		if i >= (attempts - 1) {
			break
		}

		time.Sleep(sleep)

		log.Println("Retry connecting to datastore after error:", err)
	}
	return nil, fmt.Errorf("after %d attempts, last error: %s", attempts, err)
}
