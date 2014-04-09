package rose

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	newline []byte = []byte { '\n' }
)

func fatal(src string, ix int, mesg string) {
	fmt.Printf("rose: %s\n", mesg)
	panic(mesg)
}

func diagx(mesg string, src string, ix int, die bool) {
	if die {
		fatal(src, ix, mesg)
	} else {
		fmt.Println(mesg)
	}
}

func bomb(src string, mesg string) {
	fatal(src, 0, mesg)
}

func wordlists(svect []string) [][]string {
	z := len(svect)
	lvect := make([][]string, z, z)
	for i, s := range svect {
		lvect[i] = smash_cmd(s)
	}
	return lvect
}

func bvect_to_svect(dvect [][]byte) []string {
	z := len(dvect)
	svect := make([]string, z, z)
	for i, b := range dvect {
		svect[i] = string(b)
	}
	return svect
}

func ws(c int) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// do quoting later
func smash_cmd(s string) []string {
	b := []byte(s)
	v := make([][]byte, 0, 0)
	st := -1
	for i, c := range b {
		if ws(int(c)) {
			if st >= 0 {
				z := b[st:i]
				v = append(v, z)
				st = -1
			}
			continue
		}
		if st < 0 {
			st = i
		}
	}
	if st >= 0 {
		z := b[st:len(b)]
		v = append(v, z)
	}
	return bvect_to_svect(v)
}

func strv(s string) []string {
	r := make([]string , 1, 1)
	r[0] = s
	return r
}

func bool_str(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}

func str_bool(s string) (bool, string) {
	s = strings.ToLower(s)
	b := false
	switch s {
	case "+":
		b = true
	case "-":
		b = false
	case "t":
		b = true
	case "f":
		b = false
	case "y":
		b = true
	case "n":
		b = false
	case "true":
		b = true
	case "false":
		b = false
	case "yes":
		b = true
	case "no":
		b = false
	case "on":
		b = true
	case "off":
		b = false
	default:
		return false, s + ": bad operand"
	}
	return b, ""
}

func str_int(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	} else {
		return i
	}
}

func str_intx(s string) int {
	i, err := strconv.ParseInt(s, 16, 32)
	if err != nil {
		return 0
	} else {
		return int(i)
	}
}

func str_uint(s string) uint32 {
	return uint32(str_int(s))
}

func chars_str(b []byte) string {
	s := ""
	for _, c := range b {
		s += fmt.Sprintf(" %c", c)
	}
	return s
}

func psds_str(b []psd) string {
	s := ""
	for _, c := range b {
		s += fmt.Sprintf(" %s", psds[c])
	}
	return s
}

func dptr_str(d *dptr) string {
	return fmt.Sprintf("%s %08d %s %d", psds[d.tag], d.index, poss[d.pos], d.ptr)
}

func uints_str(v []uint32) string {
	s := ""
	for _, u := range v {
		s += fmt.Sprintf(" %d", u)
	}
	return s
}

func uint_strz(u uint32, z int) string {
	if z == 8 {
		return fmt.Sprintf("%08d", u)
	} else if z == 9 {
		return fmt.Sprintf("%09d", u)
	} else {
		return fmt.Sprintf("%d", u)
	}
}

func uints_strz(v []uint32, z int) string {
	s := ""
	for _, u := range v {
		s += " " + uint_strz(u, z)
	}
	return s
}
