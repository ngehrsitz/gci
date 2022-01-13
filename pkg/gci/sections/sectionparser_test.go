package sections

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type sectionTestData struct {
	sectionDef      string
	expectedSection Section
}

func testGenerator(t *testing.T, testCases []sectionTestData) {
	for _, test := range testCases {
		parsedSection, _ := SectionParserInst.ParseStrToSection(test.sectionDef)
		assert.Equal(t, test.expectedSection, parsedSection)
	}
}
