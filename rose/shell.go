package rose

import (
	"bufio"
	"fmt"
	"os"
)

var (
	// cheese string	= "Acapella"
	cheese string	= "Brie"
)

func Shell(q bool) string {
	if !q {
		run_cmd("echo " + Version)
		run_cmd("interactive on")
		run_cmd("message on")
		run_cmd("echo " + cheese + " shell")
	}
	return run(os.Stdin)
}

func run(f *os.File) string {
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
		run_cmd(line)
	}
	return "quit"
}
