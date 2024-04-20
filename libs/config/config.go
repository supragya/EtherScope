package config

import (
	"errors"

	"github.com/supragya/EtherScope/libs/util"
	"github.com/spf13/viper"
)

type Field struct {
	Name      string
	Type      string
	Necessity string
	Info      []string
	Default   interface{}
}

var (
// nodeFields = [...]Field{
// 	{"node.startBlock", "uint64",
// 		"user defined blockheight to start sync from. this may be overriden at runtime using resume from localbackend and remoteHTTP endpoint"},
// 	{"node.skipResumeRemote", "bool",
// 		"disables fetch for blockheight from remoteHTTP"},
// 	{"node.skipResumeLocal", "bool",
// 		"disables fetch for blockheight from localbackend"},
// 	{"node.remoteResumeURL", "bool",
// 		"remoteHTTP URL for fetching blockheight to resume from"},
// 	{"general.chainID", "uint64", "ChainID of the synced chain"},
// 	{"general.networkName", "string", "Unique chain identifier"},
// 	{"general.isErigon", "bool", "Are the upstreams erigon?"},
// 	{"general.failOnNonEthError", "bool", "Fail on catastrophic error on a log event"},
// 	{"general.maxCPUParallels", "uint64", "Largest number of logical CPUs the process can be running on at a given time"},
// 	{"general.persistence", "string", "Persistence object: one of (\"postgres\", \"mq\"). MQ doesn't suppport atomic transactions"},
// 	{"general.eventsToIndex", "[]string", "Events to index"},
// 	{"general.oracleMapsRootDir", "string", "Directory storing chainlink oracle maps"},
// 	{"general.diskCacheRootDir", "string", "Directory for on-disk caches"},
// 	{"general.prometheusEndpoint", "string", "address on which to expose prometheus metrics"},

// 	{"rpc.master", "string", "Remote master upstream for RPC access"},
// 	{"rpc.slaves", "[]string", "Remote slave upstreams for RPC access"},
// 	{"rpc.timeout", "time.Duration", "Timeout per RPC call (in milliseconds)"},

// 	{"localbackend.db", "string", "path to db"},
// }

// mandatoryFieldsERC20Transfer = [...]Field{
// 	{"erc20transfer.restrictionType", "string", "Restriction type: none, from, to, both, either"},
// 	{"erc20transfer.whitelistFile", "string", "path to file listing whitelisted addresses"},
// }

//	mandatoryFieldsMessageQueue = [...]Field{
//		{"mq.secureConnection", "bool", "MQ connection over secure channel"},
//		{"mq.host", "string", "MQ host"},
//		{"mq.port", "uint64", "MQ port"},
//		{"mq.user", "string", "MQ user"},
//		{"mq.pass", "string", "MQ password"},
//		{"mq.queue", "string", "MQ queue name for channel"},
//		{"mq.queueIsDurable", "bool", "MQ queue durability"},
//		{"mq.queueAutoDelete", "bool", "MQ queue auto delete"},
//		{"mq.queueExclusive", "bool", "MQ queue exclusive"},
//		{"mq.queueNoWait", "bool", "MQ queue no-wait"},
//		{"mq.skipResume", "bool", "false if startBlock for indexing or hit url to get resumeBlock"},
//		{"mq.resumeURL", "string", "url to get resumeBlock from"},
//	}
)

func SArr(strs ...string) []string {
	return strs
}

func SFmt(strs []string) string {
	content := ""
	for _, line := range strs {
		content += " " + line
	}
	return content
}

func LoadViperConfig(file string) error {
	err := util.VerifyFileExistence(file)
	if err != nil {
		return err
	}
	viper.SetConfigFile(file)
	return viper.ReadInConfig()
}

func EnsureFieldIntegrity(section string, f Field) error {
	fieldName := section + "." + f.Name
	if !viper.IsSet(fieldName) {
		return errors.New("config error: unset mandatory field: " + fieldName + " (" + f.Type + "); description:" + SFmt(f.Info))
	}
	var castOK bool = true
	var item = viper.Get(fieldName)

	switch f.Type {
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
		return errors.New("mandatory field type invalid: " + fieldName + " (" + f.Type + "); description:" + SFmt(f.Info))
	}
	return nil
}
