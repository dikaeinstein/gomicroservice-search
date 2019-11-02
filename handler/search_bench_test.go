package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/dikaeinstein/gomicroservice-search/data"
	log "github.com/sirupsen/logrus"
)

func BenchmarkSearchHandler(b *testing.B) {
	mockStore := &data.MockStore{}
	mockStore.On("Search", "Fat Freddy's Cat").Return([]data.Kitten{
		data.Kitten{
			ID:     "2",
			Name:   "Fat Freddy's Cat",
			Weight: 20.0,
		},
	})

	s, _ := statsd.New("127.0.0.1:8125")
	search := NewSearch(mockStore, s, log.New())

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		r := httptest.NewRequest(
			http.MethodGet,
			"/search?criteria="+url.QueryEscape("Fat Freddy's Cat"),
			nil,
		)
		rr := httptest.NewRecorder()
		search.Handle(rr, r)
	}
}
