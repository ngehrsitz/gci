package sections

import (
	"github.com/daixiang0/gci/pkg/configuration"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"strings"
)

// A SectionType is used to dynamically register Sections with the parser
type SectionType interface {
	generate(sectionPrefix Section, sectionStr string, sectionSuffix Section) Section
	helpText() string
}

// Section defines a part of the formatted output.
type Section interface {
	// MatchSpecificity returns how well an Import matches to this Section
	MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity
	// Format receives the array of imports that have matched this section and formats them according to it´s rules
	Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string
	// Returns the Section that will be prefixed if this section is rendered
	sectionPrefix() Section
	// Returns the Section that will be suffixed if this section is rendered
	sectionSuffix() Section
	// Implement Stringer interface
	String() string
}

//Unbound methods that are required until interface methods are supported

// Default method for formatting a section
func inorderSectionFormat(section Section, imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	imports = importPkg.SortImportsByPath(imports)
	var output string
	if len(imports) > 0 && section.sectionPrefix() != nil {
		// imports are not passed to a prefix section to prevent rendering them twice
		output += section.sectionPrefix().Format([]importPkg.ImportDef{}, cfg)
	}
	for _, importDef := range imports {
		output += importDef.Format(cfg)
	}
	if len(imports) > 0 && section.sectionSuffix() != nil {
		// imports are not passed to a suffix section to prevent rendering them twice
		output += section.sectionSuffix().Format([]importPkg.ImportDef{}, cfg)
	}
	return output
}

// Method used for checking if the string representation of the section matches the expected aliases
// Parameters in () after the alias are also parsed
func sectionStrMatchesAlias(str string, sectionAliases []string) (match bool, parameters string) {
	lowerCaseStr := strings.ToLower(str)
	for _, alias := range sectionAliases {
		// str starts with alias
		if strings.HasPrefix(lowerCaseStr, alias) {
			leftoverStr := lowerCaseStr[len(alias):]
			if leftoverStr == "" {
				return true, ""
			}
			if strings.HasPrefix(leftoverStr, "(") && strings.HasSuffix(leftoverStr, ")") {
				return true, strings.TrimSuffix(strings.TrimPrefix(leftoverStr, "("), ")")
			}
		}
	}
	return false, ""
}
