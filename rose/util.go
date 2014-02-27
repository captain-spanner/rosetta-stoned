package rose

import (
	"fmt"
)

func fatal(s string) {
	fmt.Printf("rose: %s\n", s)
	panic(s)
}
