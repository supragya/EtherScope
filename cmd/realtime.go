package cmd

import (
	"runtime"

	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	"github.com/Blockpour/Blockpour-Geth-Indexer/indexer"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// RootCmd represents the base command when called without any subcommands
var RealtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "run geth indexer in realtime",
	Long:  `run geth indexer in realtime`,
	PreRun: func(cmd *cobra.Command, args []string) {
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
	log.Info("trying to connect to database")
	dbconn, err := db.SetupConnection()
	util.ENOK(err)
	mostRecent := dbconn.GetMostRecentPostedBlockHeight()
	log.Info("creating a new realtime indexer from ", mostRecent+1)
	var ri indexer.Indexer = indexer.NewRealtimeIndexer(mostRecent,
		viper.GetStringSlice("rpc"),
		&dbconn,
		viper.GetStringSlice("general.eventsToIndex"))
	ri.Init()
	log.Info("starting realtime indexer")
	util.ENOK(ri.Start())
}
