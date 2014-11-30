package rose

import (
	"fmt"
)

type req struct {
	cmd   func(int, []string, *Petal) ([]string, int)
	min   int
	max   int
	usage string
	desc  string
}

var (
	none   []string = make([]string, 0, 0)
	cmdtab map[string]req
)

type cmdd interface {
	Src() string
	Index() int
	Die() bool
}

func init_cmds() {
	cmdtab = map[string]req{
		"?":    {cmd_help, 0, -1, "help", "help"},
		"#":    {cmd_comment, 0, -1, "# comment until end of line", "comment"},
		"//":   {cmd_comment, 0, -1, "// comment until end of line", "comment"},
		"base": {cmd_base, 0, 1, "base <name>", "set base corpus"}, "collection": {cmd_collection, 1, 1, "collection [ <name> ]", "manage collections"},
		"corpi":       {cmd_corpi, 0, 0, "corpi", "list corpi"},
		"corpus":      {cmd_corpus, 1, 1, "corpus <name>", "add corpus"},
		"debug":       {cmd_debug, 0, 1, "debug [ <bool> ]", "manage debug"},
		"echo":        {cmd_echo, 0, -1, "echo any stuff blah blah", "echo arguments"},
		"get":         {cmd_get, 2, 2, "get <index> <word>", "get data"},
		"getu":        {cmd_getu, 2, 2, "getu <index> <word>", "get uncached data"},
		"help":        {cmd_help, 0, -1, "help", "help"},
		"index":       {cmd_index, 0, 1, "index [ <name> ]", "manage indexes"},
		"interactive": {cmd_interactive, 0, 1, "interactive [ <bool> ]", "manage interactive"},
		"lookup":      {cmd_lookup, 1, 2, "lookup <word> [ <option> ]", "lookup word is base corpus"},
		"map":         {cmd_map, 0, 1, "map <index>", "set coord map"},
		"message":     {cmd_message, 0, 1, "message [ <bool> ]", "manage message"},
		"part":        {cmd_part, 2, 2, "part <index> <word>", "get part of speach"},
		"pop":         {cmd_pop, 2, 3, "pop <index> <word> [ <depth> ]", "populate part of speach"},
		"regions":     {cmd_regions, 1, 1, "regions <directory>", "load region database"},
		"root":        {cmd_root, 1, 1, "root <directory>", "set root"},
		"run":         {cmd_run, 1, 1, "run <file>", "run commands from a file"},
		"verbose":     {cmd_verbose, 0, 1, "verbose [ <bool> ]", "manage verbose"},
		"word":        {cmd_word, 1, 1, "word <word>", "find word in indexed wordlists"},
		"xeq":         {cmd_xeq, 0, 1, "xeq [ <bool> ]", "manage xeq"},
	}
}

func Run_cmd(args []string, rose *Petal) ([]string, int) {
	mesgs, errs := rose.run(args)
	return mesgs, errs
}

func run_cmd(s string, rose *Petal) ([]string, int) {
	args := smash_cmd(s)
	v, e := Run_cmd(args, rose)
	return v, e
}

func run_cmdx(argc int, args []string, rose *Petal) (ret []string, err int) {
	ret = none
	err = 0
	if argc == 0 {
		return
	}
	if rose.xeq {
		cmd_echo(argc, args, rose)
	}
	cmd := args[0]
	set := false
	if len(cmd) >= 2 {
		if cmd[0] == '#' {
			args[0] = cmd[1:]
			cmd = "#"
			set = true
		} else if len(cmd) >= 3 && cmd[0:2] == "//" {
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
		fmt.Fprintln(rose.wr, mesg)
		err = 1
		return
	}
	if argc < cmdf.min || (cmdf.max >= 0 && argc > cmdf.max) {
		fmt.Fprintln(rose.wr, "usage:", cmdf.usage)
		err = 1
		return
	}
	return cmdf.cmd(argc, args, rose)
}
