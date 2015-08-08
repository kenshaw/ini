package parser

import (
	"bytes"
	"os"
)

// File
type File struct {
	// lines in file
	lines []*Line

	// sections in file
	sections []*Section

	// Manipulation function used on section name for AddSection, RenameSection
	SectionManipFunc func(string) string

	// Function used to normalize section name
	SectionNameFunc func(string) string

	// Comparison function used to find section in File
	// Set this to override default comparison behavior
	SectionCompFunc func(string, string) bool

	// Manipulation function used on key in File
	KeyManipFunc func(string) string

	// Comparison function used to find key in File
	KeyCompFunc func(string, string) bool

	// Manipulation function used when setting value in File
	ValueManipFunc func(string) string

	// Function is used to split a key name (such as section.key)
	NameSplitFunc func(string) (string, string)
}

// Create a new ini file from passed lines
func NewFile(lines []*Line) *File {
	// create ret object
	ret := &File{
		lines:    lines,
		sections: make([]*Section, 0),

		// copy manipulation functions
		SectionManipFunc: SectionManipFunc,
		SectionNameFunc:  SectionNameFunc,
		SectionCompFunc:  nil,
		KeyManipFunc:     KeyManipFunc,
		KeyCompFunc:      KeyCompFunc,
		ValueManipFunc:   ValueManipFunc,
		NameSplitFunc:    NameSplitFunc,
	}

	// create default section
	lastSection := NewSection(position{}, "", "", nil)
	lastSection.file = ret
	ret.sections = append(ret.sections, lastSection)

	// loop over lines and build sections/keys
	for _, l := range lines {
		switch l.item.(type) {
		case *Section:
			// get section
			lastSection, _ = l.item.(*Section)

			// save data
			lastSection.file = ret
			ret.sections = append(ret.sections, lastSection)

		case *KeyValuePair:
			// save kvp
			kvp, _ := l.item.(*KeyValuePair)
			lastSection.keys = append(lastSection.keys, kvp.key)
		}
	}

	return ret
}

// Stringer
func (f File) String() string {
	var buf bytes.Buffer

	for _, l := range f.lines {
		buf.WriteString(l.String())
	}

	return buf.String()
}

// Write to file
func (f *File) Write(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(f.String())

	return err
}

// Get line count
func (f *File) LineCount() int {
	return len(f.lines)
}

// Get raw section names in file
func (f *File) RawSectionNames() []string {
	names := make([]string, len(f.sections))
	for i, s := range f.sections {
		names[i] = s.RawName()
	}
	return names
}

// Get section names in file
// Section names are passed through SectionNameFunc
func (f *File) SectionNames() []string {
	names := make([]string, len(f.sections))
	for i, s := range f.sections {
		names[i] = s.Name()
	}
	return names
}

// Add section with raw name
func (f *File) AddSectionRaw(name string) *Section {
	// if its "", then avoid retrieving ...
	if f.sectionNameComp(name, "") {
		return f.GetSection("")
	}

	// create section
	s := NewSection(position{}, name, "", nil)
	s.file = f

	// add section data to file
	f.sections = append(f.sections, s)

	if len(f.lines) > 0 && f.lines[len(f.lines)-1].item == nil {
		// if it's a blank line on the last line, then put it there
		f.lines[len(f.lines)-1].item = s
	} else {
		// default line ending
		le := DefaultLineEnding
		if len(f.lines) > 0 {
			// take line ending from first line if present
			le = f.lines[0].le
		}

		// create the line and append to end
		l := NewLine(position{}, "", s, le)
		f.lines = append(f.lines, l)
	}

	return s
}

// Add section to file and retrieve
// Section name is passed through file's SectionManipFunc
func (f *File) AddSection(name string) *Section {
	return f.AddSectionRaw(f.SectionManipFunc(name))
}

// compares two section names
// uses f.SectionCompFunc if present, otherwise
func (f *File) sectionNameComp(a, b string) bool {
	if f.SectionCompFunc != nil {
		return f.SectionCompFunc(a, b)
	}

	return f.SectionNameFunc(a) == f.SectionNameFunc(b)
}

// Get a section and its starting line number
func (f *File) _getSection(name string) (*Section, int) {
	n := f.SectionManipFunc(name)

	// blank section isn't actually defined ...
	if f.sectionNameComp(n, "") {
		return f.sections[0], 0
	}

	// loop through lines and find section
	for idx, line := range f.lines {
		if s, ok := line.item.(*Section); ok && f.sectionNameComp(n, s.name) {
			return s, idx
		}
	}

	return nil, -1
}

// Get section from file
func (f *File) GetSection(name string) *Section {
	s, _ := f._getSection(name)
	return s
}

// Rename section in file (RAW)
func (f *File) RenameSectionRaw(name, value string) {
	s := f.GetSection(name)
	s.name = value
}

// Rename section in file
// value will be passed through the file's SectionManipFunc
func (f *File) RenameSection(name, value string) {
	f.RenameSectionRaw(name, f.SectionManipFunc(value))
}

// Remove section and all related lines from file
func (f *File) RemoveSection(name string) {
	section, start := f._getSection(name)
	if section == nil {
		return
	}

	// save copy of line ending
	le := f.lines[0].le

	// find next section
	end := start + 1
	for ; end < len(f.lines); end++ {
		if _, ok := f.lines[end].item.(*Section); ok {
			break
		}
	}

	// remove from f.lines
	f.lines = append(f.lines[:start], f.lines[end:]...)

	// if we removed all lines, then put a blank line back in
	if len(f.lines) < 1 {
		line := NewLine(position{}, "", nil, le)
		f.lines = []*Line{line}
	}

	// find in f.sections
	pos := -1
	for idx, s := range f.sections {
		if section == s {
			pos = idx
			break
		}
	}

	// remove from f.sections
	f.sections = append(f.sections[:pos], f.sections[pos+1:]...)
}

// Set key in file with name in form of section.key
// If no section is specified, then the empty (first) section is used
// Uses NameSplitFunc to split the key
func (f *File) SetKey(key, value string) {
	name, k := f.NameSplitFunc(key)

	// get the section
	section := f.GetSection(name)
	if section == nil {
		section = f.AddSection(name)
	}

	section.SetKey(k, value)
}

// Retrieve key from file with name in form of section.key
// If no section is specified, then the empty (first) section is used
// Uses file's NameSplitFunc to split the key
func (f *File) GetKey(key string) string {
	name, k := f.NameSplitFunc(key)

	// get the section
	section := f.GetSection(name)
	if section == nil {
		return ""
	}

	return section.Get(k)
}

// Remove key from file with name in form of section.key
// If no section is specified, then the empty (first) section is used
// Uses file's NameSplitFunc to split the key
func (f *File) RemoveKey(key string) {
	name, k := f.NameSplitFunc(key)

	// get the section
	section := f.GetSection(name)
	if section == nil {
		return
	}

	section.RemoveKey(k)
}
