package shapefile

func (s *Shapefile) Where(p *Point) int {
	return s.quad.Search(p)
}

func (s *Shapefile) Wherexy(x, y float64) int {
	return s.quad.Search(&Point{x: x, y: y})
}

func (s *Shapefile) Getrec(n int) []byte {
	return s.dbase.Getrec(n)
}
