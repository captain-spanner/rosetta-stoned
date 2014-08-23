package rose

import (
	"os"
)

var (
	Confp	*Petal
)

func configure() {
	sf, err := os.Open(Config)
	if err != nil {
		fatal(Config, 0, "open config failed")
	}
	defer sf.Close()
	Confp = MkPetal("Config", sf, os.Stdout, os.Stderr, nil)
	Confp.prompt = ">> "
	mesg := Confp.XeqPetal()
	if mesg != "" {
		fatal(Config, 0, mesg)
	}
}
