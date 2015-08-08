package ini

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/knq/ini/parser"
)

// Gitconfig style name manipulation allowing for subsections.
//
// Use this as the File.SectionManipFunc.
func GitSectionManipFunc(name string) string {
	n, sub := parser.NameSplitFunc(name)

	if n == "" {
		n = sub
		sub = ""
	}

	// clean up name
	n = strings.TrimSpace(strings.ToLower(n))

	// if there's a subsection in the name
	s := ""
	if sub != "" {
		s = fmt.Sprintf(" \"%s\"", sub)
	}

	return fmt.Sprintf("%s%s", n, s)
}

// Gitconfig style section name formatting.
//
// Effectively inverse of GitSectionManipFunc.
//
// Use this as the File.SectionNameFunc.
func GitSectionNameFunc(name string) string {
	// remove " from string
	n := strings.Replace(strings.TrimSpace(name), "\"", "", -1)

	// replace any spacing with .
	return regexp.MustCompile("[ \t]+").ReplaceAllString(n, ".")
}
