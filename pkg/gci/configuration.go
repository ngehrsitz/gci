package gci

import (
	"github.com/daixiang0/gci/pkg/configuration"
	sectionsPkg "github.com/daixiang0/gci/pkg/gci/sections"
)

type GciConfiguration struct {
	configuration.FormatterConfiguration
	Sections []sectionsPkg.Section
}
