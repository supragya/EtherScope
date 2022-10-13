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
		isRestricted    = viper.GetBool("erc20transfer.restrictToWhitelist")
		restrictionType = viper.GetString("erc20transfer.restrictionType")
		whitelistFile   = viper.GetString("erc20transfer.whitelistFile")
	)

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

	var whitelistFrom, whitelistsTo *map[common.Address]bool = nil, nil

	if restrictionType == "to" || restrictionType == "both" {
		whitelistsTo = &whitelistMap
	}
	if restrictionType == "from" || restrictionType == "both" {
		whitelistFrom = &whitelistMap
	}

	return &ERC20TransferRestrictions{isRestricted, whitelistFrom, whitelistsTo}
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

	if r.erc20TransferRestrictions != nil {
		if r.erc20TransferRestrictions.isRestricted {
			if r.erc20TransferRestrictions.whitelistFrom != nil {
				if _, ok := (*r.erc20TransferRestrictions.whitelistFrom)[sender]; !ok {
					return
				}
			}
			if r.erc20TransferRestrictions.whitelistTo != nil {
				if _, ok := (*r.erc20TransferRestrictions.whitelistTo)[recv]; !ok {
					return
				}
			}
		}
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
		Type:        "transfer",
		Network:     r.dbconn.ChainID,
		LogIdx:      l.Index,
		Transaction: l.TxHash,
		Time:        bm.Time,
		Height:      l.BlockNumber,
		Token:       l.Address,
		Sender:      sender,
		TxSender:    txSender,
		Receiver:    recv,
		Amount:      formattedAmount,
		AmountUSD:   tokenPrice,
	}

	AddToSynopsis(mt, bm, transfer, items, "transfer", true)
	instrumentation.TfrProcessed.Inc()
}
