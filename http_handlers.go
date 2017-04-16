package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type RequestHandler struct {
	searcher  Searcher
	templates map[string]*template.Template
}

func (s *RequestHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {

	fromDate := time.Now().AddDate(0, 0, -2)
	aggResult, err := s.searcher.ArticleGroupSalesStartHistogram("Öl", time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC))

	if err != nil {
		log.Println(err)
	}

	if err := s.templates["index"].Execute(w, struct{ Aggregations []Aggregation }{aggResult.Aggregations}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *RequestHandler) SalesStartDateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	releaseDate, err := time.Parse("2006-01-02", vars["date"])
	if err != nil {
		log.Println(err)
	}

	result, err := s.searcher.SearchArticleGroupSalesStart("Öl", time.Date(releaseDate.Year(), releaseDate.Month(), releaseDate.Day(), 0, 0, 0, 0, time.UTC), 0, 50)
	if err != nil {
		log.Println(err)
	}

	if err := s.templates["salesstartdate"].Execute(w, struct{ Articles []Article }{result.Articles}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func setupTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)

	var baseTemplate = "templates/_base.html"
	templates["index"] = template.Must(template.ParseFiles(baseTemplate, "templates/index.html"))
	templates["salesstartdate"] = template.Must(template.ParseFiles(baseTemplate, "templates/salesstartdate.html"))

	return templates
}

func NewRequestHandler(serverURL string) RequestHandler {
	return RequestHandler{
		searcher:  NewSearcher(serverURL),
		templates: setupTemplates(),
	}
}
