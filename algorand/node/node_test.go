package node

import (
	"fmt"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/logger"
	types "github.com/Blockpour/Blockpour-Geth-Indexer/algorand/types"
	util "github.com/Blockpour/Blockpour-Geth-Indexer/algorand/util"
	"github.com/stretchr/testify/assert"
)

const url = "https://flashy-quiet-card.algorand-mainnet.discover.quiknode.pro/288becf9ca16eb031fb0d515a208176d424a861e"
const token = "288becf9ca16eb031fb0d515a208176d424a861e"

// Basic Tinyman V1 swap (USDC/ALGO)
func TestProcessBlocksTinymanV1(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	funcsig := common.TinymanV1SwapSignature

	f := types.TxFilter{
		StartRound: 24869123,
		EndRound:   24869126,
		Signatures: []types.FunctionSignature{funcsig},
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.IsType(t, swaps, []types.Swap{})
}

// Tinyman V1.1 seems to structure transactions same as V1
func TestProcessBlocksTinymanV11(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	funcsig := common.TinymanV1SwapSignature

	f := types.TxFilter{
		StartRound: 26349674,
		EndRound:   26349674,
		Signatures: []types.FunctionSignature{funcsig},
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.IsType(t, swaps, []types.Swap{})
}

// Normal tinyman V2 swap (USDC/ALGO)
func TestProcessBlocksWithTinyManV2USDCALGO(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	funcsig := common.TinymanV2SwapSignature

	f := types.TxFilter{
		StartRound: 26201363,
		EndRound:   26201363,
		Signatures: []types.FunctionSignature{funcsig},
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.IsType(t, swaps, []types.Swap{})
}

// Normal tinyman V2 swap (USDC/USDT)
func TestProcessBlocksWithTinyManV2USDCUSDT(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sig := common.TinymanV2SwapSignature
	f := types.TxFilter{
		StartRound: 26168163,
		EndRound:   26168163,
		Signatures: []types.FunctionSignature{sig},
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.IsType(t, swaps, []types.Swap{})
}

// Basic algofi swap (the block contains 2 swaps of which we only test the first one)
func TestProcessBlocksAlgoFiALGOUSDC(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.AlgoFiSwapSignatures

	f := types.TxFilter{
		StartRound: 26361185,
		EndRound:   26361185,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(len(swaps))

	for _, swap := range swaps {
		swap.Print()
	}

	assert.IsType(t, swaps, []types.Swap{})
	assert.Equal(t, 2, len(swaps))
	assert.Equal(t, 0, swaps[0].Token0)
	assert.Equal(t, common.USDC.Id, swaps[0].Token1)
	assert.InDelta(t, -3.93, swaps[0].Amount0, 0.1)
	assert.InDelta(t, 0.82, swaps[0].Amount1, 0.1)
	assert.Equal(t, "ALGO", swaps[0].Name0)
	assert.Equal(t, "USDC", swaps[0].Name1)
}

// This tests for an Algofi transaction that contains a nanoswap (stablecoin to stablecoin swap)
func TestProcessBlocksAlgoFiNanoSwap(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.AlgoFiSwapSignatures

	f := types.TxFilter{
		StartRound: 26204937,
		EndRound:   26204937,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	assert.IsType(t, swaps, []types.Swap{})
	assert.Equal(t, 2, len(swaps))
	assert.Equal(t, common.USDC.Id, swaps[0].Token0)
	assert.Equal(t, common.STBL.Id, swaps[0].Token1)
	assert.InDelta(t, 75.3, swaps[0].Amount0, 0.1)
	assert.InDelta(t, -75, swaps[0].Amount1, 0.1)
	assert.Equal(t, "USDC", swaps[0].Name0)
	assert.Equal(t, "STBL", swaps[0].Name1)
}

/*
Test for a batch with multiple swaps on different exchanges
algofi + tinyman: kDHSaenH9Qam6yWgFQsBME/+zQH6UUH1dYey35RI3fw=
https://algoexplorer.io/tx/group/kDHSaenH9Qam6yWgFQsBME%2F%2BzQH6UUH1dYey35RI3fw%3D
We test we can parse the algofi swap within the transaction batch
*/
func TestProcessBlocksAlgoFiBatchSwap(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.AlgoFiSwapSignatures

	f := types.TxFilter{
		StartRound: 26210445,
		EndRound:   26210445,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	assert.IsType(t, swaps, []types.Swap{})
	assert.Equal(t, 1, len(swaps))
	assert.Equal(t, common.ALGO.Id, swaps[0].Token0)
	assert.Equal(t, common.PDAT.Id, swaps[0].Token1)
	assert.InDelta(t, -1.7, swaps[0].Amount0, 0.1)
	assert.InDelta(t, 3230306.3, swaps[0].Amount1, 0.1)
	assert.Equal(t, "ALGO", swaps[0].Name0)
	assert.Equal(t, "PDAT", swaps[0].Name1)
	assert.Equal(t, "VYZPVANXQMKKZPT3YTQUXWPCG4G33RMKKLV2OJB6A4O6MRV2KSLA", swaps[0].Transaction)
}

/*
https://algoexplorer.io/tx/group/UIwqecoB6a4IyXPyCZTQJ31%2BgwJDZXA6NuxaaXpbITQ%3D
Test for an algofi swap between 2 other swaps.
*/
func TestProcessBlocksAlgoFi6(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.AlgoFiSwapSignatures

	f := types.TxFilter{
		StartRound: 26224309,
		EndRound:   26224309,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	// for _, swap := range swaps {
	// 	swap.Print()
	// }

	assert.IsType(t, swaps, []types.Swap{})
	assert.Equal(t, 2, len(swaps))
	assert.Equal(t, common.ALGO.Id, swaps[0].Token0)
	assert.Equal(t, common.USDt.Id, swaps[0].Token1)
	assert.InDelta(t, -2.39025, swaps[0].Amount0, 0.1)
	assert.InDelta(t, 0.501907, swaps[0].Amount1, 0.1)
	assert.Equal(t, "ALGO", swaps[0].Name0)
	assert.Equal(t, "USDt", swaps[0].Name1)
	assert.Equal(t, "EFCSATXFR72G3CVXXNQ5YJEHH66O4ANQFP23QNKUYMRENLKUEWMQ", swaps[0].Transaction)
	assert.Equal(t, "algofi-swap", swaps[0].Type)
}

/*
Test getting swaps for all tinymanv1, tinymanv2, algofi.
Currently errors because of too many requests to the node (limit of 25/s i believe)
*/
func TestProcessBlocks(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.SupportedFunctionSignatures

	f := types.TxFilter{
		StartRound: 26224088,
		EndRound:   26224088,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.IsType(t, swaps, []types.Swap{})
}

/*
This tests a tinymanv2 edge case where the app tx has an additional inner transaction
which contains the change.
https://algoexplorer.io/tx/group/C0zfF2p28hmI5gJblvFLDWQLMa98FhvPCiSZ1qALMb8%3D
*/
func TestTinymanV2EdgeCase1(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.SupportedFunctionSignatures

	f := types.TxFilter{
		StartRound: 26318342,
		EndRound:   26318342,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.Equal(t, 1, len(swaps))
	assert.IsType(t, swaps, []types.Swap{})
	assert.Equal(t, common.USDC.Id, swaps[0].Token0)
	assert.Equal(t, common.ALGO.Id, swaps[0].Token1)
	assert.Equal(t, "USDC", swaps[0].Name0)
	assert.Equal(t, "ALGO", swaps[0].Name1)
	assert.InDelta(t, -230.92, swaps[0].Amount0, 0.1)
	assert.InDelta(t, 1000, swaps[0].Amount1, 0.1)
	assert.IsType(t, swaps, []types.Swap{})
	assert.Equal(t, "NY653STL65EQDLI4FYZZMW5WARPYYD36BJTNKDLAIGVLFI3QCRQQ", swaps[0].Transaction)
	assert.Equal(t, "tinyman-v2-swap", swaps[0].Type)
}

/*
edge case for a group of 2 algofi swaps.
One with an extra transfer and one with nanoswap exchange
-> https://algoexplorer.io/tx/group/ZlVJHnC%2FNux8Ot4bNZB3oCd1Pdghf1uCSj4TIocZbR4%3D
*/
func TestAlgoFiEdgeCase1(t *testing.T) {
	logger := logger.NewNopLogger()
	indexer, err := NewNode(0, url, token, 10, true, "", logger)
	util.ENOK(err)

	sigs := common.AlgoFiSwapSignatures

	f := types.TxFilter{
		StartRound: 26204937,
		EndRound:   26204937,
		Signatures: sigs,
	}

	swaps, err := indexer.ProcessTxns(f)
	if err != nil {
		fmt.Println(err)
	}

	for _, swap := range swaps {
		swap.Print()
	}

	assert.Equal(t, 2, len(swaps))
	assert.IsType(t, swaps, []types.Swap{})

	assert.Equal(t, common.USDC.Id, swaps[0].Token0)
	assert.Equal(t, common.STBL.Id, swaps[0].Token1)
	assert.Equal(t, "USDC", swaps[0].Name0)
	assert.Equal(t, "STBL", swaps[0].Name1)
	assert.InDelta(t, 75.3, swaps[0].Amount0, 0.1)
	assert.InDelta(t, -75, swaps[0].Amount1, 0.1)
	assert.Equal(t, "W2IZ3EHDRW2IQNPC33CI2CXSLMFCFICVKQVWIYLJWXCTD765RW47ONNCEY", swaps[0].Sender)
	assert.Equal(t, "W2IZ3EHDRW2IQNPC33CI2CXSLMFCFICVKQVWIYLJWXCTD765RW47ONNCEY", swaps[0].Receiver)
	assert.Equal(t, "52LP7NPYKPXJWFV7QOVAJDB4W4FAWSG2C5G7YNB7K4Q3YBXXIJUA", swaps[0].Transaction)
	assert.Equal(t, "algofi-swap", swaps[0].Type)

	assert.Equal(t, common.ALGO.Id, swaps[1].Token0)
	assert.Equal(t, common.STBL.Id, swaps[1].Token1)
	assert.Equal(t, "ALGO", swaps[1].Name0)
	assert.Equal(t, "STBL", swaps[1].Name1)
	assert.InDelta(t, -358.608025, swaps[1].Amount0, 0.1)
	assert.InDelta(t, 75, swaps[1].Amount1, 0.1)
	assert.Equal(t, "W2IZ3EHDRW2IQNPC33CI2CXSLMFCFICVKQVWIYLJWXCTD765RW47ONNCEY", swaps[1].Sender)
	assert.Equal(t, "W2IZ3EHDRW2IQNPC33CI2CXSLMFCFICVKQVWIYLJWXCTD765RW47ONNCEY", swaps[1].Receiver)
	assert.Equal(t, "HA5MRLVAOCH7HTIIEFVLR26RC2OP7MDXRJEH3SDB2HD6RN5VNP3Q", swaps[1].Transaction)
	assert.Equal(t, "algofi-swap", swaps[1].Type)
}
