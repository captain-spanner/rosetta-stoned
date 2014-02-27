package rose

import (
	"os"
)

func configure() {
	sf, err := os.Open(Config)
	if err != nil {
		fatal("open config failed")
	}
	_ = sf
}
