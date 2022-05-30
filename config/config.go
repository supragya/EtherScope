package config

import (
	"errors"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/spf13/viper"
)

var (
	mandatoryFields = [...]string{
		"general.network",
		"general.start_block",
		"general.maxBlockSpanPerCall",
		"general.chainid",

		"db.type", "db.host", "db.port", "db.user",
		"db.pass", "db.dbname", "db.sslmode",
		"db.datatable", "db.metatable",

		"rpc",
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
