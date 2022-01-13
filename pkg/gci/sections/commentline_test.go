package sections

import (
	"github.com/daixiang0/gci/pkg/gci/specificity"
	"testing"
)

func TestCommentLineSpecificity(t *testing.T) {
	testCases := []specificityTestData{
		{`""`, CommentLine(""), specificity.MisMatch{}},
		{`"x"`, CommentLine(""), specificity.MisMatch{}},
		{`"//"`, CommentLine(""), specificity.MisMatch{}},
		{`"/"`, CommentLine(""), specificity.MisMatch{}},
	}
	testSpecificity(t, testCases)
}

func TestCommentLineGenerator(t *testing.T) {
	testCases := []sectionTestData{
		{"commentline", CommentLine("")},
		{"Commentline(abc)", CommentLine("abc")},
	}
	testGenerator(t, testCases)
}
