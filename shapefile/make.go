package shapefile

import (
	"io"
)

type Shapefile struct {
	path	string
	dbase	*Dbase
}

func MakeShapefile(n string, out io.Writer) (*Shapefile, error) {
	sf := new(Shapefile)
	sf.path = n
	d, err := MakeDbase(n+".dbf", out)
	if err != nil {
		return nil, err
	}
	sf.dbase = d
	return sf, nil
}
