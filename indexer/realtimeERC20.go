package indexer

import (
	"sync"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/indexer/types"
	"github.com/ethereum/go-ethereum/core/types"
)

func (r *RealtimeIndexer) processTransfer(
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
		AmountUSD:   0, // TODO
	}

	AddToSynopsis(mt, bm, transfer, items, "transfer", true)
}
