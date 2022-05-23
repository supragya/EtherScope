package cmd

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/version"
	"github.com/spf13/cobra"
)

var cfgFile string
var logLevel string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "bgidx",
	Short:   "bgidx is a go-ethereum indexer",
	Long:    `bgidx is a go-ethereum indexer`,
	Version: version.RootCmdVersion,
	Run:     func(cmd *cobra.Command, args []string) {},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.ENOK(logger.SetLogLevel(logLevel))

		if cfgFile == "" {
			cfgFile = util.GetUserHomedir() + "/.blockpour/bgidx/config.yaml"
		}
		util.ENOK(config.LoadViperConfig(cfgFile))
		util.ENOK(config.CheckViperMandatoryFields())
	},
}

func init() {
	RootCmd.AddCommand(RealtimeCmd)
	// RootCmd.AddCommand(backfill.BackfillCmd)

	RootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "loglevel (default is INFO)")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.blockpour/bgidx/config.yaml)")
}
