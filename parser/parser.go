// ini file parser
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

// Section Name Manipulation Function
func SectionManipFunc(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// Section Name Format Function
func SectionNameFunc(name string) string {
	return strings.TrimSpace(strings.ToLower(name))
}

// Key Manipulation Function
// Takes a key name and returns the value that should used
// Default does TrimSpace(ToLower(key))
func KeyManipFunc(key string) string {
	return strings.TrimSpace(strings.ToLower(key))
}

// Key Comparison Function
// Passes keys a, b through KeyManipFunc and returns string equality
func KeyCompFunc(a, b string) bool {
	return KeyManipFunc(a) == KeyManipFunc(b)
}

// Section Name Split Function
// Splits names based on DefaultNameKeySeparator
// Returns section, key name
func NameSplitFunc(name string) (string, string) {
	idx := strings.LastIndex(name, DefaultNameKeySeparator)

	// no section name
	if idx < 0 {
		return "", name
	}

	return name[:idx], name[idx+1:]
}

// Value Manipulation Function
func ValueManipFunc(value string) string {
	return strings.TrimSpace(value)
}

// Retrieve a formatted last error
func LastError() error {
	return errors.New(fmt.Sprintf("error on line %d:%d near '%s'", lastPosition.line, lastPosition.col, lastText))
}

// helper function taken from pigeon source / examples
func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}

	return v.([]interface{})
}

// Common interface to Comment, Section, and KeyValuePair
type Item interface {
	String() string
}

// Line in a file
type Line struct {
	pos position

	ws   string // Leading whitespace (if any)
	item Item   // A Comment, Section, or KeyValuePair
	le   string
}

// New Line
func NewLine(pos position, ws string, item Item, le string) *Line {
	return &Line{
		pos: pos,

		ws:   ws,
		item: item,
		le:   le,
	}
}

// Stringer
func (l Line) String() string {
	item := ""
	if l.item != nil {
		item = l.item.String()
	}

	return fmt.Sprintf("%s%s%s", l.ws, item, l.le)
}

// Comment
type Comment struct {
	pos position

	cs      string // comment separator
	comment string // actual comment
}

// New Comment
func NewComment(pos position, cs string, comment string) *Comment {
	return &Comment{
		pos: pos,

		cs:      cs,
		comment: comment,
	}
}

// Stringer
func (c Comment) String() string {
	return fmt.Sprintf("%s%s", c.cs, c.comment)
}

// Key Value Pair
type KeyValuePair struct {
	//section *Section

	pos position

	key   string
	ws    string
	value string

	comment *Comment
}

// New Key Value Pair
func NewKeyValuePair(pos position, key, ws, value string, comment *Comment) *KeyValuePair {
	return &KeyValuePair{
		pos: pos,

		key:     key,
		ws:      ws,
		value:   value,
		comment: comment,
	}
}

// Stringer
func (kvp KeyValuePair) String() string {
	comment := ""
	if kvp.comment != nil {
		comment = kvp.comment.String()
	}

	return fmt.Sprintf("%s=%s%s%s", kvp.key, kvp.ws, kvp.value, comment)
}
