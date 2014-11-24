package shapefile

func (s *Shapefile) Where(p *Point) int {
	return s.quad.Search(p)
}
