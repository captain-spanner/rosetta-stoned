package shapefile

import (
	"io"
)

type Shapefile struct {
	path	string
	shp	*Shp
	shx	*Index
	dbase	*Dbase
}

func MakeShapefile(n string, out io.Writer) (*Shapefile, error) {
	sf := new(Shapefile)
	sf.path = n
	s, err := MakeShp(n+".shp", out)
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
	return sf, nil
}
