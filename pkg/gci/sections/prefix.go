package sections

import (
	"fmt"
	"github.com/daixiang0/gci/pkg/configuration"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"strings"
)

func init() {
	SectionParserInst.RegisterSection(PrefixType{})
}

type PrefixType struct{}

func (p PrefixType) generate(sectionPrefix Section, sectionStr string, sectionSuffix Section) Section {
	match, prefixParameter := sectionStrMatchesAlias(sectionStr, []string{"prefix", "importprefix"})
	if match {
		return Prefix{prefixParameter, sectionPrefix, sectionSuffix}
	}
	return nil
}

func (p PrefixType) helpText() string {
	return "Prefix(gitlab.com/myorg) - Groups all imports with the specified Prefix. Imports will be matched to the longest Prefix."
}

type Prefix struct {
	ImportPrefix string
	Prefix       Section
	Suffix       Section
}

func (p Prefix) sectionPrefix() Section {
	return p.Prefix
}

func (p Prefix) sectionSuffix() Section {
	return p.Suffix
}

func (p Prefix) MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity {
	if len(p.ImportPrefix) > 0 && strings.HasPrefix(spec.Path(), p.ImportPrefix) {
		return specificity.Match{len(p.ImportPrefix)}
	}
	return specificity.MisMatch{}
}

func (p Prefix) Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	return inorderSectionFormat(p, imports, cfg)
}

func (p Prefix) String() string {
	return fmt.Sprintf("ImportPrefix(%q)", p.ImportPrefix)
}
