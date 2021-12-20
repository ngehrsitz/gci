package gci

import (
	"fmt"
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/constants"
	"github.com/daixiang0/gci/pkg/gci"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"
	"github.com/spf13/cobra"
)

var (
	noInlineComments *bool
	noPrefixComments *bool
	debug            *bool
	sections         *[]string
)

func goFileArg(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{"go"}, cobra.ShellCompDirectiveFilterFileExt
}

func newGciCommand(use, short, long string, aliases []string, runFunc func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := cobra.Command{
		Use:               use,
		Aliases:           aliases,
		Short:             short,
		Long:              long,
		ValidArgsFunction: goFileArg,
		Args:              cobra.MinimumNArgs(1),
		RunE:              runFunc,
	}

	// register command as subcommand
	rootCmd.AddCommand(&cmd)

	sectionHelp := "Sections define how inputs will be processed. " +
		"Section names are case-insensitive and may contain parameters in (). " +
		fmt.Sprintf("A section can contain a Prefix and a Suffix section which is delimited by %q. ", constants.SectionSeparator) +
		"These sections can be used for formatting and will only be rendered if the main section contains an entry." +
		"\n" +
		sectionsPkg.SectionParserInst.SectionHelpTexts()
	// add flags
	noInlineComments = cmd.Flags().Bool("NoInlineComments", false, "Drops inline comments while formatting")
	noPrefixComments = cmd.Flags().Bool("NoPrefixComments", false, "Drops comment lines above an import statement while formatting")
	sections = cmd.Flags().StringSliceP("Section", "s", []string{"Standard", "NewLine:Default"}, sectionHelp)
	return &cmd
}

func newDebuggableGciCommand(use, short, long string, aliases []string, runFunc func(cmd *cobra.Command, args []string) error) *cobra.Command {
	cmd := newGciCommand(use, short, long, aliases, runFunc)
	debug = cmd.Flags().BoolP("debug", "d", false, "Enables debug output from the formatter")
	return cmd
}

func parseGciConfiguration() (*gci.GciConfiguration, error) {
	parsedSections := []sectionsPkg.Section{}
	for _, sectionStr := range *sections {
		section, err := sectionsPkg.SectionParserInst.ParseStrToSection(sectionStr)
		if err != nil {
			return nil, err
		}
		parsedSections = append(parsedSections, section)
	}
	return &gci.GciConfiguration{configuration.FormatterConfiguration{*noInlineComments, *noPrefixComments, *debug}, parsedSections}, nil
}
