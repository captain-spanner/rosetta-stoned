package rose

import (
	"shapefile"
)

func loadsf() string {
	if root == "" {
		return "root not set"
	}
	path := root + "/" + regions
	s, err := shapefile.MakeShapefile(path, nil)
	if err != nil {
		return err.Error()
	}
	regfile = s
	return ""
}
