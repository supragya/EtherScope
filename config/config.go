package config

import (
	"errors"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/spf13/viper"
)

var (
	mandatoryFields = [...]string{
		"network",
		"start_block",
	}
)

func LoadViperConfig(file string) error {
	err := util.VerifyFileExistence(file)
	if err != nil {
		return err
	}
	viper.SetConfigFile(file)
	return viper.ReadInConfig()
}

func CheckViperMandatoryFields() error {
	for _, mf := range mandatoryFields {
		if !viper.IsSet(mf) {
			return errors.New("unset mandatory field: " + mf)
		}
	}
	return nil
}
