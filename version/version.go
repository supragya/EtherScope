package version

import (
	"bytes"
	"strconv"
)

// Codename -- hardcoded
var ApplicationCodename string = "dubai"

// Build commit -- supplied compile time
var buildCommit string = "unknowncommit"

// Build time -- supplied compile time
var buildTime string = "unknowntime"

// Build time -- supplied compile time
var builder string = "unknownbuilder"

// Go version -- supplied compile time
var gover string = "unknownver"

// Persistence version -- database compatibility index.
// NOT TO be supplied compile time. Should be hardcoded.
var PersistenceVersion uint8 = 3

var RootCmdVersion string = prepareVersionString()

func prepareVersionString() string {
	var buffer bytes.Buffer

	buffer.WriteString(buildCommit + " persistence v" + strconv.Itoa(int(PersistenceVersion)))
	buffer.WriteString(" (" + ApplicationCodename + ")")

	buffer.WriteString("\ncompiled at " + buildTime + " by " + builder)
	buffer.WriteString(" using " + gover)

	return buffer.String()
}
