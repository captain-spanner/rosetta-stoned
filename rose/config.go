package rose

import (
	"os"
)

func configure() {
	sf, err := os.Open(Config)
	if err != nil {
		fatal(Config, 0, "open config failed")
	}
	defer sf.Close()
	vect := readlines(sf)
	vargs := wordlists(vect)
	Run_cmds(vargs, Config, true)
}
