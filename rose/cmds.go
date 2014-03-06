package rose

import (
	"fmt"
)

// Commands

func cmd_comment(argc int, args []string, cmdi cmdd) ([]string, int) {
	return none, 0
}

func cmd_bool(name string, argc int, args []string, vp *bool) ([]string, int) {
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
	if message && argc == 0 {
		fmt.Println(m)
	}
	return strv(m), e
}

func cmd_collection(argc int, args []string, cmdi cmdd) ([]string, int) {
	if argc == 0 {
		// v := print_indexes()
		return none, 0
	} else {
		e, m  := make_collection(args[0])
		if m != "" {
			return strv(m), e
		} else {
			return none, e
		}
	}
}

func cmd_corpus(argc int, args []string, cmdi cmdd) ([]string, int) {
	if argc == 0 {
		// v := print_indexes()
		return none, 0
	} else {
		opt := ""
		if argc > 1 {
			opt = args[1]
		}
		e, m  := make_corpus(args[0], opt)
		if m != "" {
			return strv(m), e
		} else {
			return none, e
		}
	}
}

func cmd_debug(argc int, args []string, cmdi cmdd) ([]string, int) {
	m, e := cmd_bool("debug", argc, args, &debug)
	return m, e
}

func cmd_echo(argc int, args []string, cmdi cmdd) ([]string, int) {
	e := "";
	for i, s := range args {
		if i == 0 {
			e = s
		} else {
			e = e + " " + s
		}
	}
	if message {
		fmt.Println(e)
	}
	return strv(e), 0
}

func cmd_help(argc int, args []string, cmdi cmdd) ([]string, int) {
	m := "I need somebody!"
	if message {
		fmt.Println(m)
	}
	return strv(m), 0
}

func cmd_index(argc int, args []string, cmdi cmdd) ([]string, int) {
	if argc == 0 {
		v := print_indexes()
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

func cmd_message(argc int, args []string, cmdi cmdd) ([]string, int) {
	m, e := cmd_bool("message", argc, args, &message)
	return m, e
}

func cmd_raw(argc int, args []string, cmdi cmdd) ([]string, int) {
	_, m, e := fetch_raw(args[0], args[1])
	return m, e
}

func cmd_root(argc int, args []string, cmdi cmdd) ([]string, int) {
	if root != "" {
		m := "root already set"
		diagx(m, cmdi.Src(), cmdi.Index(), cmdi.Die())
		return strv(m), 1
	}
	root = args[0]
	m := checkdir(root)
	if m != "" {
		fatal(root, 0, m)
	}
	if message {
		fmt.Printf("root = %q\n", root)
	}
	glob_set("root", root)
	return none, 0
}

func cmd_verbose(argc int, args []string, cmdi cmdd) ([]string, int) {
	m, e := cmd_bool("verbose", argc, args, &verbose)
	return m, e
}

func cmd_xeq(argc int, args []string, cmdi cmdd) ([]string, int) {
	m, e := cmd_bool("xeq", argc, args, &xeq)
	return m, e
}
