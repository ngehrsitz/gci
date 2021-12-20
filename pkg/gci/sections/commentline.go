package sections

import (
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/constants"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"strings"
)

func init() {
	SectionParserInst.RegisterSection(CommentLineType{})
}

type CommentLineType struct{}

func (c CommentLineType) generate(sectionPrefix Section, sectionStr string, sectionSuffix Section) Section {
	match, commentParameter := sectionStrMatchesAlias(sectionStr, []string{"comment"})
	if match {
		return CommentLine(commentParameter)
	}
	return nil
}

func (c CommentLineType) helpText() string {
	return "Comment(your text here) - Prints the specified indented comment"
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
