package gci

import (
	"github.com/daixiang0/gci/pkg/configuration"
	"github.com/daixiang0/gci/pkg/io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func isTestInputFile(file os.FileInfo) bool {
	return !file.IsDir() && strings.HasSuffix(file.Name(), ".in.go")
}

func testGciConfig() GciConfiguration {
	sections := LocalFlagsToSections([]string{
		"github.com/daixiang0",
		"github.com/local",
	})
	return GciConfiguration{configuration.FormatterConfiguration{false, false, false}, sections}
}

func TestRun(t *testing.T) {
	testFiles, err := io.FindFilesForPath("internal/testdata", isTestInputFile)
	if err != nil {
		t.Fatal(err)
	}
	for _, testFile := range testFiles {
		fileBaseName := strings.TrimSuffix(testFile, ".in.go")
		t.Run(fileBaseName, func(t *testing.T) {
			t.Parallel()

			_, formattedFile, err := LoadFormatGoFile(fileBaseName+".in.go", testGciConfig())
			if err != nil {
				t.Fatal(err)
			}
			expectedOutput, err := ioutil.ReadFile(fileBaseName + ".out.go")
			if err != nil {
				t.Fatal(err)
			}
			assert.Equal(t, string(formattedFile), string(expectedOutput), "output")
			assert.NoError(t, err)
		})
	}
}
