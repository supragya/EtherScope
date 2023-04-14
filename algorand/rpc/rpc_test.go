package rpc

import (
	"fmt"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
	"github.com/algorand/go-algorand-sdk/client/v2/common/models"
	"github.com/stretchr/testify/assert"
)

const url = "https://flashy-quiet-card.algorand-mainnet.discover.quiknode.pro/288becf9ca16eb031fb0d515a208176d424a861e"
const token = "288becf9ca16eb031fb0d515a208176d424a861e"

var rpc, _ = NewAlgoRPC(url, token)

func TestGetCurrentBlockHeight(t *testing.T) {
	h, err := rpc.GetCurrentBlockHeight()
	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, h, uint64(0))
}

func TestGetBlock(t *testing.T) {
	b, err := rpc.GetBlock(24869124)
	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, b, models.Block{})
}

func TestGetBlockTransactions(t *testing.T) {
	txns, err := rpc.GetBlockTransactions(24869124)
	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, txns, []models.Transaction{})
}

func TestGetBlockTimestamp(t *testing.T) {
	ts, err := rpc.GetBlockTimestamp(24869124)
	if err != nil {
		t.Error(err)
	}

	assert.IsType(t, ts, uint64(0))
}

func TestGetAssetInfo(t *testing.T) {
	asset, err := rpc.GetAssetInfo(31566704)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, asset.Name, "USDC")
	assert.Equal(t, asset.Decimals, uint64(0))
}

func TestGetTinyManV1Transactions(t *testing.T) {
	sig := common.TinymanV1SwapSignature

	txGroups, err := rpc.GetTxGroups(24869122, []types.FunctionSignature{sig})
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Found ", len(txGroups), " transaction batches")
	for i, group := range txGroups {
		fmt.Println("Found ", len(group.Transactions), " transactions in batch", i)
	}

}

func TestGetTransactionsAlgofi(t *testing.T) {
	sigs := common.AlgoFiSwapSignatures
	round := uint64(26049535)

	txn, err := rpc.GetTxGroups(round, sigs)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%s\n", txn[0].GroupID)

}

// Test the case where there are 3 transactions in the group
func TestGetTransactionsAlgofiLen3(t *testing.T) {
	sigs := common.AlgoFiSwapSignatures
	round := uint64(26204937)

	group, err := rpc.GetTxGroups(round, sigs)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Found ", len(group), " transaction batches")

	for i, g := range group {
		fmt.Println("Found ", len(g.Transactions), " transactions in batch", i)

		for _, tx := range g.Transactions {
			fmt.Println(tx.Id)
			fmt.Println(string(tx.Group))
		}

	}
}

func TestGetTransactionsAlgofi3(t *testing.T) {
	// tests for both sef and fse signatures
	sigs := common.AlgoFiSwapSignatures
	round := uint64(26204937)

	group, err := rpc.GetTxGroups(round, sigs)
	if err != nil {
		t.Error(err)
	}

	fmt.Println("Found ", len(group), " transaction batches")

	for i, g := range group {
		fmt.Println("Found ", len(g.Transactions), " transactions in batch", i)

		for _, tx := range g.Transactions {
			fmt.Println(tx.Id)
		}

	}
}

func TestGetTinymanV1PoolInfo(t *testing.T) {
	algousdc := "FPOU46NBKTWUZCNMNQNXRWNW3SMPOOK4ZJIN5WSILCWP662ANJLTXVRUKA"
	a0, a1, err := rpc.GetTinymanV1PoolInfo(algousdc)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, a0, 31566704)
	assert.Equal(t, a1, 0)
}

func TestGetTinymanV2PoolInfo(t *testing.T) {
	algousdcv2 := "2PIFZW53RHCSFSYMCFUBW4XOCXOMB7XOYQSQ6KGT3KVGJTL4HM6COZRNMM"
	a0, a1, err := rpc.GetTinymanV2PoolInfo(algousdcv2)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, a0, uint64(31566704))
	assert.Equal(t, a1, uint64(0))
}
func TestGetAlgoFiPoolInfo(t *testing.T) {
	a0, a1, err := rpc.GetAlgoFiPoolInfo(605929989)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, a0, uint64(31566704))
	assert.Equal(t, a1, uint64(0))
}

// STBL2/ALGO Pool
func TestGetAlgoFiPoolInfoAlgofiPoolReserves(t *testing.T) {
	a0, a1, err := rpc.GetAlgoFiPoolReserves(855716333)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(a0)
	fmt.Println(a1)

	// assert.Equal(t, a0, uint64(31566704))
	// assert.Equal(t, a1, uint64(0))
}
