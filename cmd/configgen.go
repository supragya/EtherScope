package cmd

import (
	"fmt"
	"os"

	"github.com/Blockpour/Blockpour-Geth-Indexer/services/node"
	"github.com/Blockpour/Blockpour-Geth-Indexer/version"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var ConfigGen = &cobra.Command{
	Use:   "configgen",
	Short: "generate config for bgidx",
	Long:  `generate config for bgidx`,
	Run:   GenConfig,
}

func GenConfig(cmd *cobra.Command, args []string) {
	// config string generation
	// file := cfgFile
	content := ""

	// Add metadata
	content += "# Config file generated from bgidx configgen\n"

	content += "\n# VERSION INFO:\n"
	for _, line := range version.GetVersionStrings() {
		content += "# " + line + "\n"
	}

	// Add statutory note
	content += "\n# NOTE: Starting with v0.5.0 (enterprise), bgidx has moved to\n"
	content += "# service based architecture and as such config sections reflect\n"
	content += "# such separation. These configs may evolve over time and as such\n"
	content += "# it is adviced to always use `bgidx configgen` as the utility\n"
	content += "# to generate example configs.\n"
	content += "\n"

	// Services config
	content += fmt.Sprintf("# ++++++++++++++++++++++ BEGIN: %s\n", node.NodeSection)
	for _, line := range node.NodeCFGHeader {
		content += fmt.Sprintf("# %s\n", line)
	}
	content += node.NodeSection + ":\n"
	for _, field := range node.NodeCFGFields {
		content += fmt.Sprintf("#   FIELD: %s\n", field.Name)
		content += fmt.Sprintf("#   EXPECTED_TYPE: %s\n", field.Type)
		for _, line := range field.Err {
			content += "#   " + line + "\n"
		}
		content += fmt.Sprintf("    %s: %v\n", field.Name, field.Default)
		content += "\n"
	}
	content += fmt.Sprintf("# ++++++++++++++++++++++ END: %s\n", node.NodeSection)

	if err := os.WriteFile(cfgFile, []byte(content), 0600); err != nil {
		panic(err)
	}
}
