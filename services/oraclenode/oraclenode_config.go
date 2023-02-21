package oraclenode

import (
	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
)

var (
	OracleNodeCFGSection   = "oraclenode"
	OracleNodeCFGNecessity = "oracle indexing"
	OracleNodeCFGHeader    = cfg.SArr("oraclenode is chainlink indexing service for bgidx",
		"oraclenode is tasked with initiating other services such as",
		"localbackend (badger-db) and outputsink (rabbitmq)")
	OracleNodeCFGFields = [...]cfg.Field{
		{
			Name:      "moniker",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("custom name given to node to differentiate messages",
				"at the outputSink"),
			Default: "blockpour-geth-node",
		},
		{
			Name:      "prodcheck",
			Type:      "bool",
			Necessity: "always needed",
			Info: cfg.SArr("ensures only tagged and released version of bgidx is",
				"allowed to run."),
			Default: true,
		},
		{
			Name:      "maxCPUParallels",
			Type:      "uint",
			Necessity: "always needed",
			Info:      cfg.SArr("maximum number of CPU threads to give to bgidx"),
			Default:   4,
		},
		{
			Name:      "network",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("evm compatible network name"),
			Default:   "ethereum-mainnet",
		},
		{
			Name:      "startBlock",
			Type:      "uint64",
			Necessity: "always needed",
			Info: cfg.SArr("user defined blockheight to start sync from.",
				"this may be overriden at runtime using resume from localbackend",
				"and remoteHTTP endpoint"),
			Default: 12864088,
		},
		{
			Name:      "skipResumeRemote",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("disables fetch for blockheight from remoteHTTP"),
			Default:   false,
		},
		{
			Name:      "remoteResumeURL",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("remoteHTTP URL for fetching blockheight to resume from"),
			Default:   "https://myremote.blockpour.com",
		},
		{
			Name:      "outputSinkType",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("type of output sink backend indexer should",
				"offload indexed information to. only possible type",
				"right now is `rabbitmq`"),
			Default: "rabbitmq",
		},
		{
			Name:      "ethRPCType",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("type of ethrpc handler to route requests",
				"through. only possible type right now is `mspool`"),
			Default: "mspool",
		},
		{
			Name:      "maxBlockSpanPerCall",
			Type:      "uint64",
			Necessity: "always needed",
			Info: cfg.SArr("number of blocks to fetch logs for at the",
				"beginning of processing loop"),
			Default: 5,
		},
		{
			Name:      "chainlinkFeedRegistry",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("contract address where feeds are registered and confirmed"),
			Default:   "0x47Fb2585D2C56Fe188D0E6ec628a38b74fCeeeDf",
		},
		{
			Name:      "feedFile",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("file saving feeds for chainlink"),
			Default:   "feeds.csv",
		},
	}
)
