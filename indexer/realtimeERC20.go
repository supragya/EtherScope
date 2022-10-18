package indexer

import (
	"encoding/json"
	"io/ioutil"
	"math/big"
	"os"
	"sync"

	"github.com/Blockpour/Blockpour-Geth-Indexer/instrumentation"
	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
)

func setupERC20TransferRestrictions(events []common.Hash) *ERC20TransferRestrictions {
	isERC20TransferToBeIndexed := false
	for _, e := range events {
		if e == itypes.ERC20TransferTopic {
			isERC20TransferToBeIndexed = true
		}
	}

	if !isERC20TransferToBeIndexed {
		return nil
	}

	var (
		restrictionType = viper.GetString("erc20transfer.restrictionType")
		whitelistFile   = viper.GetString("erc20transfer.whitelistFile")
	)

	var _type ERC20RestrictionType
	switch restrictionType {
	case "none":
		_type = None
	case "to":
		_type = To
	case "from":
		_type = From
	case "both":
		_type = Both
	case "either":
		_type = Either
	default:
		panic("unknown ERC20RestrictionType")
	}

	file, err := os.Open(whitelistFile)
	util.ENOK(err)

	_bytes, err := ioutil.ReadAll(file)
	util.ENOK(err)

	whitelist := []common.Address{}
	util.ENOK(json.Unmarshal(_bytes, &whitelist))

	whitelistMap := make(map[common.Address]bool, len(whitelist))
	for _, ra := range whitelist {
		whitelistMap[ra] = true
	}

	return &ERC20TransferRestrictions{_type, &whitelistMap}
}

func restrictIndexing(r *ERC20TransferRestrictions, from common.Address, to common.Address) bool {
	if r._type == None {
		return true
	}
	var (
		whFrom = false
		whTo   = false
	)
	if _, ok := (*r.whitelist)[from]; ok {
		whFrom = true
	}
	if _, ok := (*r.whitelist)[to]; ok {
		whTo = true
	}
	switch r._type {
	case None:
		return true
	case To:
		return whTo
	case From:
		return whFrom
	case Both:
		return whTo && whFrom
	case Either:
		return whTo || whFrom
	}
	return false
}

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

	if restrictIndexing(r.erc20TransferRestrictions, sender, recv) {
		return
	}

	callopts := GetBlockCallOpts(l.BlockNumber)

	ok, formattedAmount := r.GetFormattedAmount(amt, callopts, l.Address)
	if !ok {
		return
	}

	tokenPrice := r.da.GetRateForBlock(callopts, util.Tuple2[common.Address, *big.Float]{l.Address, formattedAmount})

	txSender, err := r.da.GetTxSender(l.TxHash, l.BlockHash, l.TxIndex)
	if util.IsEthErr(err) {
		return
	}
	util.ENOK(err)

	transfer := itypes.Transfer{
		Type:                "erc20transfer",
		Network:             r.dbconn.ChainID,
		LogIdx:              l.Index,
		Transaction:         l.TxHash,
		Time:                bm.Time,
		Height:              l.BlockNumber,
		Token:               l.Address,
		Sender:              sender,
		TxSender:            txSender,
		Receiver:            recv,
		Amount:              formattedAmount,
		AmountUSD:           tokenPrice.Price,
		PriceDerivationMeta: tokenPrice,
	}

	AddToSynopsis(mt, bm, transfer, items, "transfer", true)
	instrumentation.TfrProcessed.Inc()
}
