package priceresolver

import (
	"encoding/csv"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type DexRecords []DexRecord

// Implements interface sort.Sort
func (r DexRecords) Len() int {
	return len(r)
}

// Implements interface sort.Sort
// sorts by StartBlock
func (r DexRecords) Less(i, j int) bool {
	return r[i].StartBlock < r[j].StartBlock
}

// Implements interface sort.Sort
func (r DexRecords) Swap(i, j int) {
	tempRec := r[i]
	r[i] = r[j]
	r[j] = tempRec
}

type DexRecord struct {
	Pair       common.Address
	Token0     common.Address
	Token1     common.Address
	StartBlock int64
	Exchange   common.Address
}

func loadDexCSV(filePath string) (DexRecords, error) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	DexRecords, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	resDexRecords := []DexRecord{}

	for _, rec := range DexRecords[1:] {
		startBlock, err := strconv.Atoi(strings.Split(rec[4], "-")[0])
		if err != nil {
			panic(err)
		}
		DexRecord := DexRecord{
			Pair:       common.HexToAddress(strings.Split(rec[0], "-")[0]),
			Token0:     common.HexToAddress(strings.Split(rec[1], "-")[0]),
			Token1:     common.HexToAddress(strings.Split(rec[2], "-")[0]),
			StartBlock: int64(startBlock),
			Exchange:   common.HexToAddress(strings.Split(rec[6], "-")[0]),
		}
		resDexRecords = append(resDexRecords, DexRecord)
	}

	return resDexRecords, nil
}
