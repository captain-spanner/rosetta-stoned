package shapefile

import (
	"fmt"
	"io"
)

type Shapefile struct {
	path  string
	box   bbox
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
	if out != nil {
		sf.box = sf.shp.hdr.xybox
		sf.box.print(out)
		sf.polys[0].bounds.print(out)
		sf.polys[0].polys[0].bounds.print(out)
		fmt.Fprintln(out, "cw", sf.polys[0].polys[0].cw)
	}
	return sf, nil
}
