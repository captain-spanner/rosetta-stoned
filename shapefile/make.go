package shapefile

import (
	"io"
)

type Shapefile struct {
	path  string
	shp   *Shapes
	shx   *Index
	dbase *Dbase
	nrecs int
	polys []*polygons
}

func MakeShapefile(n string, out io.Writer) (*Shapefile, error) {
	sf := new(Shapefile)
	sf.path = n
	s, err := MakeShapes(n+".shp", out)
	if err != nil {
		return nil, err
	}
	sf.shp = s
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
	return sf, nil
}
