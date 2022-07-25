package util

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	EthErrorRegexes []*regexp.Regexp
	// ExecutionReverted    *regexp.Regexp
	// AbiErrRegex          *regexp.Regexp
	// NoContract           *regexp.Regexp
	// ErrUnmarshal         *regexp.Regexp
	FailOnNonEthError    bool
	FailOnNonEthErrorSet bool
)

// Checks if error is nil or not. Kills process if not nil
func ENOK(err error) {
	if err != nil {
		_, file, no, ok := runtime.Caller(1)
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
		ENOK(errors.New(fmt.Sprintf("%s: %v", err.Error(), info)))
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

func ExtractUintFromBytes(_bytes []byte) *big.Int {
	a := big.NewInt(0)
	a = a.SetBytes(_bytes)
	return a
}

func init() {
	EthErrors := []string{
		"execution reverted", // Should always be kept at idx 0
		"abi: cannot marshal",
		"no contract code at given address",
		"abi: attempting to unmarshall",
		"missing trie node",
	}
	for _, e := range EthErrors {
		EthErrorRegexes = append(EthErrorRegexes, regexp.MustCompile(e))
	}
}
