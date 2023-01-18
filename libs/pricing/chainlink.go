package priceresolver

import (
	"encoding/csv"
	"log"
	"os"
	"sort"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

type ChainlinkRecords []ChainlinkRecord

// Implements interface sort.Sort
func (r ChainlinkRecords) Len() int {
	return len(r)
}

// Implements interface sort.Sort
// sorts by StartBlock
func (r ChainlinkRecords) Less(i, j int) bool {
	return r[i].StartBlock < r[j].StartBlock
}

// Implements interface sort.Sort
func (r ChainlinkRecords) Swap(i, j int) {
	tempRec := r[i]
	r[i] = r[j]
	r[j] = tempRec
}

type ChainlinkRecord struct {
	From       common.Address
	To         common.Address
	Oracle     common.Address
	StartBlock int64
}

var (
	USDTokenID = common.HexToAddress("0xffffffffffffffffffffffffffffffffffffffff")
)

func loadChainlinkCSV(filePath string) (ChainlinkRecords, error) {
	recs, err := loadChainlinkCSVi(filePath)
	if err != nil {
		return ChainlinkRecords{}, err
	}
	sort.Sort(recs)
	return recs, nil
}

func loadChainlinkCSVi(filePath string) (ChainlinkRecords, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	ChainlinkRecords, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	resChainlinkRecords := []ChainlinkRecord{}

	for _, rec := range ChainlinkRecords[1:] {
		startBlock, err := strconv.Atoi(rec[3])
		if err != nil {
			panic(err)
		}
		if rec[1] == "USD" {
			rec[1] = "0xffffffffffffffffffffffffffffffffffffffff"
		}
		ChainlinkRecord := ChainlinkRecord{
			From:       common.HexToAddress(rec[0]),
			To:         common.HexToAddress(rec[1]),
			Oracle:     common.HexToAddress(rec[2]),
			StartBlock: int64(startBlock),
		}
		resChainlinkRecords = append(resChainlinkRecords, ChainlinkRecord)
	}

	return resChainlinkRecords, nil
}
