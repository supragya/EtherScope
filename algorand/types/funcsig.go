package types

import "fmt"

// for the case where there is no appId (ex. algofi where each pool has a different app id),
// we only use the key to identify the function and we set AppId to 0)
type FunctionSignature struct {
	AppId uint64
	Key   string
}

func (fs FunctionSignature) String() string {
	return fmt.Sprintf("%d-%s", fs.AppId, fs.Key)
}

func (fs FunctionSignature) Equals(other FunctionSignature) bool {
	if fs.AppId == 0 || other.AppId == 0 {
		return fs.Key == other.Key
	} else {
		return fs.AppId == other.AppId && fs.Key == other.Key
	}
}
