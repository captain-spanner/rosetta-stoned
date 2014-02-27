package rose

import (
	"fmt"
)

const (
	Copyright string	= "Copyright Bruce Ellis 2014"
	Version string	= "Pusbox 0.0"
	Config string	= "stone.conf"
)

var (
	root string
)

func init() {
	fmt.Println("Init...")
	configure()
}
