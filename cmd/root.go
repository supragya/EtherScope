package cmd

import (
	"fmt"
	"runtime"

	"github.com/Blockpour/Blockpour-Geth-Indexer/config"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/Blockpour/Blockpour-Geth-Indexer/version"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		util.ENOK(logger.SetLogLevel(logLevel))

		if cfgFile == "" {
			cfgFile = util.GetUserHomedir() + "/.blockpour/bgidx/config.yaml"
		}
		util.ENOK(config.LoadViperConfig(cfgFile))
		util.ENOK(config.CheckViperMandatoryFields())

		maxParallelsRequested := viper.GetInt("general.maxCPUParallels")
		if maxParallelsRequested > runtime.NumCPU() {
			log.Warn("running on fewer threads than requrested parallels: %v vs requested %v", runtime.NumCPU(), maxParallelsRequested)
			maxParallelsRequested = runtime.NumCPU()
		}

		runtime.GOMAXPROCS(maxParallelsRequested)
		log.Info("set runtime max parallelism: %v", maxParallelsRequested)
	},
}

func init() {
	RootCmd.AddCommand(RealtimeCmd)

	RootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "info", "loglevel (default is INFO)")
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.blockpour/bgidx/config.yaml)")
}
