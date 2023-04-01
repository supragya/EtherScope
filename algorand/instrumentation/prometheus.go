package instrumentation

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func StartPromServer() {
	address := viper.GetString("general.prometheusEndpoint")
	log.Info("starting metrics server at ", address)
	http.Handle("/", promhttp.Handler())
	http.ListenAndServe(address, nil)
}

var (
	StartingBlock  = pg("starting_block", "starting block")
	ProcessedBlock = pg("processed_block", "processed block")
	CurrentBlock   = pg("current_block", "current blockchain height")
	SwapFound      = pc("swap_found", "swap found")
	SwapProcessed  = pc("swap_processed", "swap processed")
)

func pc(name string, help string) prometheus.Counter {
	return promauto.NewCounter(prometheus.CounterOpts{
		Name: "indexer_" + name,
		Help: help,
	})
}

func pg(name string, help string) prometheus.Gauge {
	return promauto.NewGauge(prometheus.GaugeOpts{
		Name: "indexer_" + name,
		Help: help,
	})
}
