package rose

var (
	stats	[]stat	= make([]stat, 0)
)

type stat interface {
	print()
}

func add_stat(s stat) {
	stats = append(stats, s)
}
