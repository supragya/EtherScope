package node

import (
	"context"
	"sync"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	outs "github.com/Blockpour/Blockpour-Geth-Indexer/services/output_sink"
	"github.com/spf13/viper"
)

var (
	NodeSection   = "node"
	NodeCFGHeader = cfg.SArr("node is core indexing service for bgidx",
		"node is tasked with initiating other services such as",
		"localbackend (badger-db) and outputsink (rabbitmq)")
	NodeCFGFields = [...]cfg.Field{
		{
			Name: "startBlock",
			Type: "uint64",
			Err: cfg.SArr("user defined blockheight to start sync from.",
				"this may be overriden at runtime using resume from localbackend",
				"and remoteHTTP endpoint"),
			Default: 15865859,
		},
		{
			Name:    "skipResumeRemote",
			Type:    "bool",
			Err:     cfg.SArr("disables fetch for blockheight from remoteHTTP"),
			Default: false,
		},
		{
			Name:    "skipResumeLocal",
			Type:    "bool",
			Err:     cfg.SArr("disables fetch for blockheight from localbackend"),
			Default: false,
		},
		{
			Name:    "remoteResumeURL",
			Type:    "bool",
			Err:     cfg.SArr("remoteHTTP URL for fetching blockheight to resume from"),
			Default: "https://myremote.blockpour.com",
		},
	}
)

type NodeImpl struct {
	service.BaseService

	log logger.Logger
	// EthRPC       ethrpc.EthRPC             // HA upstream connection to rpc nodes, uses mspool
	LocalBackend lb.LocalBackend // Local database for caching / processing
	OutputSink   outs.OutputSink // Consumer for offloading processed data

	// Configs
	startBlock       uint64 // User defined startBlock
	skipResumeRemote bool   // skip checking remote for resume height
	skipResumeLocal  bool   // skip checking localbackend for resume height
	remoteResumeURL  string // URL to use for resume height GET request
}

// OnStart starts the Node. It implements service.Service.
func (n *NodeImpl) OnStart(ctx context.Context) error {
	if err := n.LocalBackend.Start(ctx); err != nil {
		return err
	}

	if err := n.OutputSink.Start(ctx); err != nil {
		return err
	}

	// Do height syncup using both LocalBackend and remote http
	// startHeight, err := n.getResumeHeight()

	// Loop for impl
	go n.Loop()

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

// Loop implements core indexing logic
func (n *NodeImpl) Loop() {

}

// Creates a new node service with spf13/viper fields (yaml)
// CONTRACT: NodeCFGFields enlists all the fields to be accessed in this function
func NewNodeWithViperFields(log logger.Logger) (service.Service, error) {
	// ensure field integrity for viper
	for _, mf := range NodeCFGFields {
		err := cfg.EnsureFieldIntegrity(NodeSection, mf)
		if err != nil {
			return nil, err
		}
	}

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

	node := &NodeImpl{
		log:              log.With("service", "node"),
		LocalBackend:     localBackend,
		OutputSink:       outputSink,
		startBlock:       viper.GetUint64(NodeSection + ".startBlock"),
		skipResumeRemote: viper.GetBool(NodeSection + ".skipResumeRemote"),
		skipResumeLocal:  viper.GetBool(NodeSection + ".skipResumeLocal"),
		remoteResumeURL:  viper.GetString(NodeSection + ".remoteResumeURL"),
	}
	node.BaseService = *service.NewBaseService(log, "node", node)
	return node, nil
}
