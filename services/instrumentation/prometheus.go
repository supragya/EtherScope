package instrumentation

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"

	logger "github.com/supragya/EtherScope/libs/log"
)

func StartPromServer(log logger.Logger) {
	address := viper.GetString("general.prometheusEndpoint")
	log.Info("starting metrics server", "address", address)
	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(address, nil)
}

var (
	StartingBlock  = pg("starting_block", "starting block")
	ProcessedBlock = pg("processed_block", "processed block")
	CurrentBlock   = pg("current_block", "current blockchain height")

	TfrFound    = pc("tfr_found", "transfers found")
	MintV2Found = pc("mintv2_found", "mint v2 found")
	MintV3Found = pc("mintv3_found", "mint v3 found")
	BurnV2Found = pc("burnv2_found", "burn v2 found")
	BurnV3Found = pc("burnv3_found", "burn v3 found")
	SwapV2Found = pc("swapv2_found", "swapv2 found")
	SwapV3Found = pc("swapv3_found", "swapv3 found")

	TfrProcessed    = pc("tfr_processed", "transfers processed")
	MintV2Processed = pc("mintv2_processed", "mint v2 processed")
	MintV3Processed = pc("mintv3_processed", "mint v3 processed")
	BurnV2Processed = pc("burnv2_processed", "burn v2 processed")
	BurnV3Processed = pc("burnv3_processed", "burn v3 processed")
	SwapV2Processed = pc("swapv2_processed", "swapv2 processed")
	SwapV3Processed = pc("swapv3_processed", "swapv3 processed")

	TraderJoeV2SwapFound     = pc("trader_joe_v2_swap_found", "trader joe v2 swap found")
	TraderJoeV2SwapProcessed = pc("trader_joe_v2_swap_processed", "trader joe v2 swap processed")
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
