package cmd

import (
	"time"

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
	Run:   StartRealtimeNode,
}

func StartRealtimeNode(cmd *cobra.Command, args []string) {
	log.Info("creating a new realtime indexer")
	var ri indexer.Indexer = indexer.NewRealtimeIndexer(28705483, viper.GetStringSlice("rpc"))
	ri.Init()
	log.Info("starting realtime indexer")
	util.ENOK(ri.Start())
	log.Info("starting to sleep")
	time.Sleep(time.Second * 10)
	return
}
