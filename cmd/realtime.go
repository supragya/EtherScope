package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/instrumentation"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/node"
	"github.com/spf13/cobra"
)

var globalLogger logger.Logger

// RootCmd represents the base command when called without any subcommands
var RealtimeCmd = &cobra.Command{
	Use:   "realtime",
	Short: "run geth indexer in realtime",
	Long:  `run geth indexer in realtime`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if cfgFile == "" {
			cfgFile = util.GetUserHomedir() + "/.blockpour/bgidx/config.yaml"
		}
		util.ENOK(config.LoadViperConfig(cfgFile))

		// Setup logger
		log, err := logger.NewDefaultLogger(logLevel)
		if err != nil {
			panic(err)
		}
		globalLogger = log
	},
	Run: StartRealtimeNode,
}

func StartRealtimeNode(cmd *cobra.Command, args []string) {
	var log = globalLogger

	log.Info("setting up a new indexer node")
	_n, err := node.NewNodeWithViperFields(globalLogger)
	if err != nil {
		log.Fatal(err.Error(), nil)
	}

	if err := _n.Start(context.Background()); err != nil {
		log.Fatal("error while starting node", "error", err.Error())
	}
	instrumentation.StartPromServer(log)
	handleSig(_n, log)
}

// Keep listening for SIGTERM / SIGINT and handle graceful shutdown
func handleSig(_n service.Service, log logger.Logger) {
	c := make(chan os.Signal, 4)
	go signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	const MOSTGRACE = 3
	grace := MOSTGRACE
	lastSignal := time.Now()
	for {
		signal := <-c

		if time.Since(lastSignal) > time.Minute {
			grace = MOSTGRACE
		}
		lastSignal = time.Now()

		if grace > 2 {
			log.Info("encountered os signal. shutting down services", "signal", signal)
			grace--
			_n.Stop()
			_n.Wait()
			os.Exit(1)
		} else if grace > 0 {
			log.Warn("requesting forceful shutdown, indexer may end up in an inconsistent state",
				"signal", signal,
				"grace", grace-1)
			grace--
			_n.Stop()
			_n.Wait()
			os.Exit(1)
		} else {
			log.Fatal("too many signals, shutting down forcefully", "signal", signal)
		}
	}
}
