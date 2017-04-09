package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

type SearchWrapper struct {
	searcher Searcher
}

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

	searchWrapper := &SearchWrapper{searcher: NewSearcher(esURL)}

	/*result, err := searcher.SearchArticleGroupSalesStart("Öl", time.Date(2017, 4, 7, 0, 0, 0, 0, time.UTC), 0, 50)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Found a total of %d articles\n", result.NumberOfHits)

	for _, article := range result.Articles {
		fmt.Printf("Article %s: %s\n", article.Name, article.SecondaryName)
	}*/

	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/salesstarts", searchWrapper.SalesStartsHandler)
	r.HandleFunc("/salesstart/{date}", searchWrapper.SalesStartDateHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on port 8080")
	log.Fatal(srv.ListenAndServe())
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Under contruction")
}

func (s *SearchWrapper) SalesStartsHandler(w http.ResponseWriter, r *http.Request) {
	fromDate := time.Now().AddDate(0, 0, -2)
	aggResult, err := s.searcher.ArticleGroupSalesStartHistogram("Öl", time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC))

	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(aggResult.Aggregations)
}

func (s *SearchWrapper) SalesStartDateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	releaseDate, err := time.Parse("2006-01-02", vars["date"])
	if err != nil {
		log.Println(err)
	}

	result, err := s.searcher.SearchArticleGroupSalesStart("Öl", time.Date(releaseDate.Year(), releaseDate.Month(), releaseDate.Day(), 0, 0, 0, 0, time.UTC), 0, 50)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result.Articles)
}
