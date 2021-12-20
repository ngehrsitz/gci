package sections

import (
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"testing"
)

func TestDefaultSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, DefaultSection{}, specificity.Default{}},
		{`"x"`, DefaultSection{}, specificity.Default{}},
	}
	testSpecificity(t, testCases)
}
