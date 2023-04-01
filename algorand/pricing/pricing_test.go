package pricing

import (
	"fmt"
	"testing"

	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/common"
	"github.com/Blockpour/Blockpour-Geth-Indexer/algorand/rpc"
	"github.com/stretchr/testify/assert"
)

const url = "https://flashy-quiet-card.algorand-mainnet.discover.quiknode.pro/288becf9ca16eb031fb0d515a208176d424a861e"
const token = "288becf9ca16eb031fb0d515a208176d424a861e"

var r, _ = rpc.NewAlgoRPC(url, token)
var engine = NewPricingEngine(r)

func TestGetQuotePrice(t *testing.T) {
	p, err := engine.GetQuotePrice(common.ALGO)
	if err != nil {
		t.Error(err)
	}

	fmt.Print("preices", p)

	assert.Greater(t, p, float64(0.01))
	assert.Greater(t, p, float64(0.4))
}

func TestGetQuotePriceUSDC(t *testing.T) {
	p, err := engine.GetQuotePrice(common.USDC)
	if err != nil {
		t.Error(err)
	}

	fmt.Print("prices", p)
}
