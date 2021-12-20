package gci

import (
	"fmt"
	"github.com/daixiang0/gci/pkg/io"
	"github.com/hexops/gotextdiff"
	"github.com/hexops/gotextdiff/myers"
	"github.com/hexops/gotextdiff/span"
	"golang.org/x/sync/errgroup"
	"io/ioutil"
	"os"

	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"
)

func DefaultSections() []sectionsPkg.Section {
	return []sectionsPkg.Section{sectionsPkg.StandardPackage{}, sectionsPkg.DefaultSection{sectionsPkg.NewLine{}, nil}}
}

func LocalFlagsToSections(localFlags []string) []sectionsPkg.Section {
	sections := DefaultSections()
	// Add all local arguments as ImportPrefix sections
	for _, prefix := range localFlags {
		sections = append(sections, sectionsPkg.Prefix{prefix, sectionsPkg.NewLine{}, nil})
	}
	return sections
}

func PrintFormattedFiles(paths []string, cfg GciConfiguration) error {
	return processFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		fmt.Print(string(formattedFile))
		// TODO error wrapping for difference
		return nil
	})
}

func WriteFormattedFiles(paths []string, cfg GciConfiguration) error {
	return processFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		return os.WriteFile(filePath, formattedFile, 0644)
	})
}

func DiffFormattedFiles(paths []string, cfg GciConfiguration) error {
	return processFilesInPaths(paths, cfg, func(filePath string, unmodifiedFile, formattedFile []byte) error {
		fileURI := span.URIFromPath(filePath)
		edits := myers.ComputeEdits(fileURI, string(unmodifiedFile), string(formattedFile))
		unifiedEdits := gotextdiff.ToUnified(filePath, filePath, string(unmodifiedFile), edits)
		fmt.Printf("%v", unifiedEdits)
		return nil
	})
}

type filePostFormattingFunc func(filePath string, unmodifiedFile, formattedFile []byte) error

func processFilesInPaths(paths []string, cfg GciConfiguration, fileFunc filePostFormattingFunc) error {
	var taskGroup errgroup.Group
	for _, path := range paths {
		files, err := io.FindGoFilesForPath(path)
		if err != nil {
			return err
		}
		for _, filePath := range files {
			// run file processing in parallel
			taskGroup.Go(func() error {
				unmodifiedFile, formattedFile, err := LoadFormatGoFile(filePath, cfg)
				if err != nil {
					return err
				}
				return fileFunc(filePath, unmodifiedFile, formattedFile)
			})
		}
	}
	return taskGroup.Wait()
}

func LoadFormatGoFile(filePath string, cfg GciConfiguration) (unmodifiedFile, formattedFile []byte, err error) {
	unmodifiedFile, err = ioutil.ReadFile(filePath)
	if err != nil {
		return nil, nil, err
	}

	formattedFile, err = formatGoFile(unmodifiedFile, cfg)
	if err != nil {
		return unmodifiedFile, nil, err
	}
	return unmodifiedFile, formattedFile, nil
}
