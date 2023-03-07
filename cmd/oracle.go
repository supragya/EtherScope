package cmd

import (
	"context"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/services/oraclenode"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var OracleCmd = &cobra.Command{
	Use:   "oracle",
	Short: "run oracle indexer",
	Long:  `run oracle indexer`,
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
	Run: StartOracleNode,
}

func StartOracleNode(cmd *cobra.Command, args []string) {
	var log = globalLogger

	log.Info("setting up a new oracle geth indexer node")
	_n, err := oraclenode.NewOracleNodeWithViperFields(globalLogger)
	if err != nil {
		log.Fatal(err.Error(), nil)
	}

	if err := _n.Start(context.Background()); err != nil {
		log.Fatal("error while starting oracle node", "error", err.Error())
	}

	handleSig(_n, log)
}
