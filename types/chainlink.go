package types

import (
	"encoding/gob"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ChainlinkLatestRoundData struct {
	RoundId         *big.Int
	Answer          *big.Int
	StartedAt       *big.Int
	UpdatedAt       *big.Int
	AnsweredInRound *big.Int
}

type DirectPriceDerivationInfo struct {
	LatestRoundData ChainlinkLatestRoundData
	Decimals        uint8
	ConversionPrice *big.Float
}

type CounterpartyPriceDerivationInfo struct {
	CalculationTx     common.Hash
	LogIdx            uint
	CounterpartyToken common.Address
	CounterpartyQty   *big.Float
	CounterpartyPrice *big.Float
	SelfQty           *big.Float
}

type PriceResult struct {
	Price                     *big.Float
	IsStablecoin              bool
	IsDerivedFromCounterparty bool
	CounterpartyInfo          *CounterpartyPriceDerivationInfo
	DerivationInfo            map[string]DirectPriceDerivationInfo
}

type UniV2Metadata struct {
	Amt0 *big.Int
	Amt1 *big.Int
}

func init() {
	gob.Register(UniV2Metadata{})
	gob.Register(ChainlinkLatestRoundData{})
}
