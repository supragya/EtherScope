package node

import (
	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
)

var (
	NodeCFGSection   = "node"
	NodeCFGNecessity = "general evm events indexing"
	NodeCFGHeader    = cfg.SArr("node is core indexing service for bgidx",
		"node is tasked with initiating other services such as",
		"localbackend (badger-db) and outputsink (rabbitmq)")
	NodeCFGFields = [...]cfg.Field{
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
			Default: 15865859,
		},
		{
			Name:      "skipResumeRemote",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("disables fetch for blockheight from remoteHTTP"),
			Default:   false,
		},
		{
			Name:      "skipResumeLocal",
			Type:      "bool",
			Necessity: "always needed",
			Info:      cfg.SArr("disables fetch for blockheight from localbackend"),
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
			Name:      "localBackendType",
			Type:      "string",
			Necessity: "always needed",
			Info: cfg.SArr("type of local backend indexer should use.",
				"only possible type right now is `badgerdb`"),
			Default: "badgerdb",
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
			Name:      "eventsToIndex",
			Type:      "[]string",
			Necessity: "always needed",
			Info: cfg.SArr("ethereum events to index. events listed here",
				"are not guaranteed to be the only calls made",
				"to underlying rpc for processing, but are guaranteed",
				"to be the only events presented to the output sink",
				"could be one or many of the following:",
				"- UniswapV2Swap",
				"- UniswapV2Mint",
				"- UniswapV2Burn",
				"- UniswapV3Swap",
				"- UniswapV3Mint",
				"- UniswapV3Burn",
				"- ERC20Transfer"),
			Default: "\n    - ERC20Transfer\n    - UniswapV2Swap",
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
			Name:      "pricingChainlinkOraclesDumpFile",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("dump file containing list of chainlink oracles to trust"),
			Default:   "chainlink_oracle_dumpfile.csv",
		},
		{
			Name:      "pricingDexDumpFile",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("dump file containing list of dexes to take into account historically"),
			Default:   "dex_dumpfile.csv",
		},
	}
)
