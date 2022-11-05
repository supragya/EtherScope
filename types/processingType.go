package types

type ProcessingType uint

const (
	PricingEngineRequest ProcessingType = iota
	UserRequested
)

func (p ProcessingType) ToString() string {
	switch p {
	case PricingEngineRequest:
		return "pricing"
	case UserRequested:
		return "user"
	}
	return ""
}
