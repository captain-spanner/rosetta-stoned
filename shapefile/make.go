package shapefile

import (
	"io"
	"strings"
)

const (
	mkplot = false
)

type Shapefile struct {
	path  string
	dir   string
	box   bbox
	holes int
	shp   *Shapes
	shx   *Index
	dbase *Dbase
	nrecs int
	polys []*polygons
	regs  []*region
	err   string
}

func MakeShapefile(n string, out io.Writer) (*Shapefile, error) {
	sf := new(Shapefile)
	sf.path = n
	sf.dir = n[:strings.LastIndex(n, "/")+1]
	s, err := MakeShapes(n+".shp", out)
	if err != nil {
		return nil, err
	}
	sf.shp = s
	sf.box = sf.shp.hdr.xybox
	x, err := MakeIndex(n+".shx", out)
	if err != nil {
		return nil, err
	}
	sf.shx = x
	d, err := MakeDbase(n+".dbf", out)
	if err != nil {
		return nil, err
	}
	sf.dbase = d
	sf.decode(out)
	if out != nil {
		sf.box.print(out)
		sf.polys[0].bounds.print(out)
		sf.polys[0].polys[0].bounds.print(out)
	}
	sf.regs = make([]*region, 0)
	return sf, sf.analyze()
}

func (s *Shapefile) Error() string {
	return "shapefile: " + s.err
}
