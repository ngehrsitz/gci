package gci

import (
	"fmt"
	"github.com/daixiang0/gci/pkg/constants"
	importPkg "github.com/daixiang0/gci/pkg/gci/imports"
	"strings"
)

// Recursively parses import lines into a list of ImportDefs
func parseToImportDefinitions(unformattedLines []string) ([]importPkg.ImportDef, error) {
	newImport := importPkg.ImportDef{}
	for index, unformattedLine := range unformattedLines {
		line := strings.TrimSpace(unformattedLine)
		if line == "" {
			//empty line --> starts a new import
			return parseToImportDefinitions(unformattedLines[index+1:])
		}
		if strings.HasPrefix(line, constants.CommentFlag) {
			// comment line
			newImport.PrefixComment = append(newImport.PrefixComment, line)
			continue
		}
		// split inline comment from import
		importSegments := strings.SplitN(line, constants.CommentFlag, 2)
		switch len(importSegments) {
		case 1:
			// no inline comment
		case 2:
			// inline comment present
			newImport.InlineComment = constants.CommentFlag + importSegments[1]
		default:
			return nil, fmt.Errorf("seperating the inline comment from the import yielded an invalid number of segments: %v", importSegments)
		}
		// split alias from path
		pkgArray := strings.Fields(importSegments[0])
		switch len(pkgArray) {
		case 1:
			// only a path
			newImport.QuotedPath = pkgArray[0]
		case 2:
			// alias + path
			newImport.Alias = pkgArray[0]
			newImport.QuotedPath = pkgArray[1]
		default:
			return nil, fmt.Errorf("seperating the alias from the path yielded an invalid number of segments: %v", pkgArray)
		}
		err := newImport.Validate()
		if err != nil {
			return nil, err
		}
		followingImports, err := parseToImportDefinitions(unformattedLines[index+1:])
		return append([]importPkg.ImportDef{newImport}, followingImports...), err
	}
	return nil, nil
}
