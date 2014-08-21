package rose

import (
	"bufio"
	"fmt"
	"os"
)

var (
	// cheese string	= "Acapella"
	// cheese string	= "Brie"
	// cheese string	= "Camembert"
	cheese string	= "Danablu"
)

func Shell(q bool, rose *Petal) string {
	if !q {
		run_cmd("echo " + Version, rose)
		run_cmd("interactive on", rose)
		run_cmd("message on", rose)
		run_cmd("echo " + cheese + " shell", rose)
	}
	return run(os.Stdin, rose)
}

func run(f *os.File, rose *Petal) string {
	r := bufio.NewReaderSize(f, 8192)
	for {
		if interactive {
			fmt.Printf("%s", prompt)
		}
		line, err := r.ReadString('\n')
		// fix
		if err != nil || (len(line) >=4 && line[0:4] == "quit") {
			break
		}
		run_cmd(line, rose)
	}
	return "quit"
}
