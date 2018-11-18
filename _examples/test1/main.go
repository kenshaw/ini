// _examples/test1/main.go
package main

import (
	"log"

	"github.com/knq/ini"
)

func main() {
	f, err := ini.LoadFile("blah.ini")
	if err != nil {
		log.Fatal(err)
	}

	sect := f.GetSection("newSection 7")
	if sect == nil {
		sect = f.AddSection("newSection 7")
	}

	sect.SetKey("abracadabra", "fucking magical")
	x := sect.Keys()
	log.Printf(">>> keys: %#v\n", x)

	log.Printf(">>> sect: '%s'\n", sect)

	err = f.Save()
	log.Printf(">>>> file: '%s'\n", f)
	if err != nil {
		log.Fatal(err)
	}
}
