package rose

var (
	variables map[string][]string = make(map[string][]string)
	vars *vartab = &vartab{ name: "User Table" }
)

type vartab struct {
	name string
}

func (*vartab) get(s string) ([]string, bool) {
	v, b := variables[s]
	return v, b
}

func (*vartab) set(s string, v []string) (string, bool) {
	variables[s] = v
	return "", true
}
