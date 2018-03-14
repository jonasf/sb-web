package main

import (
	"github.com/gorilla/mux"
	search "github.com/jonasf/sb-web/internal/systembolaget-beer-releases"
	"html/template"
	"log"
	"net/http"
	"time"
)

type Search interface {
	SearchArticleGroup(articleGroup string, from int, size int) (*search.SearchResult, error)
	SearchArticleGroupSalesStart(articleGroup string, startDate time.Time, from int, size int) (*search.SearchResult, error)
	ArticleGroupSalesStartHistogram(articleGroup string, startDate time.Time) (*search.SearchResult, error)
}

type RequestHandler struct {
	searcher  Search
	templates map[string]*template.Template
}

func (s *RequestHandler) HomeHandler(w http.ResponseWriter, r *http.Request) {

	fromDate := time.Now().AddDate(0, 0, -2)
	aggResult, err := s.searcher.ArticleGroupSalesStartHistogram("Öl", time.Date(fromDate.Year(), fromDate.Month(), fromDate.Day(), 0, 0, 0, 0, time.UTC))

	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if err := s.templates["index"].Execute(w, struct{ Aggregations []search.Aggregation }{aggResult.Aggregations}); err != nil {
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else if err := s.templates["salesstartdate"].Execute(w, struct{ Articles []search.Article }{result.Articles}); err != nil {
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
		searcher:  search.NewSearcher(serverURL),
		templates: setupTemplates(),
	}
}
