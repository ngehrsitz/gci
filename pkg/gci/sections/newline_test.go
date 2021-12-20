package sections

import (
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"testing"
)

func TestNewLineSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, NewLine{}, specificity.MisMatch{}},
		{`"x"`, NewLine{}, specificity.MisMatch{}},
		{`"\n"`, NewLine{}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
