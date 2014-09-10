package shapefile

import (
	"fmt"
)

func dbstr(b []byte) string {
	n := len(b)
	l := 0
	for i := n - 1; i >= 0; i-- {
		if b[i] != ' ' {
			l = i
			break
		}
	}
	return string(b[:l+1])
}

func sbyte(b byte) string {
	if b > 0 && b <= 27 {
		return fmt.Sprintf("^%c", b+'A'-1)
	} else {
		return fmt.Sprintf("0x%02X", b)
	}
}
