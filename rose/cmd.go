package rose

import (
	"fmt"
)

type req struct {
	cmd	func(int, []string, cmdd) ([]string, int)
	min	int
	max	int
	usage	string
	desc	string
}

var (
	none []string = make([]string, 0, 0)
	cmdf *cmdb = &cmdb{ die: false }
	cmdt *cmdb = &cmdb{ die: true }
	cmdtab map[string]req
)

type cmdd interface {
	Src() string
	Index() int
	Die() bool
}

type cmdv struct {
	src	string
	index	int
	die	bool
}

type cmdb struct {
	die	bool
}

func (c *cmdv) Src() string {
	return c.src
}

func (c *cmdv) Index() int {
	return c.index
}

func (c *cmdv) Die() bool {
	return c.die
}

func (c *cmdb) Src() string {
	return ""
}

func (c *cmdb) Index() int {
	return 0
}

func init_cmds() {
	cmdtab = map[string]req {
		"?":		{ cmd_help, 0, -1, "help", "help" },
		"#":		{ cmd_comment, 0, -1, "# comment until end of line", "comment" },
		"//":		{ cmd_comment, 0, -1, "// comment until end of line", "comment" },
		"collection":	{ cmd_collection, 0, 1, "collection [ <name> ]", "manage collections" },
		"corpus":	{ cmd_corpus, 0, 2, "corpus [ <name> [ <option> ] ]", "manage corpi" },
		"debug":	{ cmd_debug, 0, 1, "debug [ <bool> ]", "manage debug" },
		"echo":		{ cmd_echo, 0, -1, "echo any stuff blah blah", "echo arguments" },
		"get":		{ cmd_get, 2, 2, "get <index> <word>", "get data" },
		"getu":		{ cmd_getu, 2, 2, "getu <index> <word>", "get uncached data" },
		"help":		{ cmd_help, 0, -1, "help", "help" },
		"index":	{ cmd_index, 0, 1, "index [ <name> ]", "manage indexes" },
		"interactive":	{ cmd_interactive, 0, 1, "interactive [ <bool> ]", "manage interactive" },
		"message":	{ cmd_message, 0, 1, "message [ <bool> ]", "manage message" },
		"part":		{ cmd_part, 2, 2, "part <index> <word>", "get part of speach" },
		"pop":		{ cmd_pop, 2, 2, "pop <index> <word>", "populate part of speach" },
		"run":		{ cmd_run, 1, 1, "run <file>", "run commands from a file" },
		"root":		{ cmd_root, 1, 1, "root <directory>", "set root" },
		"verbose":	{ cmd_verbose, 0, 1, "verbose [ <bool> ]", "manage verbose" },
		"xeq":		{ cmd_xeq, 0, 1, "xeq [ <bool> ]", "manage xeq" },
	}
}

func (c *cmdb) Die() bool {
	return c.die
}

func Run_cmd(args []string) ([]string, int) {
	return run_cmdx(len(args), args, cmdf)
}

func Run_cmds(vect [][]string, src string, die bool) (ret [][]string, errc int, errv []int) {
	ret = make([][]string, 0)
	errc = 0
	errv = make([]int, 0)
	for i, args := range vect {
		cmd := &cmdv{ src: src, index: i + 1, die: die }
		r, e := run_cmdx(len(args), args, cmd)
		ret = append(ret, r)
		errc += e
		errv = append(errv, e)
	}
	return
}

func run_cmd(s string) ([]string, int) {
	args := smash_cmd(s)
	v, e := Run_cmd(args)
	return v, e
}

func run_cmdx(argc int, args []string, cmdi cmdd) (ret []string, err int) {
	ret = none
	err = 0
	if argc == 0 {
		return
	}
	if xeq {
		cmd_echo(argc, args, cmdi)
	}
	cmd := args[0]
	set := false
	if len(cmd) >= 2 {
		if cmd[0] == '#' {
			args[0] = cmd[1:]
			cmd = "#"
			set = true
		} else if len(cmd) >=3 && cmd[0:2] == "//" {
			args[0] = cmd[2:]
			cmd = "//"
			set = true
		}
	}
	if !set {
		argc--
		args = args[1:]
	}
	cmdf, found := cmdtab[cmd]
	if !found {
		mesg := cmd + ": unknown command"
		if cmdi.Die() {
			fatal(cmdi.Src(), cmdi.Index(), mesg)
		} else {
			fmt.Println(mesg)
			err = 1
			return
		}
	}
	return cmdf.cmd(argc, args, cmdi)
}
