package analyzer

import (
	"bytes"
	"fmt"
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/gci"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"strings"
)

const (
	NoInlineCommentsFlag = "noInlineComments"
	NoPrefixCommentsFlag = "noPrefixComments"
	SectionsFlag         = "Sections"
	SectionSeperator     = ";"
)

var (
	noInlineComments bool
	noPrefixComments bool
	sectionsStr      string
)

func init() {
	Analyzer.Flags.BoolVar(&noInlineComments, NoInlineCommentsFlag, false, "If comments in the same line as the input should be present")
	Analyzer.Flags.BoolVar(&noPrefixComments, NoPrefixCommentsFlag, false, "If comments above an input should be present")
	Analyzer.Flags.StringVar(&sectionsStr, SectionsFlag, "", "Specify the Sections format that should be used to check the file formatting")
}

var Analyzer = &analysis.Analyzer{
	Name:     "gci",
	Doc:      "A tool that control golang package import order and make it always deterministic.",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      runAnalysis,
}

func runAnalysis(pass *analysis.Pass) (interface{}, error) {
	// TODO input validation

	goFileReferences := make(map[string]*token.File, len(pass.Files))
	// pass does not supply a list of all go files
	for _, pkgFile := range pass.Files {
		// find pkg in the list of all files
		pass.Fset.Iterate(func(file *token.File) bool {
			if file.Base() == int(pkgFile.Package) {
				goFileReferences[file.Name()] = file
			}
			return true
		})
	}
	expectedNumFiles := len(pass.Files)
	foundNumFiles := len(goFileReferences)
	if expectedNumFiles != foundNumFiles {
		return nil, fmt.Errorf("Expected %d files in Analyzer input, Found %d", expectedNumFiles, foundNumFiles)
	}

	// read configuration options
	gciCfg, err := parseGciConfiguration()
	if err != nil {
		return nil, err
	}

	for filePath, file := range goFileReferences {
		unmodifiedFile, formattedFile, err := gci.LoadFormatGoFile(filePath, *gciCfg)
		if err != nil {
			return nil, err
		}
		// search for a difference
		fileRunes := bytes.Runes(unmodifiedFile)
		formattedRunes := bytes.Runes(formattedFile)
		diffIdx := compareRunes(fileRunes, formattedRunes)
		switch diffIdx {
		case -1:
			// no difference
		default:
			diffPos := file.Position(file.Pos(diffIdx))
			// prevent invalid access to array
			fileRune := "nil"
			formattedRune := "nil"
			if len(fileRunes)-1 >= diffIdx {
				fileRune = fmt.Sprintf("%q", fileRunes[diffIdx])
			}
			if len(formattedRunes)-1 >= diffIdx {
				formattedRune = fmt.Sprintf("%q", formattedRunes[diffIdx])
			}
			pass.Reportf(file.Pos(diffIdx), "Expected %s, Found %s at %s[line %d,col %d]", formattedRune, fileRune, filePath, diffPos.Line, diffPos.Column)
		}
	}
	return nil, nil
}

func compareRunes(a, b []rune) (differencePos int) {
	// check shorter rune slice first to prevent invalid array access
	shorterRune := a
	if len(b) < len(a) {
		shorterRune = b
	}
	// check for differences up to where the length is identical
	for idx := 0; idx < len(shorterRune); idx++ {
		if a[idx] != b[idx] {
			return idx
		}
	}
	// check that we have compared two equally long rune arrays
	if len(a) != len(b) {
		return len(shorterRune) + 1
	}
	return -1
}

func parseGciConfiguration() (*gci.GciConfiguration, error) {
	parsedSections := []sectionsPkg.Section{}
	if sectionsStr == "" {
		parsedSections = gci.DefaultSections()
	} else {
		for _, sectionStr := range strings.Split(sectionsStr, SectionSeperator) {
			section, err := sectionsPkg.SectionParserInst.ParseStrToSection(sectionStr)
			if err != nil {
				return nil, err
			}
			parsedSections = append(parsedSections, section)
		}
	}
	return &gci.GciConfiguration{configuration.FormatterConfiguration{noInlineComments, noPrefixComments, false}, parsedSections}, nil
}
