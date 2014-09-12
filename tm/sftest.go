package main

import (
	"fmt"
	"os"
	"shapefile"
)

const (
	path	string	= "/home/brucee/TM/TM_WORLD_BORDERS-0.3"
)

func main() {
	_, err := shapefile.MakeShapefile(path, os.Stdout)
	if err != nil {
		fmt.Printf("%s: %s\n", path, err.Error())
	}
}
