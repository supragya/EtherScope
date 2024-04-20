package main

import (
	"github.com/supragya/EtherScope/cmd"
	"github.com/supragya/EtherScope/libs/util"
)

func main() {
	util.ENOK(cmd.RootCmd.Execute())
}
