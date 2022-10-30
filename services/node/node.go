package node

import (
	"context"
	"sync"

	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	outs "github.com/Blockpour/Blockpour-Geth-Indexer/services/output_sink"
)

type NodeImpl struct {
	service.BaseService

	log logger.Logger
	// EthRPC       ethrpc.EthRPC             // HA upstream connection to rpc nodes, uses mspool
	LocalBackend lb.LocalBackend // Local database for caching / processing
	OutputSink   outs.OutputSink // Consumer for offloading processed data
}

// OnStart starts the Node. It implements service.Service.
func (n *NodeImpl) OnStart(ctx context.Context) error {
	if err := n.LocalBackend.Start(ctx); err != nil {
		return err
	}

	if err := n.OutputSink.Start(ctx); err != nil {
		return err
	}

	return nil
}

// OnStop stops the Node. It implements service.Service
func (n *NodeImpl) OnStop() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		n.LocalBackend.Stop()
		n.LocalBackend.Stop()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		n.OutputSink.Stop()
		n.OutputSink.Stop()
	}()

	wg.Wait()
}

// func NewNode(localBackend localbackend.LocalBackend,
// 	outputSink outputsink.OutputSink) (*Node, error) {
// 	return &Node{
// 		LocalBackend: localBackend,
// 		OutputSink:   outputSink,
// 	}, nil
// }

func NewNodeWithViperFields(log logger.Logger) (service.Service, error) {
	// Setup local backend
	localBackend, err := lb.NewBadgerDBWithViperFields(log.With("service", "localbackend"))
	if err != nil {
		return nil, err
	}

	// Setup output link
	outputSink, err := outs.NewRabbitMQOutputSinkWithViperFields(log.With("service", "outputsink"))
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
	node := &NodeImpl{
		log:          log.With("service", "node"),
		LocalBackend: localBackend,
		OutputSink:   outputSink,
	}
	node.BaseService = *service.NewBaseService(log, "node", node)
	return node, nil
}
