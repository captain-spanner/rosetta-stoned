package rose

import (
	"fmt"
)

type req struct {
	cmd	func(int, []string, string, int, bool) []string
	min	int
	max	int
	usage	string
}

var (
	none []string = make([]string, 0, 0)
	cmdtab map[string]req = map[string]req {
		"?":	{ cmd_help, 0, -1, "help"},
		"#":	{ cmd_comment, 0, -1, "# comment until end of line"},
		"//":	{ cmd_comment, 0, -1, "// comment until end of line"},
		"help":	{ cmd_help, 0, -1, "help"},
		"echo":	{ cmd_echo, 0, -1, "echo any stuff blah blah"},
		"root":	{ cmd_root, 1, 1, "root <directory>"},
	}
)

func Run_cmd(args []string) []string {
	return run_cmdx(len(args), args, "", 0, false)
}

func Run_cmds(vect [][]string, src string, die bool) (ret [][]string) {
	ret = make([][]string, 0)
	for i, args := range vect {
		r := run_cmdx(len(args), args, src, i + 1, die)
		ret = append(ret, r)
	}
	return ret
}

func run_cmdx(argc int, args []string, src string, ix int, die bool) (ret []string) {
	ret = none
	if argc == 0 {
		return
	}
	if verbose {
		cmd_echo(argc, args, src, 0, die)
	}
	cmd := args[0]
	argc--
	args = args[1:]
	cmdf, found := cmdtab[cmd]
	if !found {
		mesg := cmd + ": unknown command"
		if die {
			fatal(src, ix, mesg)
		} else {
			fmt.Println(mesg)
			return
		}
	}
	return cmdf.cmd(argc, args, src, ix, die)
}

// Commands

func cmd_comment(argc int, args []string, src string, ix int, die bool) []string {
	return nil
}

func cmd_echo(argc int, args []string, src string, ix int, die bool) []string {
	e := "";
	for i, s := range args {
		if i == 0 {
			e = s
		} else {
			e = e + " " + s
		}
	}
	if debug || verbose {
		fmt.Println(e)
	}
	r := make([]string , 1, 1)
	r[0] = e
	return r
}

func cmd_help(argc int, args []string, src string, ix int, die bool) []string {
	fmt.Println("I need somebody!")
	return none
}

func cmd_root(argc int, args []string, src string, ix int, die bool) []string {
	return none
}
