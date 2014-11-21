package shapefile

import (
	"fmt"
)

func (s *Shapefile) Process() {
	var x, y float64
	sz := len(s.regs)
	e := make([]int, sz+1, sz+1)
	n := 0
	for {
		r, err := fmt.Scanf("%f %f\n", &x, &y)
		if r != 2 || err != nil {
			break
		}
		c := s.quad.Search(&point{x: x, y: y})
		if c < 0 {
			e[sz]++
		} else {
			e[c]++
		}
		n++
	}
	for i, v := range e {
		if v == 0 {
			continue
		}
		if i == sz {
			fmt.Printf("Pirate land: %d (%d%%)\n", v, v*100/n)
		} else {
			fmt.Printf("Reg %d: %d\n", i, v)
		}
	}
}
