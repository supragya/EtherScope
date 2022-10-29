package cmd

import (
	"runtime"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	"github.com/Blockpour/Blockpour-Geth-Indexer/indexer"
	"github.com/Blockpour/Blockpour-Geth-Indexer/instrumentation"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/node"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var globalLogger logger.Logger

// RootCmd represents the base command when called without any subcommands
var RealtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "run geth indexer in realtime",
	Long:  `run geth indexer in realtime`,
	PreRun: func(cmd *cobra.Command, args []string) {
		// Setup logger
		log, err := logger.NewDefaultLogger(logLevel)
		if err != nil {
			panic(err)
		}
		globalLogger = log

		maxParallelsRequested := viper.GetInt("general.maxCPUParallels")
		if maxParallelsRequested > runtime.NumCPU() {
			log.Warn("running on fewer threads than requested parallels: ", runtime.NumCPU(), " vs requested ", maxParallelsRequested)
			maxParallelsRequested = runtime.NumCPU()
		}

		runtime.GOMAXPROCS(maxParallelsRequested)
		log.Info("set runtime max parallelism: ", maxParallelsRequested)
	},
	Run: StartRealtimeNode,
}

func StartRealtimeNode(cmd *cobra.Command, args []string) {
	node.SetupNodeWithViperFields(globalLogger)
	var log = globalLogger

	// Setup local backend

	// Setup output link
	log.Info("trying to connect to database")
	dbconn, err := db.SetupConnection()
	util.ENOK(err)
	mostRecent := dbconn.GetMostRecentPostedBlockHeight()

	// Setup indexer
	var ri indexer.Indexer = indexer.NewRealtimeIndexer(mostRecent,
		viper.GetString("rpc.master"),
		viper.GetStringSlice("rpc.slaves"),
		viper.GetDuration("rpc.timeout"),
		&dbconn,
		viper.GetStringSlice("general.eventsToIndex"))
	ri.Init()

	// Start services
	go instrumentation.StartPromServer(log.With("module", "promserver"))
	util.ENOK(ri.Start())

	// Keep listening for SIGTERM / SIGINT and handle graceful shutdown
}
