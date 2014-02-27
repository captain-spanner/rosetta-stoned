package rose

type req struct {
	min	int
	max	int
	usage	string
}

var (
	cmdtab map[string]func([]string, string, bool) []string
	reqtab map[string]req
)

func init() {
	reqtab = map[string]req {
		"root": { 1, 1, "root directory"},
	}
	add_cmd("root", cmd_root)
}

func add_cmd(cmd string, cmdf func([]string, string, bool) []string) {
	cmdtab[cmd] = cmdf
}

func cmd_root(args []string, src string, die bool) []string {
	return nil
}
