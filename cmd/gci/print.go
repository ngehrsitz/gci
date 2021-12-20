package gci

import (
	"github.com/daixiang0/gci/pkg/gci"

	"github.com/spf13/cobra"
)

// printCmd represents the print command
var printCmd = newGciCommand(
	"print",
	"Outputs the formatted file to STDOUT",
	"Print outputs the formatted file. If you want to apply the changes to a file use write instead!",
	[]string{"output"},
	func(cmd *cobra.Command, args []string) error {
		gciCfg, err := parseGciConfiguration()
		if err != nil {
			return err
		}
		return gci.PrintFormattedFiles(args, *gciCfg)
	},
)
