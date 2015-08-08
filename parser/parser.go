/*
	Pigeon-based parser that implements a ini file parser.

	Please see http://godoc.org/github.com/knq/ini for the frontend package.
*/
package parser

//go:generate ./generate.sh

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// Default Key Whitespace
	DefaultLeadingKeyWhitespace = "\t"

	// Default line ending for new lines added to file
	DefaultLineEnding = "\n"

	// Separator token for section.name style keys
	DefaultNameKeySeparator = "."

	// last position
	lastPosition position

	// last text
	lastText string
)

// Section Name Manipulation Function.
//
// This function is used when a section name is created or altered.
//
// Override on a per-File basis by setting File.SectionManipFunc.
func SectionManipFunc(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// Section Name Format Function.
//
// This function is used to format (normalized) the section name for output.
//
// Override on a per-File basis by setting File.SectionNameFunc.
func SectionNameFunc(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// Key Manipulation Function.
//
// Takes a key name and returns the value that used. By default does
// TrimSpace(ToLower(key)).
//
// This function is used when a key is created or altered.
//
// Override on a per-File basis by setting File.KeyManipFunc.
func KeyManipFunc(key string) string {
	return strings.TrimSpace(strings.ToLower(key))
}

// Key Comparison Function.
//
// Passes keys a, b through KeyManipFunc and returns string equality.
//
// This function is used when key names are compared.
//
// Override on a per-File basis by setting File.KeyCompFunc.
func KeyCompFunc(a, b string) bool {
	return KeyManipFunc(a) == KeyManipFunc(b)
}

// Section Name Split Function.
//
// Splits names based on DefaultNameKeySeparator.
//
// Returns section, key name.
//
// This function is used to split keys when being retrieved or set on a File.
//
// Override on a per-File basis by setting File.NameSplitFunc.
func NameSplitFunc(name string) (string, string) {
	idx := strings.LastIndex(name, DefaultNameKeySeparator)

	// no section name
	if idx < 0 {
		return "", name
	}

	return name[:idx], name[idx+1:]
}

// Value Manipulation Function.
//
// This function is used when a key is set.
//
// Override on a per-File basis by setting File.ValueManipFunc.
func ValueManipFunc(value string) string {
	return strings.TrimSpace(value)
}

// Retrieve the last error encountered during parsing.
func LastError() error {
	return errors.New(fmt.Sprintf("error on line %d:%d near '%s'", lastPosition.line, lastPosition.col, lastText))
}

// Helper function taken from pigeon source / examples
func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	return v.([]interface{})
}

// Common interface to Comment, Section, and KeyValuePair.
type Item interface {
	String() string
}

// Line in a File.
type Line struct {
	pos position

	ws   string // Leading whitespace (if any)
	item Item   // A Comment, Section, or KeyValuePair
	le   string
}

// New Line.
func NewLine(pos position, ws string, item Item, le string) *Line {
	return &Line{
		pos: pos,

		ws:   ws,
		item: item,
		le:   le,
	}
}

// Line Stringer.
func (l Line) String() string {
	item := ""
	if l.item != nil {
		item = l.item.String()
	}

	return fmt.Sprintf("%s%s%s", l.ws, item, l.le)
}

// Comment in a File.
type Comment struct {
	pos position

	cs      string // comment separator
	comment string // actual comment
}

// New Comment.
func NewComment(pos position, cs string, comment string) *Comment {
	return &Comment{
		pos: pos,

		cs:      cs,
		comment: comment,
	}
}

// Comment Stringer.
func (c Comment) String() string {
	return fmt.Sprintf("%s%s", c.cs, c.comment)
}

// Key Value Pair in a File.
type KeyValuePair struct {
	//section *Section

	pos position

	key   string
	ws    string
	value string

	comment *Comment
}

// New Key Value Pair.
func NewKeyValuePair(pos position, key, ws, value string, comment *Comment) *KeyValuePair {
	return &KeyValuePair{
		pos: pos,

		key:     key,
		ws:      ws,
		value:   value,
		comment: comment,
	}
}

// Key Value Pair Stringer.
func (kvp KeyValuePair) String() string {
	comment := ""
	if kvp.comment != nil {
		comment = kvp.comment.String()
	}

	return fmt.Sprintf("%s=%s%s%s", kvp.key, kvp.ws, kvp.value, comment)
}
