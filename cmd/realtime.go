package cmd

import (
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/indexer"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var RealtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "run geth indexer in realtime",
	Long:  `run geth indexer in realtime`,
	Run:   StartRealtimeNode,
}

func StartRealtimeNode(cmd *cobra.Command, args []string) {
	// _, err := db.SetupConnection()
	// util.ENOK(err)
	ri := indexer.NewRealtimeIndexer(2, []string{"hello"})
	log.Info("Starting realtime indexer")
	util.ENOK(ri.Start())
	log.Info("Starting to sleep")
	time.Sleep(time.Second * 10)
	return
}
