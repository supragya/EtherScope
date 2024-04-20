package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/supragya/EtherScope/libs/config"
	"github.com/supragya/EtherScope/services/ethrpc"
	localbackend "github.com/supragya/EtherScope/services/local_backend"
	"github.com/supragya/EtherScope/services/node"
	"github.com/supragya/EtherScope/services/oraclenode"
	outputsink "github.com/supragya/EtherScope/services/output_sink"
	"github.com/supragya/EtherScope/version"
	"github.com/spf13/cobra"
)

// RootCmd represents the base command when called without any subcommands
var ConfigGen = &cobra.Command{
	Use:   "configgen",
	Short: "generate config for escope",
	Long:  `generate config for escope`,
	Run:   GenConfig,
}

func GenConfig(cmd *cobra.Command, args []string) {
	// config string generation
	// file := cfgFile
	content := ""

	// Add metadata
	content += "# Config file generated from escope configgen\n"

	content += "\n# VERSION INFO:\n"
	for _, line := range version.GetVersionStrings() {
		content += "# " + line + "\n"
	}

	// Add statutory note
	content += "\n# NOTE: Starting with v0.5.0 (enterprise), escope has moved to\n"
	content += "# service based architecture and as such config sections reflect\n"
	content += "# such separation. These configs may evolve over time and as such\n"
	content += "# it is adviced to always use `escope configgen` as the utility\n"
	content += "# to generate example configs.\n"
	content += "\n"

	// Services config
	content += sectionGen(node.NodeCFGSection, node.NodeCFGNecessity,
		node.NodeCFGHeader, node.NodeCFGFields[:])
	content += sectionGen(oraclenode.OracleNodeCFGSection, oraclenode.OracleNodeCFGNecessity,
		oraclenode.OracleNodeCFGHeader, oraclenode.OracleNodeCFGFields[:])
	content += sectionGen(localbackend.BadgerCFGSection, localbackend.BadgerCFGNecessity,
		localbackend.BadgerCFGHeader, localbackend.BadgerCFGFields[:])
	content += sectionGen(outputsink.RabbitMQCFGSection, outputsink.RabbitMQCFGNecessity,
		outputsink.RabbitMQCFGHeader, outputsink.RabbitMQCFGFields[:])
	content += sectionGen(ethrpc.EthRPCMSPoolCFGSection, ethrpc.EthRPCMSPoolCFGNecessity,
		ethrpc.EthRPCMSPoolCFGHeader, ethrpc.EthRPCMSPoolCFGFields[:])

	if err := os.WriteFile(cfgFile, []byte(content), 0600); err != nil {
		panic(err)
	}
}

func sectionGen(section, necessity string, header []string, fields []config.Field) string {
	content := ""
	content += "# +" + strings.Repeat("-", 78) + "+\n"
	content += fmt.Sprintf("# | %-15s: %-60s|\n", "SERVICE", section)
	content += fmt.Sprintf("# | %-15s: %-60s|\n", "NECESSITY", necessity)
	isFirst := true
	for _, line := range header {
		if isFirst {
			content += fmt.Sprintf("# | %-15s: %-60s|\n", "INFO", line)
			isFirst = false
		} else {
			content += fmt.Sprintf("# | %-15s  %-60s|\n", "", line)
		}
	}
	content += "# +" + strings.Repeat("-", 78) + "+\n"

	content += section + ":\n"
	for _, field := range fields {
		content += fmt.Sprintf("#   %-5s: %-20s\n", "FIELD", field.Name)
		content += fmt.Sprintf("#   %-5s: %-20s\n", "TYPE", field.Type)
		isFirst := true
		for _, line := range field.Info {
			if isFirst {
				content += fmt.Sprintf("#   %-5s: %-60s\n", "INFO", line)
				isFirst = false
			} else {
				content += fmt.Sprintf("#   %-5s  %-60s\n", "", line)
			}
		}
		content += fmt.Sprintf("    %s: %v\n", field.Name, field.Default)
		content += "\n"
	}
	content += "\n"

	return content
}
