package main

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/cmd"
	"github.com/Blockpour/Blockpour-Geth-Indexer/logger"
	"github.com/Blockpour/Blockpour-Geth-Indexer/util"
)

func main() {
	util.ENOK(cmd.RootCmd.Execute())
}

func init() {
	logger.SetupLog()
}
