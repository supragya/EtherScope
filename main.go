package main

import (
	"github.com/Blockpour/Blockpour-Geth-Indexer/cmd"
	"github.com/Blockpour/Blockpour-Geth-Indexer/libs/util"
)

func main() {
	util.ENOK(cmd.RootCmd.Execute())
}
