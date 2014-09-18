package shapefile

type quad struct {
	box  bbox
	qbox []*bbox
	down []*quad
	only *subreg
	full []*region
}

type subreg struct {
	box bbox
	reg *region
}

func MakeQuad(b *bbox) *quad {
	q := new(quad)
	q.box = *b
	return q
}
