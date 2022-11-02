package ethrpc

import (
	"context"
	"math/big"
	"time"

	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
)

var (
	EthRPCMSPoolCFGSection   = "ethRPCMSPool"
	EthRPCMSPoolCFGNecessity = "needed if `node.ethrpc` == mspool"
	EthRPCMSPoolCFGHeader    = cfg.SArr("mspool is master slave arch based",
		"ethrpc provider for indexer which switches between ",
		"master node and multiple slave nodes for high",
		"availability. expects nodes to be in sync")
	EthRPCMSPoolCFGFields = [...]cfg.Field{
		{
			Name:      "master",
			Type:      "string",
			Necessity: "always needed",
			Info:      cfg.SArr("master node for mspool. Recommended websockets"),
			Default:   "https://rpc.ankr.com/eth",
		},
		{
			Name:      "slaves",
			Type:      "[]string",
			Necessity: "always needed",
			Info:      cfg.SArr("slave nodes for mspool. Recommended websockets"),
			Default:   "\n    - https://rpc.ankr.com/eth\n    - https://rpc.ankr.com/eth",
		},
		{
			Name:      "timeout",
			Type:      "time.Duration",
			Necessity: "always needed",
			Info: cfg.SArr("timeout per rpc call. Should never exceed",
				"ceil(windowSize x timeStep / tolerance), else a switchover",
				"may never happen"),
			Default: "800ms",
		},
		{
			Name:      "windowSize",
			Type:      "uint",
			Necessity: "always needed",
			Info: cfg.SArr("number of timeSteps to consider while rating",
				"health of a upstream"),
			Default: 80,
		},
		{
			Name:      "timeStep",
			Type:      "time.Duration",
			Necessity: "always needed",
			Info: cfg.SArr("failure reporting timestep. Max of 1 failure is",
				"reported internally to mspool no matter # of threads."),
			Default: "100ms",
		},
		{
			Name:      "tolerance",
			Type:      "uint",
			Necessity: "always needed",
			Info: cfg.SArr("number of timeSteps to fail before we switch",
				"over"),
			Default: 5,
		},
		{
			Name:      "retryTimesteps",
			Type:      "uint",
			Necessity: "always needed",
			Info: cfg.SArr("number of timeSteps to wait before giving an",
				"rpc a second chance"),
			Default: 300,
		},
		{
			Name:      "periodicRecording",
			Type:      "time.Duration",
			Necessity: "always needed",
			Info: cfg.SArr("periodically the ethrpc service displays status",
				"of the pool. This is the only reason why ethrpc is a",
				"service. Setting this to 0ms will turn periodic display off"),
			Default: "10s",
		},
	}
)

type MSPoolEthRPCImpl struct {
	service.BaseService

	// Config
	log               logger.Logger
	master            string
	slaves            []string
	timeout           time.Duration
	mspoolcfg         MSPoolConfig
	periodicRecording time.Duration

	// Internal Data Structures
	pool         *MasterSlavePool[ethclient.Client]
	localBackend lb.LocalBackend
}

// OnStart starts the badgerdb LocalBackend. It implements service.Service.
func (n *MSPoolEthRPCImpl) OnStart(ctx context.Context) error {
	pool, err := NewEthClientMasterSlavePool(n.master, n.slaves, n.mspoolcfg, n.timeout)
	if err != nil {
		return err
	}
	n.pool = pool

	if n.periodicRecording.Nanoseconds() == 0 {
		n.log.Info("mspool ethrpc reporting turned off since periodicRecording is zero")
		return nil
	}

	go n.pool.PeriodicRecording(n.periodicRecording, n.log)

	return nil
}

// OnStop stops the badgerdb LocalBackend. It implements service.Service
func (n *MSPoolEthRPCImpl) OnStop() {
}

func NewMSPoolEthRPCWithViperFields(log logger.Logger, localBackend lb.LocalBackend) (EthRPC, error) {
	lb := &MSPoolEthRPCImpl{
		log:     log,
		master:  viper.GetString(EthRPCMSPoolCFGSection + ".master"),
		slaves:  viper.GetStringSlice(EthRPCMSPoolCFGSection + ".slaves"),
		timeout: viper.GetDuration(EthRPCMSPoolCFGSection + ".timeout"),
		mspoolcfg: MSPoolConfig{
			WindowSize:     viper.GetUint32(EthRPCMSPoolCFGSection + ".windowSize"),
			ToleranceCount: viper.GetUint32(EthRPCMSPoolCFGSection + ".tolerance"),
			TimeStep:       viper.GetDuration(EthRPCMSPoolCFGSection + ".timeStep"),
			RetryTimesteps: viper.GetUint32(EthRPCMSPoolCFGSection + ".retryTimesteps"),
		},
		periodicRecording: viper.GetDuration(EthRPCMSPoolCFGSection + ".periodicRecording"),
		localBackend:      localBackend,
	}
	lb.BaseService = *service.NewBaseService(log, "ethrpc", lb)
	return lb, nil
}

// Non-cached RPC access to get sender address for any eth transaction
func (n *MSPoolEthRPCImpl) GetTxSender(txHash common.Hash,
	blockHash common.Hash,
	txIdx uint) (common.Address, error) {
	tx, err := Do(n.pool,
		func(ctx context.Context, c *ethclient.Client) (*types.Transaction, error) {
			tx, _, err := c.TransactionByHash(ctx, txHash)
			return tx, err
		}, nil)
	if err != nil {
		return common.Address{}, err
	}

	sender, err := Do(n.pool,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			return c.TransactionSender(ctx, tx, blockHash, txIdx)
		}, common.Address{})
	return sender, err
}

// Non-cached RPC access to get current block height
func (n *MSPoolEthRPCImpl) GetCurrentBlockHeight() (uint64, error) {
	return Do(n.pool,
		func(ctx context.Context, c *ethclient.Client) (uint64, error) {
			return c.BlockNumber(ctx)
		}, 0)
}

// Non-cached RPC access to get block timestamp
func (n *MSPoolEthRPCImpl) GetBlockTimestamp(height uint64) (uint64, error) {
	header, err := Do(n.pool,
		func(ctx context.Context, c *ethclient.Client) (*types.Header, error) {
			return c.HeaderByNumber(ctx, big.NewInt(int64(height)))
		}, nil)
	if err != nil {
		n.log.Error("block timestamp error ", "error", err)
		return 0, err
	}
	return header.Time, err
}

// Non-cached RPC access to get filtered logs
func (n *MSPoolEthRPCImpl) GetFilteredLogs(fq ethereum.FilterQuery) ([]types.Log, error) {
	return Do(n.pool,
		func(ctx context.Context, c *ethclient.Client) ([]types.Log, error) {
			return c.FilterLogs(ctx, fq)
		}, []types.Log{})
}
