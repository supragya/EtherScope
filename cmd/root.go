package cmd

import (
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Incorrect invocation. See bgidx --help for subcommands.")
	},
}

func init() {
	RootCmd.AddCommand(RealtimeCmd)
	RootCmd.AddCommand(ConfigGen)

	RootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "loglevel (default is INFO)")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.blockpour/bgidx/config.yaml)")
}
