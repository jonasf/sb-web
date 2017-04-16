package main

import (
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

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

	requestHandler := NewRequestHandler(esURL)

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
