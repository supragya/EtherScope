package config

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/util"
	"github.com/spf13/viper"
)

func LoadConfig(file string) {
	viper.SetConfigFile(file)
	err := viper.ReadInConfig()
	util.ENOK(err)

	err = validateConfig()
	util.ENOK(err)
}
