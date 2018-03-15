package main

import (
	"flag"
	"os"
)

type Configuration struct {
	ElasticsearchURL   string
	ElasticsearchIndex string
}

func LoadConfig() Configuration {
	elasticsearchURL := flag.String("es-url", "http://localhost:9200", "ElasticSearch URL")
	elasticsearchIndexName := flag.String("es-index-name", "articles", "ElasticSearch index name")
	flag.Parse()

	configuration := Configuration{
		ElasticsearchURL:   os.Getenv("ELASTICSEARCH_URL"),
		ElasticsearchIndex: os.Getenv("ELASTICSEARCH_INDEX"),
	}

	if configuration.ElasticsearchURL == "" {
		configuration.ElasticsearchURL = *elasticsearchURL
	}
	if configuration.ElasticsearchIndex == "" {
		configuration.ElasticsearchIndex = *elasticsearchIndexName
	}

	return configuration
}
