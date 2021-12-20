package gci

import (
	"github.com/daixiang0/gci/pkg/gci"

	"github.com/spf13/cobra"
)

// writeCmd represents the write command
var writeCmd = newGciCommand(
	"write",
	"Formats the specified files in-place",
	"Write modifies the specified files in-place",
	[]string{"overwrite"},
	func(cmd *cobra.Command, args []string) error {
		gciCfg, err := parseGciConfiguration()
		if err != nil {
			return err
		}
		return gci.WriteFormattedFiles(args, *gciCfg)
	},
)
