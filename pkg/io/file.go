package io

import (
	"io/fs"
	"os"
	"path/filepath"
)

type fileCheckFunction func(file os.FileInfo) bool

func FindFilesForDirectory(dirPath string, fileCheckFun fileCheckFunction) ([]string, error) {
	var filePaths []string
	filepath.WalkDir(dirPath, func(path string, entry fs.DirEntry, err error) error {
		file, err := entry.Info()
		if !entry.IsDir() && fileCheckFun(file) {
			filePaths = append(filePaths, filepath.Clean(path))
		}
		return nil
	})
	return filePaths, nil
}

func FindFilesForPath(path string, fileCheckFun fileCheckFunction) ([]string, error) {
	switch entry, err := os.Stat(path); {
	case err != nil:
		return nil, err
	case entry.IsDir():
		return FindFilesForDirectory(path, fileCheckFun)
	case fileCheckFun(entry):
		return []string{filepath.Clean(path)}, nil
	default:
		return []string{}, nil
	}
}

func FindGoFilesForPath(path string) ([]string, error) {
	return FindFilesForPath(path, isGoFile)
}

func isGoFile(file os.FileInfo) bool {
	return !file.IsDir() && filepath.Ext(file.Name()) == ".go"
}
