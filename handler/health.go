package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// Health represents the health handler
type Health struct {
	statsd *statsd.Client
}

// Get returns the health status of this service
func (h *Health) Get(rw http.ResponseWriter, r *http.Request) {
	defer func(startTime time.Time) {
		h.statsd.Timing("health.timing", time.Since(startTime), nil, 1)
	}(time.Now())

	h.statsd.Incr("health.success", nil, 1)
	fmt.Fprintln(rw, "OK")
}

// NewHealth creates a new Health handler
func NewHealth(statsd *statsd.Client) *Health {
	return &Health{statsd: statsd}
}
