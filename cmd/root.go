package cmd

import (
	"fmt"

	"github.com/supragya/EtherScope/version"
	"github.com/spf13/cobra"
)

var cfgFile string
var logLevel string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:     "escope",
	Short:   "escope is a go-ethereum indexer",
	Long:    `escope is a go-ethereum indexer`,
	Version: version.RootCmdVersion,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Incorrect invocation. See escope --help for subcommands.")
	},
}

func init() {
	RootCmd.AddCommand(RealtimeCmd)
	RootCmd.AddCommand(OracleCmd)
	RootCmd.AddCommand(ConfigGen)
	RootCmd.AddCommand(AlgoIndexerCmd)

	RootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "loglevel (default is INFO)")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.supragya/escope/config.yaml)")
}
