package node

import (
	"context"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
)

type NodeImpl struct {
	service.BaseService

	log logger.Logger
	// EthRPC       ethrpc.EthRPC             // HA upstream connection to rpc nodes, uses mspool
	// LocalBackend localbackend.LocalBackend // Local database for caching / processing
	// OutputSink   outputsink.OutputSink     // Consumer for offloading processed data
}

// OnStart starts the Node. It implements service.Service.
func (n *NodeImpl) OnStart(ctx context.Context) error {
	n.log.Info("i have started")
	return nil
}

// OnStop stops the Node. It implements service.Service
func (n *NodeImpl) OnStop() {
	n.log.Error("i have stopped")
}

// func NewNode(localBackend localbackend.LocalBackend,
// 	outputSink outputsink.OutputSink) (*Node, error) {
// 	return &Node{
// 		LocalBackend: localBackend,
// 		OutputSink:   outputSink,
// 	}, nil
// }

func NewNodeWithViperFields(log logger.Logger) (*NodeImpl, error) {
	// Setup local backend
	// localBackend, err := localbackend.NewBadgerDBWithViperFields()
	// if err != nil {
	// 	return nil, err
	// }

	// Setup output link
	// outputSink, err := outputsink.NewRabbitMQOutputSinkWithViperFields()
	// if err != nil {
	// 	return nil, err
	// }

	// Setup indexer
	// var ri indexer.Indexer = indexer.NewRealtimeIndexer(mostRecent,
	// 	viper.GetString("rpc.master"),
	// 	viper.GetStringSlice("rpc.slaves"),
	// 	viper.GetDuration("rpc.timeout"),
	// 	&dbconn,
	// 	viper.GetStringSlice("general.eventsToIndex"))
	// ri.Init()
	node := &NodeImpl{log: log}
	node.BaseService = *service.NewBaseService(log, "rinode", node)
	return node, nil
}
