// examples/test2/main.go
package main

import (
	"fmt"
	"log"

	"github.com/knq/ini"
)

var (
	data = `
	firstkey = one

	[some section]
	key = blah ; comment

	[another section]
	key = blah`

	gitconfig = `
	[difftool "gdmp"]
	cmd = ~/gdmp/x "$LOCAL" "$REMOTE"
	`
)

func main() {
	f, err := ini.LoadString(data)
	if err != nil {
		log.Fatal(err)
	}

	s := f.GetSection("some section")

	fmt.Printf("some section.key: %s\n", s.Get("key"))
	s.SetKey("key2", "another value")
	f.Write("out.ini")

	// create a gitconfig parser
	g, err := ini.LoadString(gitconfig)
	if err != nil {
		log.Fatal(err)
	}

	// setup gitconfig name/key manipulation functions
	g.SectionManipFunc = ini.GitSectionManipFunc
	g.SectionNameFunc = ini.GitSectionNameFunc

	fmt.Printf("difftool.gdmp.cmd: %s\n", g.GetKey("difftool.gdmp.cmd"))
}
