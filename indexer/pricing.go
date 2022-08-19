package indexer

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/supragya/gograph"
)

type OracleMap struct {
	Network        string   `json:"Network"`
	ChainID        int      `json:"ChainID"`
	StableCoinsUSD []string `json:"StableCoinsUSD"`
	Tokens         []struct {
		ID       string `json:"id"`
		Contract string `json:"contract"`
	} `json:"Tokens"`
	Oracles []struct {
		From     string `json:"from"`
		To       string `json:"to"`
		Contract string `json:"contract"`
	} `json:"Oracles"`
}

type Pricing struct {
	oracleMapsRootDir string
	diskCacheRootDir  string
	networkName       string
	oracleFile        string
	cacheFile         string
	oracleHash        string
	graph             *gograph.Graph[string, string]
	oracleMap         OracleMap
}

func GetPricingEngine() *Pricing {
	pricing := Pricing{}

	pricing.oracleMapsRootDir = viper.GetString("general.oracleMapsRootDir")
	pricing.networkName = viper.GetString("general.networkName")
	pricing.diskCacheRootDir = viper.GetString("general.diskCacheRootDir")
	pricing.oracleFile = pricing.oracleMapsRootDir + "/oraclemaps_" + pricing.networkName + ".json"

	// Open oracle file
	fd, err := os.Open(pricing.oracleFile)
	util.ENOK(err)
	defer fd.Close()
	fileBytes, err := io.ReadAll(fd)
	util.ENOK(err)

	// Load bare oracle Map
	util.ENOK(json.Unmarshal(fileBytes, &pricing.oracleMap))

	pricing.oracleHash = hex.EncodeToString(util.SHA256Hash(fileBytes))

	pricing.cacheFile = fmt.Sprintf("%s/prcg_%s_%s.bincache",
		pricing.diskCacheRootDir,
		pricing.networkName,
		pricing.oracleHash)

	if _, err := os.Stat(pricing.cacheFile); err != nil {
		// Cache doesn't exists
		log.Info("prewarming pricing graph cache")
		tempGraph := gograph.NewGraphStringUintString(false)
		for _, o := range pricing.oracleMap.Oracles {
			tempGraph.AddEdge(o.From, o.To, 1, o.Contract)
		}
		tempGraph.SaveToDisk(pricing.cacheFile)
	}

	// Load up pricing
	pricing.graph = gograph.NewGraphStringUintString(false)
	util.ENOK(pricing.graph.ReadFromDisk(pricing.cacheFile))
	log.Info("loaded up pricing graph")

	return &pricing
}
