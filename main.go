package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port string = ":8080"

func main() {
	searcher := NewSearcher()

	/*result, err := searcher.SearchArticleGroup("Öl", 0, 10)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Found a total of %d articles\n", result.numberOfHits)

	for _, article := range result.articles {
		fmt.Printf("Article %s: %s\n", article.Name, article.SecondaryName)
	}
	*/

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Under contruction")
	})

	http.HandleFunc("/salesstarts", func(w http.ResponseWriter, r *http.Request) {
		fromDate := time.Now().AddDate(0, 0, -2)
		aggResult, err := searcher.ArticleGroupSalesStartHistogram("Öl", time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC))

		if err != nil {
			log.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(aggResult.Aggregations)
	})

	log.Println("Starting server on port: ", port)
	log.Println(http.ListenAndServe(port, nil))
}
