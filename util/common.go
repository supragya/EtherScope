package util

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"
	"time"

	itypes "github.com/Blockpour/Blockpour-Geth-Indexer/types"
	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	EthErrorRegexes              []*regexp.Regexp
	ContextDeadlineExceededRegex *regexp.Regexp
	IOTimeoutRegex               *regexp.Regexp
	FailOnNonEthError            bool
	FailOnNonEthErrorSet         bool
)

type Tuple2[A any, B any] struct {
	First  A
	Second B
}

type Tuple3[A any, B any, C any] struct {
	First  A
	Second B
	Third  C
}

// Checks if error is nil or not. Kills process if not nil
func ENOK(err error) {
	ENOKS(2, err)
}

func ENOKS(skip int, err error) {
	if err != nil {
		_, file, no, ok := runtime.Caller(skip)
		if ok {
			fileSplit := strings.Split(file, "/")
			log.WithFields(log.Fields{
				"file": fileSplit[len(fileSplit)-1],
				"line": no,
			}).Fatalln(err)
		}
		log.Fatalln(err)
	}
}

func ENOKF(err error, info interface{}) {
	if err != nil {
		ENOK(fmt.Errorf("%s: %v", err.Error(), info))
	}
}

// Check if error (if any) is ethereum error
// Also takes into account boolean flag `failOnNonEthError` in cfg
// If false, silently fail and continue to next event
func IsEthErr(err error) bool {
	if !FailOnNonEthErrorSet {
		FailOnNonEthError = viper.GetBool("general.failOnNonEthError")
		FailOnNonEthErrorSet = true
	}

	if err != nil {

		// Else, actually check if known Eth error.
		e := err.Error()
		for _, r := range EthErrorRegexes {
			if r.MatchString(e) {
				return true
			}
		}

		// Everything is EthError if `failOnNonEthError` is false
		if !FailOnNonEthError {
			log.Warn("NoFail umatched: ", e)
			return true
		}
	}
	return false
}

func IsRPCCallTimedOut(err error) bool {
	return ContextDeadlineExceededRegex.MatchString(err.Error()) ||
		IOTimeoutRegex.MatchString(err.Error())
}

func IsExecutionReverted(err error) bool {
	return err != nil && EthErrorRegexes[0].MatchString(err.Error())
}

func GetUser() (*user.User, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	if os.Geteuid() == 0 {
		// Root, try to retrieve SUDO_USER if exists
		if u := os.Getenv("SUDO_USER"); u != "" {
			usr, err = user.Lookup(u)
			if err != nil {
				return nil, err
			}
		}
	}

	return usr, nil
}

func GetUserHomedir() string {
	home, err := GetUser()
	ENOK(err)
	return home.HomeDir
}

func VerifyFileExistence(file string) error {
	_, err := os.Stat(file)
	return err
}

func DivideBy10pow(num *big.Int, pow uint8) *big.Float {
	pow10 := big.NewFloat(math.Pow10(int(pow)))
	numfloat := new(big.Float).SetInt(num)
	return new(big.Float).Quo(numfloat, pow10)
}

func ExtractAddressFromLogTopic(hash common.Hash) common.Address {
	return common.BytesToAddress(hash[12:])
}

func ExtractIntFromBytes(_bytes []byte) *big.Int {
	isNeg := (_bytes[0] >> 7) == 1

	var magnitude []byte
	if isNeg {
		magnitude = GetMagnitudeForNeg(_bytes)
	} else {
		magnitude = _bytes
	}

	a := big.NewInt(0)
	a = a.SetBytes(magnitude)
	if isNeg {
		a = a.Neg(a)
	}
	return a
}

func GetMagnitudeForNeg(_bytes []byte) []byte {
	foundOne := false
	for byteIdx := len(_bytes) - 1; byteIdx >= 0; byteIdx-- {
		for bitIdx := 0; bitIdx < 8; bitIdx++ {
			if foundOne {
				// Flip
				_bytes[byteIdx] ^= 1 << bitIdx
			} else if ((_bytes[byteIdx] << (7 - bitIdx)) >> 7) == 1 {
				// Spare this one
				foundOne = true
			}
		}
	}
	return _bytes
}

func ConstructTopics(eventsToIndex []string) ([]common.Hash, error) {
	topicsList := []common.Hash{}
	for _, t := range eventsToIndex {
		switch t {
		case "UniswapV2Swap":
			topicsList = append(topicsList, itypes.UniV2SwapTopic)
		case "UniswapV2Mint":
			topicsList = append(topicsList, itypes.UniV2MintTopic)
		case "UniswapV2Burn":
			topicsList = append(topicsList, itypes.UniV2BurnTopic)
		case "UniswapV3Swap":
			topicsList = append(topicsList, itypes.UniV3SwapTopic)
		case "UniswapV3Mint":
			topicsList = append(topicsList, itypes.UniV3MintTopic)
		case "UniswapV3Burn":
			topicsList = append(topicsList, itypes.UniV3BurnTopic)
		case "Transfer":
			topicsList = append(topicsList, itypes.ERC20TransferTopic)
		default:
			return []common.Hash{}, fmt.Errorf("unknown topic for construction: %s", t)
		}
	}
	return topicsList, nil
}

func SHA256Hash(_bytes []byte) []byte {
	hasher := sha256.New()
	hasher.Write(_bytes)
	return hasher.Sum(nil)
}

func NewCtx(timeOut time.Duration) context.Context {
	ctx, _ := context.WithTimeout(context.Background(), timeOut)
	return ctx
}

func init() {
	EthErrors := []string{
		"execution reverted", // Should always be kept at idx 0
		"abi: cannot marshal",
		"no contract code at given address",
		"abi: attempting to unmarshall",
		"missing trie node",
		"no contract code at given address",
	}
	for _, e := range EthErrors {
		EthErrorRegexes = append(EthErrorRegexes, regexp.MustCompile(e))
	}

	ContextDeadlineExceededRegex = regexp.MustCompile("context deadline exceeded")
	IOTimeoutRegex = regexp.MustCompile("i/o timeout")
}
