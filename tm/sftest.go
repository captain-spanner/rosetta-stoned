package main

import (
	"fmt"
	"shapefile"
)

const (
	path	string	= "/home/brucee/TM/TM_WORLD_BORDERS-0.3"
	test	bool	= true
)

func main() {
	s, err := shapefile.MakeShapefile(path, nil)
	if err != nil {
		fmt.Printf("%s: %s\n", path, err.Error())
	}
	if test {
		s.Stoch(10000000, 1127)
	}
}
