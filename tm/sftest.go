package main

import (
	"fmt"
	"shapefile"

	"flag"
	"log"
	"os"
	"runtime/pprof"
)

const (
	// path	string	= "/home/brucee/TM/TM_WORLD_BORDERS-0.3"
	path string = "/n/fossil/index/TM/TM_WORLD_BORDERS-0.3"
	test bool   = true
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	s, err := shapefile.MakeShapefile(path, nil)
	if err != nil {
		fmt.Printf("%s: %s\n", path, err.Error())
	}
	if test {
		s.Stoch(10000000, 11271)
	}
}
