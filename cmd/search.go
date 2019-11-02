package main

import (
	"net/http"
	"os"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/dikaeinstein/gomicroservice-search/config"
	"github.com/dikaeinstein/gomicroservice-search/data"
	"github.com/dikaeinstein/gomicroservice-search/handler"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const address = ":8082"

func main() {
	var logger = &log.Logger{
		Out:       os.Stdout,
		Formatter: new(log.TextFormatter),
		Level:     log.DebugLevel,
	}

	cfg := config.New()
	store, err := data.NewMySQLStore(cfg.MysqlConnection)
	if err != nil {
		log.Fatal(err)
	}

	c, err := statsd.New(cfg.DogStatsD)
	if err != nil {
		log.Fatal(err)
	}
	c.Namespace = "gomicroservice.search."

	search := handler.NewSearch(store, c, logger)
	health := handler.NewHealth(c)

	r := mux.NewRouter()
	r.HandleFunc("/search", search.Handle).Methods(http.MethodGet)
	r.HandleFunc("/healthz", health.Get).Methods(http.MethodGet)

	logger.WithField("service", "search").
		Infof("Starting server, listening on %s", address)
	log.WithField("service", "search").
		Fatal(http.ListenAndServe(address, r))
}
