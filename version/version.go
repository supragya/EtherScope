package version

import (
	"bytes"
	"fmt"
)

// Codename -- hardcoded
var ApplicationCodename string = "enterprise"

// Build commit -- supplied compile time
var buildCommit string = "unknowncommit"

const UNTAGGED_GITTAG = "untagged"

var gittag string = UNTAGGED_GITTAG

// Build time -- supplied compile time
var buildTime string = "unknowntime"

// Build time -- supplied compile time
var builder string = "unknownbuilder"

// Go version -- supplied compile time
var gover string = "unknownver"

// Persistence version -- database compatibility index.
// NOT TO be supplied compile time. Should be hardcoded.
var PersistenceVersion uint8 = 5

var RootCmdVersion string = prepareVersionString()

func GetVersionStrings() []string {
	return []string{
		fmt.Sprintf("%s persistence v%d (%s) tagged %s", buildCommit, PersistenceVersion, ApplicationCodename, gittag),
		fmt.Sprintf("compiled at %s by %s using %s", buildTime, builder, gover),
	}
}

func prepareVersionString() string {
	var buffer bytes.Buffer

	firstLine := true
	for _, line := range GetVersionStrings() {
		if firstLine {
			buffer.WriteString(line)
			firstLine = false
		}
		buffer.WriteString("\n" + line)
	}

	return buffer.String()
}

func GetGitTag() string {
	return gittag
}
