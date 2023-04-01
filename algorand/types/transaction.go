package types

import "github.com/algorand/go-algorand-sdk/client/v2/common/models"

type TxGroup struct {
	GroupID           []byte
	Transactions      []models.Transaction
	FunctionSignature FunctionSignature
}

type TxFilter struct {
	StartRound uint64
	EndRound   uint64
	Signatures []FunctionSignature
}
