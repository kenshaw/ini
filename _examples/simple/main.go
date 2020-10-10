// _examples/simple/main.go
package main

import (
	"fmt"
	"log"

	"github.com/kenshaw/ini"
)

var (
	simple = `  [mysect1 ] ; second

[ section2 ] # third
# acomment

# yet another comment
`

	data = `    ; zero comment
	first= second # first comment

blah          =;blah comment

[section] # yes
sectionkey=yeah

; section with no comment and no keys
[asection]

[  sect123   ]
awesome = yes ; awesome comment
secondAwesome = blah # yet another




[sectb] ; another

blah = yes #sectb.blahcomment

[sectc] #sectc comment
           key1 awesome    =             something ;blah
key2  ja www  =          another #yes
key3             =                     a         value         
  `

	data2 = `; comment
	baadkey
`

	data3 = `
	# some comment

[section] ; comment
key = value # another comment

		yet_another_key = something
`
)

func main() {
	inifile, err := ini.LoadString(data)
	if err != nil {
		log.Fatal(err)
	}

	val := inifile.GetKey("sect123.awesome")
	fmt.Printf(">> val: %s\n", val)
}
