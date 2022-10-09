package indexer

import (
	"math/big"
	"sync"

	"github.com/Blockpour/Blockpour-Geth-Indexer/instrumentation"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

func (r *RealtimeIndexer) processERC20Transfer(
	l types.Log,
	items *[]interface{},
	bm *itypes.BlockSynopsis,
	mt *sync.Mutex,
) {
	ok, sender, recv, amt := InfoTransfer(l)

	if !ok {
		return
	}

	callopts := GetBlockCallOpts(l.BlockNumber)
	ok, formattedAmount := r.GetFormattedAmount(amt, callopts, l.Address)

	if !ok {
		return
	}

	tokenPrice := r.da.GetRateForBlock(callopts, util.Tuple2[common.Address, *big.Float]{l.Address, formattedAmount})

	transfer := itypes.Transfer{
		Type:        "transfer",
		Network:     r.dbconn.ChainID,
		LogIdx:      l.Index,
		Transaction: l.TxHash,
		Time:        bm.Time,
		Height:      l.BlockNumber,
		Token:       l.Address,
		Sender:      sender,
		Receiver:    recv,
		Amount:      formattedAmount,
		AmountUSD:   tokenPrice,
	}

	AddToSynopsis(mt, bm, transfer, items, "transfer", true)
	instrumentation.TfrProcessed.Inc()
}
