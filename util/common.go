package util

import (
	"math"
	"math/big"
	"os"
	"os/user"
	"regexp"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"
)

var (
	AbiErrRegex *regexp.Regexp
	NoContract  *regexp.Regexp
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

// Check if error (if any) is ethereum error
func IsEthErr(err error) bool {
	if err != nil {
		e := err.Error()
		if e == "execution reverted" ||
			AbiErrRegex.MatchString(e) ||
			NoContract.MatchString(e) {
			return true
		}
	}
	return false
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

func init() {
	AbiErrRegex = regexp.MustCompile("abi: cannot marshal.*")
	NoContract = regexp.MustCompile("no contract code at given address")
}
