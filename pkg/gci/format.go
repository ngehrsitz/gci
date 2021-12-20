package gci

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/daixiang0/gci/pkg/constants"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"strings"
)

// Formats the import section of a Go file
func formatGoFile(input []byte, cfg GciConfiguration) ([]byte, error) {
	startIndex := bytes.Index(input, []byte(constants.ImportStartFlag))
	if startIndex < 0 {
		return nil, errors.New("No import statement present in File!")
	}
	endIndexFromStart := bytes.Index(input[startIndex:], []byte(constants.ImportEndFlag))
	if endIndexFromStart < 0 {
		return nil, errors.New("Import statement not closed!")
	}
	endIndex := startIndex + endIndexFromStart

	unformattedImports := input[startIndex+len(constants.ImportStartFlag) : endIndex]
	formattedImports, err := formatImportBlock(unformattedImports, cfg)
	if err != nil {
		return nil, err
	}

	output := []byte{}
	output = append(output, input[:startIndex+len(constants.ImportStartFlag)]...)
	output = append(output, formattedImports...)
	output = append(output, input[endIndex+1:]...)
	return output, nil
}

// Takes unsorted imports as byte array and formats them according to the specified sections
func formatImportBlock(input []byte, cfg GciConfiguration) ([]byte, error) {
	//strings.ReplaceAll(input, "\r\n", linebreak)
	lines := strings.Split(string(input), constants.Linebreak)
	imports, err := parseToImportDefinitions(lines)
	if err != nil {
		return nil, fmt.Errorf("An error occured while trying to parse imports: %w", err)
	}
	// create mapping between sections and imports
	sectionMap := make(map[sectionsPkg.Section][]importPkg.ImportDef, len(cfg.Sections))
	// find matching section for every importSpec
	for _, i := range imports {
		// determine match specificity for every available section
		var bestSection sectionsPkg.Section
		var bestSectionSpecificity specificity.MatchSpecificity = specificity.MisMatch{}
		for _, section := range cfg.Sections {
			sectionSpecificity := section.MatchSpecificity(i)
			if sectionSpecificity.IsMoreSpecific(specificity.MisMatch{}) && sectionSpecificity.Equal(bestSectionSpecificity) {
				// specificity is identical
				return nil, fmt.Errorf("Import %s matched section %s and %s equally", i, bestSection, section)
			}
			if sectionSpecificity.IsMoreSpecific(bestSectionSpecificity) {
				// better match found
				bestSectionSpecificity = sectionSpecificity
				bestSection = section
			}
		}
		if bestSection == nil {
			return nil, fmt.Errorf("No section found for Import: %v", i)
		}
		sectionMap[bestSection] = append(sectionMap[bestSection], i)
	}
	// generate output by formatting the sections
	var output string
	for _, section := range cfg.Sections {
		output += section.Format(sectionMap[section], cfg.FormatterConfiguration)
	}
	return []byte(output), nil
}
