package rose

var (
	globals	map[string][]string = make(map[string][]string)
	globs *globtab = &globtab{ name: "Sytem Table" }
)

type globtab struct {
	name string
}

func (*globtab) get(s string) ([]string, bool) {
	v, b := globals[s]
	return v, b
}

func (*globtab) set(s string, v []string) (string, bool) {
	return s + ": readonly", false
}

func glob_set(n string, v string) {
	globals[n] = strv(v)
}

func glob_flag(n string, f bool) {
	b := "false"
	if f {
		b = "true"
	}
	glob_set(n, b)
}
