package cmd

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/db"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
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
	_, err := db.SetupConnection()
	util.ENOK(err)
}
