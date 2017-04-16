package main

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
)

type SearchStub struct{}

func (histogramStub *SearchStub) ArticleGroupSalesStartHistogram(articleGroup string, startDate time.Time) (*SearchResult, error) {

	agg := make([]Aggregation, 1)
	agg[0] = Aggregation{Key: "2017-05-05", Count: 2}
	return &SearchResult{
		Aggregations: agg,
	}, nil
}

func (histogramStub *SearchStub) SearchArticleGroup(articleGroup string, from int, size int) (*SearchResult, error) {
	return nil, nil
}

func (histogramStub *SearchStub) SearchArticleGroupSalesStart(articleGroup string, startDate time.Time, from int, size int) (*SearchResult, error) {
	articles := make([]Article, 1)
	articles[0] = Article{Name: "Chimay"}

	return &SearchResult{
		NumberOfHits: 1,
		Articles:     articles,
	}, nil
}

func TestHomeHandlerReturnAggregation(t *testing.T) {

	requestHandler := &RequestHandler{searcher: &SearchStub{}, templates: setupTemplates()}

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(requestHandler.HomeHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

func TestSalesStartDateHandlerReturnData(t *testing.T) {

	requestHandler := &RequestHandler{searcher: &SearchStub{}, templates: setupTemplates()}

	req, err := http.NewRequest("GET", "/salesstart/2017-05-05", nil)

	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/salesstart/{date}", requestHandler.SalesStartDateHandler)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
}

type SearchFailStub struct{}

func (histogramStub *SearchFailStub) ArticleGroupSalesStartHistogram(articleGroup string, startDate time.Time) (*SearchResult, error) {
	return nil, errors.New("Stuff went terribly wrong")
}

func (histogramStub *SearchFailStub) SearchArticleGroup(articleGroup string, from int, size int) (*SearchResult, error) {
	return nil, nil
}

func (histogramStub *SearchFailStub) SearchArticleGroupSalesStart(articleGroup string, startDate time.Time, from int, size int) (*SearchResult, error) {
	return nil, errors.New("Stuff went terribly wrong")
}

func TestHomeHandlerSearchFail(t *testing.T) {

	requestHandler := &RequestHandler{searcher: &SearchFailStub{}, templates: setupTemplates()}

	req, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	http.HandlerFunc(requestHandler.HomeHandler).ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusInternalServerError, status)
	}
}

func TestSalesStartDateHandlerSearchFail(t *testing.T) {

	requestHandler := &RequestHandler{searcher: &SearchFailStub{}, templates: setupTemplates()}

	req, err := http.NewRequest("GET", "/salesstart/2017-05-05", nil)

	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()

	r := mux.NewRouter()
	r.HandleFunc("/salesstart/{date}", requestHandler.SalesStartDateHandler)
	r.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusInternalServerError, status)
	}
}
