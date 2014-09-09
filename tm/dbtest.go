package main

import (
	"fmt"
	"os"
	"shapefile"
)

const (
	path	string	= "/home/brucee/TM/TM_WORLD_BORDERS-0.3.dbf"
)

func main() {
	_, err := shapefile.MakeDbase(path, os.Stdout)
	if err != nil {
		fmt.Printf("%s: %s\n", path, err.Error())
	}
}
