package sections

import (
	"github.com/daixiang0/gci/pkg/configuration"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"github.com/daixiang0/gci/pkg/gci/specificity"
)

func init() {
	SectionParserInst.RegisterSection(StandardPackageType{})
}

type StandardPackageType struct{}

func (s StandardPackageType) generate(sectionPrefix Section, sectionStr string, sectionSuffix Section) Section {
	match, _ := sectionStrMatchesAlias(sectionStr, []string{"std", "standard"})
	if match {
		return StandardPackage{sectionPrefix, sectionSuffix}
	}
	return nil
}

func (s StandardPackageType) helpText() string {
	return "Std|Standard - Captures all standard packages if they do not match another section"
}

type StandardPackage struct {
	prefix Section
	suffix Section
}

func (s StandardPackage) sectionPrefix() Section {
	return s.prefix
}

func (s StandardPackage) sectionSuffix() Section {
	return s.suffix
}

func (s StandardPackage) MatchSpecificity(spec importPkg.ImportDef) specificity.MatchSpecificity {
	if isStandardPackage(spec.Path()) {
		return specificity.StandardPackageMatch{}
	}
	return specificity.MisMatch{}
}

func (s StandardPackage) Format(imports []importPkg.ImportDef, cfg configuration.FormatterConfiguration) string {
	return inorderSectionFormat(s, imports, cfg)
}

func (s StandardPackage) String() string {
	return "Standard"
}

func isStandardPackage(pkg string) bool {
	_, ok := standardPackages[pkg]
	return ok
}
