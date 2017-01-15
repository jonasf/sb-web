package main

import "fmt"

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
}
