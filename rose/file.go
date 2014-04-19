package rose

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

const (
	Bsize	 int	= 1024
)

func readfileb(file *os.File) ([]byte, string) {
	fi, err := file.Stat()
	if err != nil {
		return nil, "stat failed"
	}
	if !fi.Mode().IsRegular() {
		return nil, "not a regular file"
	}
	size := fi.Size()
	data := make([]byte, size, size)
	n, err := io.ReadFull(file, data)
	if err != nil {
		return nil, "read error"
	}
	if int64(n) != size {
		return nil, "short read"
	}
	return data, ""
}

func readlines(file *os.File) ([]string, string) {
	data, err := readfileb(file)
	if err != "" {
		return nil, err
	}
	dvect := bytes.Split(data, newline)
	return bvect_to_svect(dvect), ""
}

func readwordlist(p string) []string {
	w := p + "/" + wordlist
	file, err := fileopen(w, false)
	if err != "" {
fmt.Printf("no %s\n", w)
		return nil
	}
	data, err := readfileb(file)
	if err != "" {
fmt.Printf("error reading %s\n", w)
		return nil
	}
	dvect := bytes.Split(data, newline)
	return bvect_to_svect(dvect)
}

func checkdir(s string) string {
	f, err := fileopen(s, true)
	if f != nil {
		f.Close()
	}
	if err != "" {
		return s + ": " + err
	} else {
		return ""
	}
}

func fileopen(n string, dir bool) (*os.File, string) {
	f, err := os.Open(n)
	if err != nil {
		return nil, "open failed"
	}
	fi, err := f.Stat()
	if err != nil {
		return f, "stat failed"
	}
	m := fi.Mode()
	if dir {
		if !m.IsDir() {
			return f, "not a directory"
		}
	} else {
		if !m.IsRegular() {
			return f, "not a file"
		}
	}
	return f, ""
}

func readpbytes(p string, s string) ([]byte, string) {
	n := p + "/" + s
	f, err := fileopen(n, false)
	if err != "" {
		if f != nil {
			f.Close()
		}
		return nil, n + ": " + err
	}
	defer f.Close()
	b, err := readfileb(f)
	if err != "" {
		return nil, n + ": read failed"
	}
	return b, ""
}

func readpstr(p string, s string) (string, string) {
	b, err := readpbytes(p, s)
	if err != "" {
		return "", err
	}
	l := len(b)
	if l > 0 && b[l-1] == '\n' {
		b = b[:l-1]
	}
	return string(b), ""
}
