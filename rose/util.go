package rose

import (
	"bytes"
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

func bytes2x(b []byte) []byte {
	v := bytes.Split(b, []byte("  "))
	if len(v) == 1 {
		return v[0]
	} else {
		z := len(b) + 1
		r := make([]byte, z, z)
		l := len(v[0])
		copy(r[:l], v[0])
		copy(r[l:l+3], []byte(" X "))
		copy(r[l+3:], v[1])
		return r
	}
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

func str_uint(s string) uint32 {
	return uint32(str_int(s))
}
