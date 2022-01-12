package sections

import (
	"errors"
	"github.com/daixiang0/gci/pkg/configuration"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
)

// A SectionType is used to dynamically register Sections with the parser
type SectionType struct {
	generatorFun  func(parameter string, sectionPrefix, sectionSuffix Section) (Section, error)
	aliases       []string
	parameterHelp string
	description   string
}

func (t *SectionType) withoutParameter() *SectionType {
	sectionGeneratorFun := t.generatorFun
	t.generatorFun = func(parameter string, sectionPrefix, sectionSuffix Section) (Section, error) {
		if parameter != "" {
			return nil, errors.New("Default section does not take a parameter")
		}
		return sectionGeneratorFun(parameter, sectionPrefix, sectionSuffix)
	}
	t.parameterHelp = ""
	return t
}

func (t *SectionType) standAloneSection() *SectionType {
	nextFun := t.generatorFun
	t.generatorFun = func(parameter string, sectionPrefix, sectionSuffix Section) (Section, error) {
		if sectionPrefix != nil {
			return nil, errors.New("Section may not contain a Prefix")
		}
		if sectionSuffix != nil {
			return nil, errors.New("Section may not contain a Suffix")
		}
		return nextFun(parameter, sectionPrefix, sectionSuffix)
	}
	return t
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
