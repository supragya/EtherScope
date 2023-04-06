package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/node"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var AlgoIndexerCmd = &cobra.Command{
	Use:   "algoindexer",
	Short: "run geth indexer in realtime",
	Long:  `run geth indexer in realtime`,
	Run:   StartIndexer,
}

func StartIndexer(cmd *cobra.Command, args []string) {
	log, err := logger.NewDefaultLogger("info")
	if err != nil {
		log.Fatal(err.Error(), nil)
	}

	n, err := node.NewNode(
		viper.GetUint64("general.startBlock"),
		viper.GetString("rpc.algodUrl"),
		viper.GetString("rpc.indexerUrl"),
		viper.GetString("rpc.token"),
		viper.GetUint64("general.maxBlockSpanPerCall"),
		viper.GetBool("mq.skipResume"),
		viper.GetString("mq.resumeURL"),
		log,
	)

	if err != nil {
		log.Fatal(err.Error(), nil)
	}

	err = n.Start(context.Background())
	if err != nil {
		log.Fatal("error while starting indexer", "error", err.Error())
	}

	handleSigAlgorand(n, log)
}

// Keep listening for SIGTERM / SIGINT and handle graceful shutdown
// TODO: replace with the handleSig function in realtime.go
func handleSigAlgorand(n *node.Node, log logger.Logger) {
	c := make(chan os.Signal, 4)
	go signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for {
		signal := <-c
		log.Info("encountered os signal. shutting down services", "signal", signal)
		n.Stop()
		n.Wait()
		os.Exit(1)
	}
}
