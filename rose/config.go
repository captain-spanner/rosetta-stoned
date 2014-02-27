package rose

import (
	"os"
)

func configure() {
	sf, err := os.Open(Config)
	if err != nil {
		fatal(Config, 0, "open config failed")
	}
	_ = sf
}
