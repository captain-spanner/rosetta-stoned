package shapefile

import (
	"io"
)

type Index struct {
	path    string
	size    int
	body    []byte
	hdr	*Header
	err	string
}

func MakeIndex(n string, out io.Writer) (*Index, error) {
	s := new(Index)
	s.path = n
	body, err := ReadFile(n)
	if err != nil {
		return nil, err
	}
	s.body = body
	s.size = len(body)
	err = lencheck(Hdrsize, s.size, "header")
	if err != nil {
		return nil, err
	}
	s.hdr = MakeHeader(body, out)
	return s, nil
}

func (s *Index) Error() string {
	return "shx: "+s.err
}
