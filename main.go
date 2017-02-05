package main

import (
	"fmt"
	"time"
)

//import "fmt"

func main() {
	searcher := NewSearcher()

	result, err := searcher.SearchArticleGroup("Öl", 0, 10)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Found a total of %d articles\n", result.numberOfHits)

	for _, article := range result.articles {
		fmt.Printf("Article %s: %s\n", article.Name, article.SecondaryName)
	}

	aggResult, err := searcher.ArticleGroupSalesStartHistogram("Öl", time.Date(2017, 1, 1, 12, 0, 0, 0, time.UTC))

	if err != nil {
		panic(err)
	}

	for _, agg := range aggResult.aggregations {
		fmt.Printf("Aggregation %s: %d\n", agg.key, agg.count)
	}
}
