package version

import (
	"bytes"
)

// Application version  -- supplied compile time
var ApplicationVersion string = "unknownversion"

// Build commit -- supplied compile time
var buildCommit string = "unknowncommit"

// Build time -- supplied compile time
var buildTime string = "unknowntime"

// Build time -- supplied compile time
var builder string = "unknownbuilder"

var RootCmdVersion string = prepareVersionString()

func prepareVersionString() string {
	var buffer bytes.Buffer
	buffer.WriteString(ApplicationVersion + " build " + buildCommit)
	buffer.WriteString("\ncompiled at " + buildTime + " by " + builder)
	return buffer.String()
}
