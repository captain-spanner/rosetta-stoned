package rose

import (
	"fmt"
)

const (
	Copyright string	= "Copyright Bruce Ellis 2014"
	Version string		= "Pusbox 0.0"
	Config string		= "stone.conf"

	debug bool		= true
	verbose bool		= true
)

var (
	root string
)

func init() {
	if debug {
		fmt.Println("Init...")
	}
	configure()
}
