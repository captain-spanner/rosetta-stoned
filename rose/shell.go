package rose

import (
	"os"
)

var (
	// cheese string	= "Acapella"
	// cheese string	= "Brie"
	// cheese string	= "Camembert"
	cheese string	= "Danablu"
)

func Shell(q bool) string {
	rose := MkPetal("Shell", os.Stdin, nil, nil, Confp)
	if !q {
		run_cmd("echo " + Version, rose)
		run_cmd("interactive on", rose)
		run_cmd("message on", rose)
		run_cmd("echo " + cheese + " shell", rose)
	}
	ret := rose.XeqPetal()
	if server {
		<- fine
	}
	return ret
}
