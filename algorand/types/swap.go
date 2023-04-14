package types

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/util"
)

type Swap struct {
	Type         string
	Network      int
	Transaction  string
	Height       uint64
	Time         uint64
	Sender       string
	TxSender     string
	Receiver     string
	PairContract string
	Token0       string
	Token1       string
	Amount0      float64
	Amount1      float64
	Name0        string
	Name1        string
	AmountUSD    float64
	Reserve0     float64
	Reserve1     float64
	Price0       *PriceResult
	Price1       *PriceResult
	Decimals0    uint64
	Decimals1    uint64
	GroupId      string // somehow can't manage to decode group id.
	PoolAppId    uint64
	Protocol     string
}

/*
Implementing to match swap structure of 0.5.0. Personally think
it would be better to flatten the swap structure.
*/
type PriceResult struct {
	Price float64
	Path  []interface{}
}

type OrderedSwaps map[int][]Swap

// will be modified. was trying to get a human readable output.
// don't want to change the stringer interface but i want an interface to
// get some sort of summarized human readable output.  Essentially not printing
// FeeTx, ReceiverTx.
func (s Swap) Print() {
	util.PrintJSON(map[string]interface{}{
		"type":      s.Type,
		"txId":      s.Transaction,
		"height":    s.Height,
		"time":      s.Time,
		"sender":    s.Sender,
		"receiver":  s.Receiver,
		"pair":      s.PairContract,
		"name0":     s.Name0,
		"name1":     s.Name1,
		"amount0":   s.Amount0,
		"amount1":   s.Amount1,
		"address0":  s.Token0,
		"address1":  s.Token1,
		"price0":    s.Price0,
		"price1":    s.Price1,
		"decimals0": s.Decimals0,
		"decimals1": s.Decimals1,
		"poolAppId": s.PoolAppId,
		"protocol":  s.Protocol,
	})
}
