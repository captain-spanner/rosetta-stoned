package rose

import (
	"bytes"
	"fmt"
	"strings"
	"os"
)

const (
	whitespace string	= " \t\r\n"
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

func readfileb(file *os.File) []byte {
	fi, err := file.Stat()
	name := file.Name()
	if err != nil {
		bomb(name, "stat failed")
	}
	if !fi.Mode().IsRegular() {
		bomb(name, "not a regular file")
	}
	size := fi.Size()
	data := make([]byte, size, size)
	n, err := file.Read(data)
	if err != nil {
		panic(err)
	}
	if int64(n) != size {
		bomb(name, "short read")
	}
	return data
}

func readlines(file *os.File) []string {
	data := readfileb(file)
	dvect := bytes.Split(data, newline)
	return bvect_to_svect(dvect)
}

func wordlists(svect []string) [][]string {
	z := len(svect)
	lvect := make([][]string, z, z)
	for i, s := range svect {
		lvect[i] = strings.Split(s, whitespace)
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
