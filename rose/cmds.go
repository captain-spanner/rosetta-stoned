package rose

import (
	"fmt"
)

// Commands

func cmd_comment(argc int, args []string, rose *Petal) ([]string, int) {
	return none, 0
}

func cmd_base(argc int, args []string, rose *Petal) ([]string, int) {
	return none, 0
}

func cmd_bool(name string, argc int, args []string, vp *bool, rose *Petal) ([]string, int) {
	e := 0
	m := ""
	if argc == 0 {
		v, f := var_get(name)
		if !f {
			e = 1
			m = "unknown"
		} else {
			m = v
		}
	} else {
		b, v := str_bool(args[0])
		if v != "" {
			e = 1
			m = v
		} else {
			m = bool_str(b)
			glob_set(name, m)
			*vp = b
		}
	}
	m = name + " " + m
	if rose.message && argc == 0 {
		fmt.Println(m)
	}
	return strv(m), e
}

func cmd_collection(argc int, args []string, rose *Petal) ([]string, int) {
	e, m  := make_collection(args[0])
	if m != "" {
		return strv(m), e
	} else {
		return none, e
	}
}

func cmd_corpi(argc int, args []string, rose *Petal) ([]string, int) {
	return none, 0
}

func cmd_corpus(argc int, args []string, rose *Petal) ([]string, int) {
	e, m  := make_corpus(args[0])
	if m != "" {
		return strv(m), e
	} else {
		return none, e
	}
}

func cmd_debug(argc int, args []string, rose *Petal) ([]string, int) {
	m, e := cmd_bool("debug", argc, args, &debug, rose)
	return m, e
}

func cmd_echo(argc int, args []string, rose *Petal) ([]string, int) {
	e := "";
	for i, s := range args {
		if i == 0 {
			e = s
		} else {
			e = e + " " + s
		}
	}
	if rose.message {
		fmt.Fprintln(rose.wr, e)
	}
	return strv(e), 0
}

func cmd_get(argc int, args []string, rose *Petal) ([]string, int) {
	v, m, e := fetch_get(args[0], args[1])
	if rose.message {
		if e != 0 || len(v) == 0 {
			fmt.Printf("%s: %s: not found\n", args[1], args[0])
		} else {
			fmt.Printf("%s", string(v))
		}
	}
	return m, e
}

func cmd_getu(argc int, args []string, rose *Petal) ([]string, int) {
	_, m, e := fetch_raw(args[0], args[1])
	return m, e
}

func cmd_help(argc int, args []string, rose *Petal) ([]string, int) {
	m := "I need somebody!"
	if rose.message {
		fmt.Println(m)
	}
	return strv(m), 0
}

func cmd_index(argc int, args []string, rose *Petal) ([]string, int) {
	if argc == 0 {
		v := rose.print_indexes()
		return v, 0
	} else {
		e, m  := make_index(args[0])
		if m != "" {
			return strv(m), e
		} else {
			return none, e
		}
	}
}

func cmd_interactive(argc int, args []string, rose *Petal) ([]string, int) {
	m, e := cmd_bool("interactive", argc, args, &rose.interactive, rose)
	return m, e
}

func cmd_lookup(argc int, args []string, rose *Petal) ([]string, int) {
	return none, 0
}

func cmd_message(argc int, args []string, rose *Petal) ([]string, int) {
	m, e := cmd_bool("rose.message", argc, args, &rose.message, rose)
	return m, e
}

func fetch_part(args []string, rose *Petal) (part, []string, int) {
	p, m, e := part_get(args[0], args[1])
	if p != nil {
		err := p.Error()
		if err != "" {
			if rose.message {
				fmt.Printf("%s: %s\n", args[1], err)
			}
			return nil, strv(err), 1
		}
	}
	return p, strv(m), e
}

func cmd_part(argc int, args []string, rose *Petal) ([]string, int) {
	p, m, e := fetch_part(args, rose)
	if e != 0 || p == nil {
		if rose.message {
			fmt.Printf("%s: %s: not found\n", args[1], args[0])
		}
	} else {
		if rose.message {
			p.Print()
		}
	}
	return m, e
}

func cmd_pop(argc int, args []string, rose *Petal) ([]string, int) {
	r := 1
	if argc == 3 {
		t := str_int(args[2])
		if t > 0 {
			r = t
		}
	}
	p, m, e := fetch_part(args, rose)
	if e != 0 || p == nil {
		if rose.message {
			fmt.Printf("%s: %s: not found\n", args[1], args[0])
		}
		return m, e
	}
	v, e := p.Populate(r)
	return v, e
}

func cmd_root(argc int, args []string, rose *Petal) ([]string, int) {
	if root != "" {
		m := "root already set"
		return strv(m), 1
	}
	root = args[0]
	m := checkdir(root)
	if m != "" {
		fatal(root, 0, m)
	}
	if rose.message {
		fmt.Printf("root = %q\n", root)
	}
	glob_set("root", root)
	return none, 0
}

func cmd_run(argc int, args []string, rose *Petal) ([]string, int) {
	n := args[0]
	f, m := fileopen(n, false)
	if m != "" {
		d := fmt.Sprintf("%s: %s", n, m)
		if rose.message {
			fmt.Println(d)
		}
		return strv(d), 1
	}
	p := MkPetal("File: " + n, f, nil, nil, rose)
	p.interactive = false
	p.XeqPetal()
	return none, 0
}

func cmd_verbose(argc int, args []string, rose *Petal) ([]string, int) {
	m, e := cmd_bool("verbose", argc, args, &rose.verbose, rose)
	return m, e
}

func cmd_word(argc int, args []string, rose *Petal) ([]string, int) {
	rose.list_ixword(args[0])
	return nil, 0
}

func cmd_xeq(argc int, args []string, rose *Petal) ([]string, int) {
	m, e := cmd_bool("xeq", argc, args, &rose.xeq, rose)
	return m, e
}
