package rose

import (
	"fmt"
)

const (
	Copyright string	= "Copyright Bruce Ellis 2014"
	Version string		= "Pusbox 0.0"
	Config string		= "stone.conf"

	hError		= hashc(iota)
	hHashed
	hIndexed
	hLiteral
)

var (
	root	string

	debug bool		= true
	message bool		= true
	verbose bool		= true

	hashes map[hashc]string = map[hashc]string {
		hError:		"error",
		hHashed:	"hashed",
		hIndexed:	"indexed",
		hLiteral:	"literal",
	}
)

type hashc int

func init() {
	if debug {
		fmt.Println("Init...")
	}
	hashes = map[hashc]string {
		hError:		"error",
		hHashed:	"hashed",
		hIndexed:	"indexed",
		hLiteral:	"literal",
	}
	init_syms()
	configure()
}
