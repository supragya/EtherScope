package node

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/ethrpc"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	localbackend "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	outputsink "github.com/Blockpour/Blockpour-Geth-Indexer/services/output_sink"
)

type Node struct {
	service.Service

	log          logger.Logger
	EthRPC       ethrpc.EthRPC             // HA upstream connection to rpc nodes, uses mspool
	LocalBackend localbackend.LocalBackend // Local database for caching / processing
	OutputSink   outputsink.OutputSink     // Consumer for offloading processed data
}

func NewNode(localBackend localbackend.LocalBackend,
	outputSink outputsink.OutputSink) (*Node, error) {
	return &Node{
		LocalBackend: localBackend,
		OutputSink:   outputSink,
	}, nil
}

func NewNodeWithViperFields(log logger.Logger) (*Node, error) {
	// Setup local backend
	localBackend, err := localbackend.NewBadgerDBWithViperFields()
	if err != nil {
		return nil, err
	}

	// Setup output link
	outputSink, err := outputsink.NewRabbitMQOutputSinkWithViperFields()
	if err != nil {
		return nil, err
	}

	// Setup indexer
	// var ri indexer.Indexer = indexer.NewRealtimeIndexer(mostRecent,
	// 	viper.GetString("rpc.master"),
	// 	viper.GetStringSlice("rpc.slaves"),
	// 	viper.GetDuration("rpc.timeout"),
	// 	&dbconn,
	// 	viper.GetStringSlice("general.eventsToIndex"))
	// ri.Init()
	return &Node{
		log: log,
		EthRPC: ,
	}
}
