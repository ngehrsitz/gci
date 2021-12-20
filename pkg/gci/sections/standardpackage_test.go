package sections

import (
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"testing"
)

func TestStandardSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`"context"`, StandardPackage{}, specificity.StandardPackageMatch{}},
		{`"contexts"`, StandardPackage{}, specificity.MisMatch{}},
		{`"crypto"`, StandardPackage{}, specificity.StandardPackageMatch{}},
		{`"crypto1"`, StandardPackage{}, specificity.MisMatch{}},
		{`"crypto/ae"`, StandardPackage{}, specificity.MisMatch{}},
		{`"crypto/aes"`, StandardPackage{}, specificity.StandardPackageMatch{}},
		{`"crypto/aes2"`, StandardPackage{}, specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}
