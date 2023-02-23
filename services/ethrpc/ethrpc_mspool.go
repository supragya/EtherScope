package ethrpc

import (
	"context"
	"math/big"
	"time"

	"github.com/Blockpour/Blockpour-Geth-Indexer/assets/abi/ERC20"
	"github.com/Blockpour/Blockpour-Geth-Indexer/assets/abi/chainlink"
	"github.com/Blockpour/Blockpour-Geth-Indexer/assets/abi/univ2pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/assets/abi/univ3pair"
	"github.com/Blockpour/Blockpour-Geth-Indexer/assets/abi/univ3positionsnft"
	cfg "github.com/Blockpour/Blockpour-Geth-Indexer/libs/config"
	logger "github.com/Blockpour/Blockpour-Geth-Indexer/libs/log"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/service"
	lb "github.com/Blockpour/Blockpour-Geth-Indexer/services/local_backend"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	lru "github.com/hashicorp/golang-lru"
	"github.com/spf13/viper"
	"golang.org/x/sync/semaphore"
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
			Name:      "maxParallels",
			Type:      "uint",
			Necessity: "always needed",
			Info:      cfg.SArr("maximum parallel rpc calls to upstream at a time"),
			Default:   20,
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
		{
			Name:      "cacheSizeContractTokens",
			Type:      "uint32",
			Necessity: "always needed",
			Info: cfg.SArr("number of slots for LRU cache in-memory for storing",
				"contract token sides for queried contracts"),
			Default: 20000,
		},
		{
			Name:      "cacheSizeERC20",
			Type:      "uint32",
			Necessity: "always needed",
			Info: cfg.SArr("number of slots for LRU cache in-memory for storing",
				"erc20 decimal places for queried contracts"),
			Default: 20000,
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
	maxParallels      uint

	// Internal Data Structures
	pool         *MasterSlavePool[ethclient.Client]
	localBackend lb.LocalBackend
	sem          *semaphore.Weighted

	// In-memory caches
	cacheContractTokens *lru.ARCCache
	cacheERC20          *lru.ARCCache
	cacheERC20Name      *lru.ARCCache
}

// OnStart starts the badgerdb LocalBackend. It implements service.Service.
func (n *MSPoolEthRPCImpl) OnStart(ctx context.Context) error {
	pool, err := NewEthClientMasterSlavePool(n.master, n.slaves, n.mspoolcfg, n.timeout, n.log)
	if err != nil {
		return err
	}
	n.pool = pool
	n.sem = semaphore.NewWeighted(int64(n.maxParallels))

	if n.periodicRecording.Nanoseconds() == 0 {
		n.log.Info("mspool ethrpc reporting turned off since periodicRecording is zero")
		return nil
	}

	go n.pool.PeriodicRecording(n.periodicRecording)

	return nil
}

// OnStop stops the badgerdb LocalBackend. It implements service.Service
func (n *MSPoolEthRPCImpl) OnStop() {
}

func NewMSPoolEthRPCWithViperFields(log logger.Logger, localBackend lb.LocalBackend) (EthRPC, error) {
	// ensure field integrity for viper
	for _, mf := range EthRPCMSPoolCFGFields {
		err := cfg.EnsureFieldIntegrity(EthRPCMSPoolCFGSection, mf)
		if err != nil {
			return nil, err
		}
	}

	cacheContractTokens, err := lru.NewARC(viper.GetInt(EthRPCMSPoolCFGSection + ".cacheSizeContractTokens"))
	if err != nil {
		return nil, err
	}
	cacheERC20, err := lru.NewARC(viper.GetInt(EthRPCMSPoolCFGSection + ".cacheSizeERC20"))
	if err != nil {
		return nil, err
	}
	cacheERC20Name, err := lru.NewARC(viper.GetInt(EthRPCMSPoolCFGSection + ".cacheSizeERC20"))
	if err != nil {
		return nil, err
	}
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
		maxParallels:        viper.GetUint(EthRPCMSPoolCFGSection + ".maxParallels"),
		periodicRecording:   viper.GetDuration(EthRPCMSPoolCFGSection + ".periodicRecording"),
		localBackend:        localBackend,
		cacheContractTokens: cacheContractTokens,
		cacheERC20:          cacheERC20,
		cacheERC20Name:      cacheERC20Name,
	}
	lb.BaseService = *service.NewBaseService(log, "ethrpc", lb)
	return lb, nil
}

// Non-cached RPC access to get sender address for any eth transaction
func (n *MSPoolEthRPCImpl) GetTxSender(txHash common.Hash,
	blockHash common.Hash,
	txIdx uint) (common.Address, error) {
	tx, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (*types.Transaction, error) {
			tx, _, err := c.TransactionByHash(ctx, txHash)
			return tx, err
		}, nil)
	if err != nil {
		return common.Address{}, err
	}

	sender, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			return c.TransactionSender(ctx, tx, blockHash, txIdx)
		}, common.Address{})
	return sender, err
}

// Non-cached RPC access to get current block height
func (n *MSPoolEthRPCImpl) GetCurrentBlockHeight() (uint64, error) {
	return Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (uint64, error) {
			return c.BlockNumber(ctx)
		}, 0)
}

// Non-cached RPC access to get block timestamp
func (n *MSPoolEthRPCImpl) GetBlockTimestamp(height uint64) (uint64, error) {
	header, err := Do(n.pool,
		n.sem,
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
		n.sem,
		func(ctx context.Context, c *ethclient.Client) ([]types.Log, error) {
			return c.FilterLogs(ctx, fq)
		}, []types.Log{})
}

// Cached RPC access to get token sides for uniswap v2
func (n *MSPoolEthRPCImpl) GetTokensUniV2(pairContract common.Address, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	// Cache checkup
	lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{pairContract, *callopts}
	if ret, ok := n.cacheContractTokens.Get(lookupKey); ok {
		retI := ret.(itypes.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	token0, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ2pair.NewUniv2pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token0(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	token1, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ2pair.NewUniv2pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token1(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	n.cacheContractTokens.Add(lookupKey, itypes.Tuple2[common.Address, common.Address]{token0, token1})
	return token0, token1, nil
}

// Cached RPC access to get decimals for ERC20 addresses
func (n *MSPoolEthRPCImpl) GetERC20Decimals(erc20Address common.Address, callopts *bind.CallOpts) (uint8, error) {
	lookupKey := erc20Address
	if ret, ok := n.cacheERC20.Get(lookupKey); ok {
		retI := ret.(uint8)
		return retI, nil
	}

	decimals, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (uint8, error) {
			erc20, err := ERC20.NewERC20(erc20Address, c)
			if err != nil {
				return 0, nil
			}
			callopts.Context = ctx
			return erc20.Decimals(callopts)
		}, 0)

	if err != nil {
		return 0, err
	}

	n.cacheERC20.Add(lookupKey, decimals)
	return decimals, nil
}

// Non-cached RPC access to get balances for tuple (holderAddress, tokenAddress)
func (n *MSPoolEthRPCImpl) GetERC20Balances(requests []itypes.Tuple2[common.Address, common.Address],
	callopts *bind.CallOpts) ([]*big.Int, error) {
	results := []*big.Int{}

	for _, req := range requests {
		balance, err := Do(n.pool,
			n.sem,
			func(ctx context.Context, c *ethclient.Client) (*big.Int, error) {
				token, err := ERC20.NewERC20(req.Second, c)
				if err != nil {
					return big.NewInt(0), nil
				}
				callopts.Context = ctx
				return token.BalanceOf(callopts, req.First)
			}, nil)
		if err != nil {
			return results, err
		}
		results = append(results, balance)
	}
	return results, nil
}

// Cached RPC access to get name for erc20 name
func (n *MSPoolEthRPCImpl) GetERC20Name(erc20Address common.Address, callopts *bind.CallOpts) (string, error) {
	if erc20Address == common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff") {
		return "USD", nil
	}
	lookupKey := erc20Address
	if ret, ok := n.cacheERC20Name.Get(lookupKey); ok {
		retI := ret.(string)
		return retI, nil
	}

	name, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (string, error) {
			token, err := ERC20.NewERC20(erc20Address, c)
			if err != nil {
				return "unknown", err
			}
			callopts.Context = ctx
			return token.Name(callopts)
		}, "unknown")
	if err != nil {
		return "unknown", err
	}

	n.cacheERC20Name.Add(lookupKey, name)
	return name, nil
}

func (n *MSPoolEthRPCImpl) GetTokensUniV3(pairContract common.Address,
	callopts *bind.CallOpts) (common.Address, common.Address, error) {
	lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{pairContract, *callopts}
	if ret, ok := n.cacheContractTokens.Get(lookupKey); ok {
		retI := ret.(itypes.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	token0, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ3pair.NewUniv3pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token0(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	token1, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (common.Address, error) {
			pc, err := univ3pair.NewUniv3pair(pairContract, c)
			if err != nil {
				return common.Address{}, err
			}
			callopts.Context = ctx
			return pc.Token1(callopts)
		}, common.Address{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	n.cacheContractTokens.Add(lookupKey, itypes.Tuple2[common.Address, common.Address]{token0, token1})
	return token0, token1, nil
}

func (n *MSPoolEthRPCImpl) GetTokensUniV3NFT(nftContract common.Address, tokenID *big.Int, callopts *bind.CallOpts) (common.Address, common.Address, error) {
	// Cache checkup
	lookupKey := itypes.Tuple2[common.Address, bind.CallOpts]{nftContract, *callopts}
	if ret, ok := n.cacheContractTokens.Get(lookupKey); ok {
		retI := ret.(itypes.Tuple2[common.Address, common.Address])
		return retI.First, retI.Second, nil
	}

	tokens, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (itypes.Tuple2[common.Address, common.Address], error) {
			pc, err := univ3positionsnft.NewUniv3positionsnft(nftContract, c)
			if err != nil {
				return itypes.Tuple2[common.Address, common.Address]{}, err
			}
			callopts.Context = ctx
			positions, err := pc.Positions(callopts, tokenID)
			return itypes.Tuple2[common.Address, common.Address]{positions.Token0, positions.Token1}, err
		}, itypes.Tuple2[common.Address, common.Address]{})
	if err != nil {
		return common.Address{}, common.Address{}, err
	}

	n.cacheContractTokens.Add(lookupKey, tokens)
	return tokens.First, tokens.Second, nil

}

func (n *MSPoolEthRPCImpl) GetChainlinkRoundData(
	contractAddress common.Address, callopts *bind.CallOpts) (itypes.ChainlinkLatestRoundData, error) {
	roundData, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) (itypes.ChainlinkLatestRoundData, error) {
			oracle, err := chainlink.NewChainlink(contractAddress, c)
			if err != nil {
				return itypes.ChainlinkLatestRoundData{}, err
			}
			callopts.Context = ctx
			return oracle.LatestRoundData(callopts)
		}, itypes.ChainlinkLatestRoundData{})
	return roundData, err
}

// CodeAt returns the contract bytecode associated with the given account. If the account isn't a contract, returns nil.
func (n *MSPoolEthRPCImpl) IsContract(Address common.Address, callopts *bind.CallOpts) (bool, error) {
	isAddressContract, err := Do(n.pool,
		n.sem,
		func(ctx context.Context, c *ethclient.Client) ([]byte, error) {
			return c.CodeAt(ctx, Address, nil)
		}, []byte{})
	if len(isAddressContract) > 0 {
		return true, err
	}
	return false, err
}
