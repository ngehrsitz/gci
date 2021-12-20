package imports

import (
	"errors"
	"fmt"
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/constants"
	"sort"
	"strings"
	"unicode"
)

type ImportDef struct {
	Alias         string
	QuotedPath    string
	PrefixComment []string
	InlineComment string
}

func (i ImportDef) Path() string {
	return strings.TrimSuffix(strings.TrimPrefix(i.QuotedPath, string('"')), string('"'))
}

// Checks whether the contents are valid for an import
func (i ImportDef) Validate() error {
	err := checkAlias(i.Alias)
	if err != nil {
		return err
	}
	if !strings.HasPrefix(i.QuotedPath, string('"')) {
		return errors.New("Path is missing starting quotes!")
	}
	if !strings.HasSuffix(i.QuotedPath, string('"')) {
		return errors.New("Path is missing closing quotes!")
	}
	return nil
}

func checkAlias(alias string) error {
	for _, r := range alias {
		if !unicode.IsLetter(r) && r != '_' {
			return fmt.Errorf("Found non-letter character %q in Alias: %s", r, alias)
		}
	}
	return nil
}
func (i ImportDef) String() string {
	return i.Format(configuration.FormatterConfiguration{false, false, false})
}

func (i ImportDef) Format(cfg configuration.FormatterConfiguration) string {
	lineprefix := constants.Indent
	var output string
	if cfg.NoPrefixComments == false {
		for _, prefixComment := range i.PrefixComment {
			output += lineprefix + prefixComment + constants.Linebreak
		}
	}
	output += lineprefix
	if i.Alias != "" {
		output += i.Alias + constants.Blank
	}
	output += i.QuotedPath
	if cfg.NoInlineComments == false {
		if i.InlineComment != "" {
			output += constants.Blank + i.InlineComment
		}
	}
	output += constants.Linebreak
	return output
}

func SortImportsByPath(imports []ImportDef) []ImportDef {
	sort.Slice(
		imports,
		func(i, j int) bool {
			return sort.StringsAreSorted([]string{imports[i].Path(), imports[j].Path()})
		},
	)
	return imports
}
