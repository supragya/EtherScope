package version

import (
	"bytes"
)

// Application version  -- supplied compile time
var ApplicationVersion string = "unknownversion"
var ApplicationCodename string = "buffalo"

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
var PersistenceVersion uint8 = 1

var RootCmdVersion string = prepareVersionString()

func prepareVersionString() string {
	var buffer bytes.Buffer
	buffer.WriteString(ApplicationVersion + " build " + buildCommit + "(" + ApplicationCodename + ")")
	buffer.WriteString("\ncompiled at " + buildTime + " by " + builder)
	buffer.WriteString("\npersistence version " + string(PersistenceVersion))
	buffer.WriteString("\nusing " + gover)
	return buffer.String()
}
