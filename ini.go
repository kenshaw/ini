// Simple package to read/write/manipulate ini files
// Mainly a frontend to github.com/knq/ini/parser
package ini

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/knq/ini/parser"
)

// ini file data
type File struct {
	*parser.File        // ini file data
	Filename     string // filename to read/write from/to
}

// Create a new file
func NewFile() *File {
	lines := make([]*parser.Line, 0)
	inifile := parser.NewFile(lines)

	return &File{
		File:     inifile,
		Filename: "",
	}
}

// Save file (write to filename)
func (f *File) Save() error {
	if f.Filename == "" {
		return errors.New("no filename supplied")
	}

	return f.Write(f.Filename)
}

// gets and sanitizes the file data from source
func sanitizeData(r io.Reader) ([]byte, error) {
	// read all data in
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	// add '\n' to end of data if not present
	if len(data) < 1 || !bytes.Equal(data[len(data)-1:], []byte("\n")) {
		data = append(data, '\n')
	}

	return data, nil
}

// passes the filename/reader to the ini.Parser
func parse(name, filename string, r io.Reader) (*File, error) {
	// sanitize data first (make sure file ends with '\n')
	data, err := sanitizeData(r)
	if err != nil {
		return nil, err
	}

	// pass through ini/parser package
	d, err := parser.Parse(name, data)
	if err != nil {
		return nil, err
	}

	// convert to *parser.File
	inifile, ok := d.(*parser.File)
	if !ok {
		return nil, errors.New(fmt.Sprintf("unable to load %s %s", name, parser.LastError().Error()))
	}

	// create file
	file := &File{
		File:     inifile,
		Filename: filename,
	}

	return file, nil
}

// Load ini file from a reader
func Load(r io.Reader) (*File, error) {
	return parse("<io.Reader>", "", r)
}

// Load ini file from string
func LoadString(inistr string) (*File, error) {
	r := strings.NewReader(inistr)
	return parse("<string>", "", r)
}

// Load ini from a file
// If the file doesn't exist, then an empty File is returned
// It is the caller's job to then write the file to disk
func LoadFile(filename string) (*File, error) {
	// check if the file exists, return a new file if it doesn't
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		file := NewFile()
		file.Filename = filename
		return file, nil
	}

	// if file exists, read and parse it
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return parse(filename, filename, f)
}
