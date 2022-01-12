package sections

import (
	"errors"
	"fmt"
	"github.com/daixiang0/gci/pkg/constants"
	"strings"
)

var SectionParserInst = SectionParser{}

type SectionParser struct {
	sectionTypes []SectionType
}

func (s *SectionParser) RegisterSection(sectionType *SectionType) {
	// TODO test for duplicate aliases
	s.sectionTypes = append(s.sectionTypes, *sectionType)
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

//TODO use error wrapping
func (s *SectionParser) parseSectionComponentsToSection(sectionPrefixStr string, sectionStr string, sectionSuffixStr string) (Section, error) {
	var sectionPrefix, sectionSuffix Section
	var err error
	if len(sectionPrefixStr) > 0 {
		sectionPrefix, err = s.parseSectionStr(sectionPrefixStr, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	if len(sectionSuffixStr) > 0 {
		sectionSuffix, err = s.parseSectionStr(sectionSuffixStr, nil, nil)
		if err != nil {
			return nil, err
		}
	}
	section, err := s.parseSectionStr(sectionStr, sectionPrefix, sectionSuffix)
	if err != nil {
		return nil, err
	}
	return section, nil
}

func (s *SectionParser) parseSectionStr(sectionStr string, prefixSection, suffixSection Section) (Section, error) {
	// create map of all aliases
	aliasMap := map[string]SectionType{}
	for _, sectionType := range s.sectionTypes {
		for _, alias := range sectionType.aliases {
			aliasMap[strings.ToLower(alias)] = sectionType
		}
	}
	// parse everything before the parameter brackets
	sectionComponents := strings.Split(sectionStr, "(")
	sectionType, exists := aliasMap[strings.ToLower(sectionComponents[0])]
	if !exists {
		return nil, fmt.Errorf("Section alias %q not found", sectionComponents[0])
	}
	switch len(sectionComponents) {
	case 1:
		return sectionType.generatorFun("", prefixSection, suffixSection)
	case 2:
		if strings.HasSuffix(sectionComponents[1], ")") {
			return sectionType.generatorFun(strings.TrimSuffix(sectionComponents[1], ")"), prefixSection, suffixSection)
		} else {
			return nil, errors.New("Section parameter is missing closing \")\"")
		}
	}
	return nil, fmt.Errorf("Found more than one %q parameter start sequences: %d", "(", len(sectionComponents))
}

func (s *SectionParser) SectionHelpTexts() string {
	help := ""
	for _, sectionType := range s.sectionTypes {
		aliasesWithParameters := []string{}
		for _, alias := range sectionType.aliases {
			parameterSuffix := ""
			if sectionType.parameterHelp != "" {
				parameterSuffix = "(" + sectionType.parameterHelp + ")"
			}
			aliasesWithParameters = append(aliasesWithParameters, alias+parameterSuffix)
		}
		help += fmt.Sprintf("%s - %s\n", strings.Join(aliasesWithParameters, " | "), sectionType.description)
	}
	return help
}
