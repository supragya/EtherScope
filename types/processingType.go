package types

type ProcessingType uint

const (
	PricingEngineRequest ProcessingType = iota
	UserRequested
)
