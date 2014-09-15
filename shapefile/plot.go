package shapefile

import (
	"fmt"
	"os"
)

const (
	plfmt = "plot/poly%03d.plot"
)

func (s *Shapefile) mkplotfile(i int, p *polygons) {
	fn := s.dir + fmt.Sprintf(plfmt, i)
	f, err := os.Create(fn)
	if err != nil {
		s.ploterr(err)
		return
	}
	defer f.Close()
	fmt.Fprint(f, "pol {")
	c := p.count
	for i := 0; i < c; i++ {
		p.polys[i].print(f)
	}
	fmt.Fprint(f, " }\n")
}

func (p *polygon) print(f *os.File) {
	fmt.Fprint(f, " {")
	c := p.count
	for i := 0; i < c; i++ {
		if i != 0 {
			fmt.Fprint(f, " ")
		}
		fmt.Fprintf(f, "%f %f", p.points[i].x, p.points[i].y)
	}
	fmt.Fprint(f, "}")
}

func (s *Shapefile) ploterr(err error) {
	s.err = "plot: " + err.Error()
}
