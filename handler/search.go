package handler

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/dikaeinstein/gomicroservice-search/data"
	log "github.com/sirupsen/logrus"
)

type searchRequest struct {
	Criteria string `json:"criteria"`
}

type searchResponse struct {
	Kittens []data.Kitten `json:"kittens"`
}

// Search is an http handler for the search service
type Search struct {
	dataStore data.Store
	statsd    *statsd.Client
	logger    *log.Logger
}

// NewSearch creates a new Search handler
func NewSearch(d data.Store, s *statsd.Client, l *log.Logger) *Search {
	return &Search{d, s, l}
}

var defaultFields = log.Fields{
	"service": "search",
	"handler": "search",
}

// Handle uses the search query to retrieve list of kittens
func (s *Search) Handle(w http.ResponseWriter, r *http.Request) {
	defer func(startTime time.Time) {
		s.statsd.Timing("search.timing.total", time.Since(startTime), nil, 1)
	}(time.Now())

	request := searchRequest{}
	v, err := url.ParseQuery(r.URL.RawQuery)
	request.Criteria = v.Get("criteria")
	if err != nil || len(request.Criteria) < 1 {
		s.statsd.Incr("search.badrequest", nil, 1)

		s.logger.WithFields(defaultFields).Error(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	startTime := time.Now()
	kittens, err := s.dataStore.Search(request.Criteria)
	s.statsd.Timing("search.timing.data", time.Since(startTime), nil, 1)
	if err != nil {
		s.logger.WithFields(defaultFields).Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if kittens == nil {
		kittens = []data.Kitten{}
	}

	w.Header().Add("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(&searchResponse{Kittens: kittens})

	s.statsd.Incr("search.success", nil, 1)
}
