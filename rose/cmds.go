package rose

import (
	"fmt"
)

// Commands

func cmd_comment(argc int, args []string, src string, ix int, die bool) ([]string, int) {
	return nil, 0
}

func cmd_echo(argc int, args []string, src string, ix int, die bool) ([]string, int) {
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
	return r, 0
}

func cmd_help(argc int, args []string, src string, ix int, die bool) ([]string, int) {
	fmt.Println("I need somebody!")
	return none, 0
}

func cmd_root(argc int, args []string, src string, ix int, die bool) ([]string, int) {
	if root != "" {
		diagx("root already set", src, ix, die)
		return none, 1
	}
	root = args[0]
	if verbose {
		fmt.Printf("root = %q\n", root)
	}
	return none, 0
}
