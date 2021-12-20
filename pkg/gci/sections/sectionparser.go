package sections

import (
	"fmt"
	"github.com/daixiang0/gci/pkg/constants"
	"strings"
)

var SectionParserInst = SectionParser{}

type SectionParser struct {
	sectionTypes []SectionType
}

func (s *SectionParser) RegisterSection(sectionType SectionType) {
	s.sectionTypes = append(s.sectionTypes, sectionType)
}

func (s *SectionParser) ParseStrToSection(sectionStr string) (Section, error) {
	trimmedSection := strings.TrimSpace(sectionStr)
	sectionSegments := strings.Split(trimmedSection, constants.SectionSeparator)
	switch len(sectionSegments) {
	case 1:
		// section
		return s.parseSectionComponentsToSection("", sectionSegments[0], "")
	case 2:
		// prefix + section
		return s.parseSectionComponentsToSection(sectionSegments[0], sectionSegments[1], "")
	case 3:
		// prefix + section + suffix
		return s.parseSectionComponentsToSection(sectionSegments[0], sectionSegments[1], sectionSegments[2])
	}
	return nil, fmt.Errorf("Expected section Definition in format [FormattingSection:]Section[:FormattingSection], got %s", sectionSegments)

}

func (s *SectionParser) parseSectionComponentsToSection(sectionPrefixStr string, sectionStr string, sectionSuffixStr string) (Section, error) {
	var sectionPrefix, sectionSuffix Section
	if len(sectionPrefixStr) > 0 {
		sectionPrefix = s.parseSectionStr(sectionPrefixStr, nil, nil)
		if sectionPrefix == nil {
			return nil, fmt.Errorf("Could not parse %q as a valid prefix section for section %q", sectionPrefixStr, sectionStr)
		}
	}
	if len(sectionSuffixStr) > 0 {
		sectionSuffix = s.parseSectionStr(sectionSuffixStr, nil, nil)
		if sectionSuffix == nil {
			return nil, fmt.Errorf("Could not parse %q as a valid suffix section for section %q", sectionSuffixStr, sectionStr)
		}
	}
	section := s.parseSectionStr(sectionStr, sectionPrefix, sectionSuffix)
	if section == nil {
		return nil, fmt.Errorf("Could not parse %q as a valid section", sectionStr)
	}
	return section, nil
}

func (s *SectionParser) parseSectionStr(sectionStr string, prefixSection, suffixSection Section) Section {
	for _, sectionType := range s.sectionTypes {
		section := sectionType.generate(prefixSection, sectionStr, suffixSection)
		if section != nil {
			return section
		}
	}
	return nil
}

func (s *SectionParser) SectionHelpTexts() string {
	help := ""
	for _, sectionType := range s.sectionTypes {
		help += sectionType.helpText() + "\n"
	}
	return help
}
