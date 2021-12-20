package sections

import (
	"github.com/daixiang0/gci/pkg/configuration"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func init() {
	SectionParserInst.RegisterSection(DefaultSectionType{})
}

type DefaultSectionType struct{}

func (d DefaultSectionType) generate(sectionPrefix Section, sectionStr string, sectionSuffix Section) Section {
	match, _ := sectionStrMatchesAlias(sectionStr, []string{"def", "default"})
	if match {
		return DefaultSection{sectionPrefix, sectionSuffix}
	}
	return nil
}

func (d DefaultSectionType) helpText() string {
	return "Def|Default - Contains all imports that could not be matched to another section type"
}

type DefaultSection struct {
	Prefix Section
	Suffix Section
}

func (d DefaultSection) sectionPrefix() Section {
	return d.Prefix
}

func (d DefaultSection) sectionSuffix() Section {
	return d.Suffix
}

func (d DefaultSection) MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity {
	return specificity.Default{}
}

func (d DefaultSection) Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	return inorderSectionFormat(d, imports, cfg)
}

func (d DefaultSection) String() string {
	return "Default"
}
