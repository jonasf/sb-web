package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const port string = ":8080"

func main() {

	//Delayed start to wait for dependencies to start
	delayedStartDuration := os.Getenv("DELAYED-START")
	if delayedStartDuration != "" {
		duration, _ := time.ParseDuration(delayedStartDuration)
		log.Println("Delaying start for ", duration)
		time.Sleep(duration)
	}

	esURL := os.Getenv("ES-URL")

	if esURL == "" {
		esURLPtr := flag.String("es-url", "http://localhost:9200", "ElasticSearch URL")
		flag.Parse()
		esURL = *esURLPtr
	}
	log.Println("Elasticsearch server address: ", esURL)

	searcher := NewSearcher(esURL)

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
