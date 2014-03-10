package rose

import (
	"fmt"
)

type symtab interface {
	get(string) ([]string, bool)
	set(string, []string) (string, bool)
}

var (
	symtabs []symtab
)

func init_syms() {
	symtabs = make([]symtab, 2, 2)
	symtabs[0] = globs
	symtabs[1] = vars
	glob_set("Copyright", Copyright)
	glob_set("Version", Version)
	glob_set("Config", Config)
	glob_flag("debug", debug)
	glob_flag("message", message)
	glob_flag("verbose", verbose)
	glob_flag("verbose", interactive)
	glob_set("prompt", prompt)
}

func var_getx(s string) ([]string, bool) {
	for _, t := range symtabs {
		v, b := t.get(s)
		if b {
			return v, true
		}
	}
	return nil, false
}

func var_setx(s string, v []string) (string, bool) {
	for _, t := range symtabs {
		e, b := t.set(s, v)
		if !b && e != "" {
			return e, false
		}
		return "", true
	}
	return "", false
}

func var_get(s string) (string, bool) {
	v, b := var_getx(s)
	if !b {
		return "", false
	}
	if len(v) > 1 && message {
		fmt.Printf("%s: not a scalar\n", s)
	}
	if len(v) == 0 {
		if message {
			fmt.Printf("%s: empty vector\n", s)
		}
		return "empty vector", false
	}
	return v[0], true
}

func var_set(s string, v string) (string, bool) {
	a := make([]string, 1, 1)
	a[0] = v
	e, b := var_setx(s, a)
	return e, b
}


