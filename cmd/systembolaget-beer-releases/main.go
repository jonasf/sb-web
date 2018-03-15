package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func main() {

	configuration := LoadConfig()

	log.Println("Using Elasticsearch Server:", configuration.ElasticsearchURL)
	log.Println("Using Elasticsearch Index Name:", configuration.ElasticsearchIndex)

	requestHandler := NewRequestHandler(configuration.ElasticsearchURL, configuration.ElasticsearchIndex)

	r := mux.NewRouter()
	r.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public/"))))
	r.HandleFunc("/", requestHandler.HomeHandler)
	r.HandleFunc("/salesstart/{date}", requestHandler.SalesStartDateHandler)

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Listening on port 8080")
	log.Fatal(srv.ListenAndServe())
}
