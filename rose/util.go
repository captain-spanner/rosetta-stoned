package rose

import (
	"fmt"
)

func fatal(src string, ix int, s string) {
	fmt.Printf("rose: %s\n", s)
	panic(s)
}
