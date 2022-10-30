package config

import (
	"errors"

	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
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
		{"general.networkName", "string", "Unique chain identifier"},
		{"general.isErigon", "bool", "Are the upstreams erigon?"},
		{"general.failOnNonEthError", "bool", "Fail on catastrophic error on a log event"},
		{"general.maxCPUParallels", "uint64", "Largest number of logical CPUs the process can be running on at a given time"},
		{"general.persistence", "string", "Persistence object: one of (\"postgres\", \"mq\"). MQ doesn't suppport atomic transactions"},
		{"general.eventsToIndex", "[]string", "Events to index"},
		{"general.oracleMapsRootDir", "string", "Directory storing chainlink oracle maps"},
		{"general.diskCacheRootDir", "string", "Directory for on-disk caches"},
		{"general.prometheusEndpoint", "string", "address on which to expose prometheus metrics"},

		{"rpc.master", "string", "Remote master upstream for RPC access"},
		{"rpc.slaves", "[]string", "Remote slave upstreams for RPC access"},
		{"rpc.timeout", "time.Duration", "Timeout per RPC call (in milliseconds)"},
	}

	mandatoryFieldsERC20Transfer = [...]Field{
		{"erc20transfer.restrictionType", "string", "Restriction type: none, from, to, both, either"},
		{"erc20transfer.whitelistFile", "string", "path to file listing whitelisted addresses"},
	}

	mandatoryFieldsMessageQueue = [...]Field{
		{"mq.secureConnection", "bool", "MQ connection over secure channel"},
		{"mq.host", "string", "MQ host"},
		{"mq.port", "uint64", "MQ port"},
		{"mq.user", "string", "MQ user"},
		{"mq.pass", "string", "MQ password"},
		{"mq.queue", "string", "MQ queue name for channel"},
		{"mq.queueIsDurable", "bool", "MQ queue durability"},
		{"mq.queueAutoDelete", "bool", "MQ queue auto delete"},
		{"mq.queueExclusive", "bool", "MQ queue exclusive"},
		{"mq.queueNoWait", "bool", "MQ queue no-wait"},
		{"mq.skipResume", "bool", "false if startBlock for indexing or hit url to get resumeBlock"},
		{"mq.resumeURL", "string", "url to get resumeBlock from"},
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

	for _, event := range viper.GetStringSlice("general.eventsToIndex") {
		if event == "ERC20Transfer" {
			for _, mf := range mandatoryFieldsERC20Transfer {
				err := ensureFieldIntegrity(mf)
				if err != nil {
					return err
				}
			}
		}
	}

	dbType := viper.GetString("general.persistence")
	if dbType == "mq" {
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
		return errors.New("config error: unset mandatory field: " + f.name + " (" + f._type + "); description: " + f.err)
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
