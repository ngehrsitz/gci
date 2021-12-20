package sections

import (
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"testing"
)

func TestPrefixSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`"foo/pkg/bar"`, Prefix{"", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"foo", nil, nil}, specificity.Match{3}},
		{`"foo/pkg/bar"`, Prefix{"bar", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"github.com/foo/bar", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"github.com/foo", nil, nil}, specificity.MisMatch{}},
		{`"foo/pkg/bar"`, Prefix{"github.com/bar", nil, nil}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
