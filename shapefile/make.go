package shapefile

import (
	"io"
)

type Shapefile struct {
	path	string
	shp	*Shp
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
	d, err := MakeDbase(n+".dbf", out)
	if err != nil {
		return nil, err
	}
	sf.dbase = d
	return sf, nil
}
