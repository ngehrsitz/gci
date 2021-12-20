package sections

import (
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/constants"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func init() {
	SectionParserInst.RegisterSection(NewLineType{})
}

type NewLineType struct{}

func (n NewLineType) generate(sectionPrefix Section, sectionStr string, sectionSuffix Section) Section {
	match, _ := sectionStrMatchesAlias(sectionStr, []string{"nl", "newline"})
	if match {
		return NewLine{}
	}
	return nil
}

func (n NewLineType) helpText() string {
	return "NL|NewLine - Prints an empty line"
}

type NewLine struct{}

func (n NewLine) sectionPrefix() Section {
	return nil
}

func (n NewLine) sectionSuffix() Section {
	return nil
}

func (n NewLine) MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity {
	return specificity.MisMatch{}
}

func (n NewLine) Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	return constants.Linebreak
}

func (n NewLine) String() string {
	return "NewLine"
}
