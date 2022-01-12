package sections

import (
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/constants"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"strings"
)

func init() {
	commentLineType := &SectionType{
		generatorFun: func(parameter string, sectionPrefix, sectionSuffix Section) (Section, error) {
			return CommentLine(parameter), nil
		},
		aliases:       []string{"Comment", "CommentLine"},
		parameterHelp: "your text here",
		description:   "Prints the specified indented comment",
	}
	SectionParserInst.RegisterSection(commentLineType.standAloneSection())
}

type CommentLine string

func (c CommentLine) MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity {
	return specificity.MisMatch{}
}

func (c CommentLine) Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	comment := string(c)
	if !strings.HasSuffix(comment, constants.Linebreak) {
		comment += constants.Linebreak
	}
	return comment
}

func (c CommentLine) sectionPrefix() Section {
	return nil
}

func (c CommentLine) sectionSuffix() Section {
	return nil
}

func (c CommentLine) String() string {
	return "CommentLine"
}
