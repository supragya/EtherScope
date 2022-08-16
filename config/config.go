package config

import (
	"errors"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/spf13/viper"
)

type Field struct {
	name  string
	_type string
	err   string
}

var (
	mandatoryFields = [...]Field{
		{"general.network", "string", "Type of network being supported. Only \"evm\" supported as of now"},
		{"general.startBlock", "uint64", "Block to start indexing from. Overridden if postgres tables have higher heights already synced"},
		{"general.maxBlockSpanPerCall", "uint64", "Max number of block to sync once in a batch"},
		{"general.chainID", "uint64", "ChainID of the synced chain"},
		{"general.failOnNonEthError", "bool", "Fail on catastrophic error on a log event"},
		{"general.persistence", "string", "Persistence object: one of (\"postgres\", \"mq\"). MQ doesn't suppport atomic transactions"},
		{"general.eventsToIndex", "[]string", "Events to index"},

		{"rpc", "[]string", "Remote upstreams for RPC access"},
	}

	mandatoryFieldsPostgres = [...]Field{
		{"postgres.host", "string", "Postgresql DB host"},
		{"postgres.port", "uint64", "Postgresql DB port"},
		{"postgres.user", "string", "Postgresql DB user"},
		{"postgres.pass", "string", "Postgresql DB password"},
		{"postgres.dbname", "string", "Postgresql DB dbname"},
		{"postgres.sslmode", "string", "Postgresql DB sslmode. One of (\"enable\", \"disable\")"},
		{"postgres.datatable", "string", "Postgresql DB to store mint, burn, swap in"},
		{"postgres.metatable", "string", "Postgresql DB to store meta information about indexing"},
	}
	mandatoryFieldsMessageQueue = [...]Field{
		{"mq.host", "string", "MQ host"},
		{"mq.port", "uint64", "MQ port"},
		{"mq.user", "string", "MQ user"},
		{"mq.pass", "string", "MQ password"},
		{"mq.queue", "string", "MQ queue name for channel"},
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
		err := ensureFieldIntegrity(mf)
		if err != nil {
			return err
		}
	}

	dbType := viper.GetString("general.persistence")
	if dbType == "postgres" {
		for _, mf := range mandatoryFieldsPostgres {
			err := ensureFieldIntegrity(mf)
			if err != nil {
				return err
			}
		}
	} else if dbType == "mq" {
		for _, mf := range mandatoryFieldsMessageQueue {
			err := ensureFieldIntegrity(mf)
			if err != nil {
				return err
			}
		}
	} else {
		return errors.New("unknown persistence type: " + dbType)
	}
	return nil
}

func ensureFieldIntegrity(f Field) error {
	if !viper.IsSet(f.name) {
		return errors.New("unset mandatory field: " + f.name + " (" + f._type + "); description: " + f.err)
	}
	var castOK bool = true
	var item = viper.Get(f.name)

	switch f._type {
	case "string":
		_, castOK = item.(string)
	case "uint64":
		_, castOK = item.(int)
	case "bool":
		_, castOK = item.(bool)
	case "[]string":
		var subItems []interface{}
		subItems, castOK = item.([]interface{})
		if !castOK {
			break
		}
		for _, subItem := range subItems {
			_, cOK := subItem.(string)
			castOK = castOK && cOK
		}
	}

	if !castOK {
		return errors.New("mandatory field type invalid: " + f.name + " (" + f._type + "); description: " + f.err)
	}
	return nil
}
