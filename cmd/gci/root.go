package gci

import (
	"fmt"
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/gci"
	"github.com/spf13/cobra"
	"os"
)

var (
	diffMode   *bool
	writeMode  *bool
	localFlags *[]string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gci [-d | -w] [-local localPackageURLs] path ...",
	Short: "Gci controls golang package import order and makes it always deterministic",
	Long: `Gci enables automatic formatting of imports in a deterministic manner

If you want to integrate this as part of your CI take a look at golangci-lint.`,
	ValidArgsFunction: goFileArg,
	Args:              cobra.MinimumNArgs(1),
	Version:           "0.5",
	RunE:              runInCompatibilityMode,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	//emulate the old CLI
	diffMode = rootCmd.Flags().BoolP("diff", "d", false, "display diffs instead of rewriting files")
	writeMode = rootCmd.Flags().BoolP("write", "w", false, "write result to (source) file instead of stdout")
	localFlags = rootCmd.Flags().StringSliceP("local", "l", []string{}, "put imports beginning with this string after 3rd-party packages, separate imports by comma")
}

func runInCompatibilityMode(cmd *cobra.Command, args []string) error {
	// Workaround since the Deprecation message in Cobra can not be printed to STDERR
	fmt.Fprintln(os.Stderr, "Using the old parameters is deprecated, please use the named subcommands!")

	if *writeMode && *diffMode {
		return fmt.Errorf("diff and write must not be specified at the same time")
	}
	// generate section specification from old localFlags format
	sections := gci.LocalFlagsToSections(*localFlags)
	cfg := gci.GciConfiguration{configuration.FormatterConfiguration{false, false, false}, sections}
	if *writeMode {
		return gci.WriteFormattedFiles(args, cfg)
	}
	if *diffMode {
		return gci.DiffFormattedFiles(args, cfg)
	}
	return gci.PrintFormattedFiles(args, cfg)
}
