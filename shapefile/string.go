package shapefile

import (
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
