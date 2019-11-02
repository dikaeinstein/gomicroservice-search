package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/dikaeinstein/gomicroservice-search/data"
	log "github.com/sirupsen/logrus"
)

func setupTest(r *searchRequest) (*http.Request, *httptest.ResponseRecorder, *Search, *data.MockStore) {
	mockStore := &data.MockStore{}

	s, _ := statsd.New("127.0.0.1:8125")
	h := NewSearch(mockStore, s, log.New())

	rw := httptest.NewRecorder()

	if r == nil {
		return httptest.NewRequest(
			http.MethodGet,
			"/search",
			nil,
		), rw, h, mockStore
	}

	return httptest.NewRequest(
		http.MethodGet,
		"/search?criteria="+url.QueryEscape(r.Criteria),
		nil,
	), rw, h, mockStore
}

func TestSearchHandlerReturnsBadRequestWhenNoSearchCriteriaIsSent(t *testing.T) {
	r, rw, handler, _ := setupTest(nil)

	handler.Handle(rw, r)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v", rw.Code)
	}
}

func TestSearchHandlerReturnsBadRequestWhenBlankSearchCriteriaIsSent(t *testing.T) {
	r, rw, handler, _ := setupTest(&searchRequest{""})

	handler.Handle(rw, r)

	if rw.Code != http.StatusBadRequest {
		t.Errorf("Expected BadRequest got %v", rw.Code)
	}
}

func TestSearchHandlerCallsDataStoreWithValidQuery(t *testing.T) {
	r, rw, handler, mockStore := setupTest(&searchRequest{"Fat Freddy's Cat"})
	mockStore.On("Search", "Fat Freddy's Cat").Return(make([]data.Kitten, 0))

	handler.Handle(rw, r)

	mockStore.AssertExpectations(t)
}

func TestSearchHandlerReturnsKittensWithValidQuery(t *testing.T) {
	r, rw, handler, mockStore := setupTest(&searchRequest{"Fat Freddy's Cat"})
	mockStore.On("Search", "Fat Freddy's Cat").Return(make([]data.Kitten, 1))

	handler.Handle(rw, r)

	response := searchResponse{}
	json.Unmarshal(rw.Body.Bytes(), &response)

	if http.StatusOK != rw.Code {
		t.Errorf("request failed with code: %d", rw.Code)
	}
	if len(response.Kittens) != 1 {
		t.Errorf("expected %d, got %d", 1, len(response.Kittens))
	}
}
