package rose

import (
	"fmt"
)

type req struct {
	cmd	func(int, []string, string, int, bool) ([]string, int)
	min	int
	max	int
	usage	string
}

var (
	none []string = make([]string, 0, 0)
	cmdtab map[string]req = map[string]req {
		"?":		{ cmd_help, 0, -1, "help"},
		"#":		{ cmd_comment, 0, -1, "# comment until end of line"},
		"//":		{ cmd_comment, 0, -1, "// comment until end of line"},
		"help":		{ cmd_help, 0, -1, "help"},
		"index":	{ cmd_index, 0, 1, "index [ <name> ]"},
		"echo":		{ cmd_echo, 0, -1, "echo any stuff blah blah"},
		"root":		{ cmd_root, 1, 1, "root <directory>"},
	}
)

func Run_cmd(args []string) ([]string, int) {
	return run_cmdx(len(args), args, "", 0, false)
}

func Run_cmds(vect [][]string, src string, die bool) (ret [][]string, errc int, errv []int) {
	ret = make([][]string, 0)
	errc = 0
	errv = make([]int, 0)
	for i, args := range vect {
		r, e := run_cmdx(len(args), args, src, i + 1, die)
		ret = append(ret, r)
		errc += e
		errv = append(errv, e)
	}
	return
}

func run_cmdx(argc int, args []string, src string, ix int, die bool) (ret []string, err int) {
	ret = none
	err = 0
	if argc == 0 {
		return
	}
	if verbose {
		cmd_echo(argc, args, src, 0, die)
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
		if die {
			fatal(src, ix, mesg)
		} else {
			fmt.Println(mesg)
			err = 1
			return
		}
	}
	return cmdf.cmd(argc, args, src, ix, die)
}
