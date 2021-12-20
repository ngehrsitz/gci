package gci

import (
	"github.com/daixiang0/gci/pkg/gci"

	"github.com/spf13/cobra"
)

// diffCmd represents the diff command
var diffCmd = newDebuggableGciCommand(
	"diff",
	"Prints a git style diff to STDOUT",
	"Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.",
	[]string{},
	func(cmd *cobra.Command, args []string) error {
		gciCfg, err := parseGciConfiguration()
		if err != nil {
			return err
		}
		return gci.DiffFormattedFiles(args, *gciCfg)
	},
)
