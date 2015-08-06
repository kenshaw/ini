package parser

import (
	"bytes"
	"os"
)

// File
type File struct {
	lines []*Line // parsed content

	sections []*Section // computed sections
}

// Create a new ini file from passed lines
func NewFile(lines []*Line) *File {
	// create ret object
	ret := &File{
		lines:    lines,
		sections: make([]*Section, 0),
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
			//lastSection.values[kvp.key] = l
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
// Section names are passed through KeyManipFunc
func (f *File) SectionNames() []string {
	names := make([]string, len(f.sections))
	for i, s := range f.sections {
		names[i] = s.Name()
	}
	return names
}

// Add section to file
func (f *File) AddSection(name string) *Section {
	// create key
	k := KeyManipFunc(name)

	// if its "", then avoid retrieving ...
	if k == "" {
		return f.GetSection("")
	}

	// create section
	s := NewSection(position{}, k, "", nil)
	s.file = f

	// add section data to file
	f.sections = append(f.sections, s)

	if f.lines[len(f.lines)-1].item == nil {
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

// Get a section and its starting line number
func (f *File) _getSection(name string) (*Section, int) {
	// blank section isn't actually defined ...
	if name == "" {
		return f.sections[0], 0
	}

	// loop through lines and find section
	for idx, line := range f.lines {
		if s, ok := line.item.(*Section); ok && KeyCompFunc(name, s.name) {
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

// Rename section in file
func (f *File) RenameSection(key, value string) {
	s := f.GetSection(key)
	s.name = KeyManipFunc(value)
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
	s, k := NameSplitFunc(key)

	// get the section
	section := f.GetSection(s)
	if section == nil {
		section = f.AddSection(s)
	}

	section.SetKey(k, value)
}

// Retrieve key from file with name in form of section.key
// If no section is specified, then the empty (first) section is used
// Uses NameSplitFunc to split the key
func (f *File) GetKey(key string) string {
	s, k := NameSplitFunc(key)

	// get the section
	section := f.GetSection(s)
	if section == nil {
		return ""
	}

	return section.Get(k)
}

// Remove key from file with name in form of section.key
// If no section is specified, then the empty (first) section is used
// Uses NameSplitFunc to split the key
func (f *File) RemoveKey(key string) {
	s, k := NameSplitFunc(key)

	// get the section
	section := f.GetSection(s)
	if section == nil {
		return
	}

	section.RemoveKey(k)
}
