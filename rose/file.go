package rose

import (
	"bytes"
	"os"
)

const (
	Bsize	 int	= 1024
)

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

func checkdir(s string) string {
	f, err := os.Open(s)
	if err != nil {
		return s + ": open failed"
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return s + ": stat failed"
	}
	if !fi.IsDir() {
		return s + ": not a directory"
	}
	return ""
}

func readpbytes(p string, s string, z int) ([]byte, string) {
	n := p + "/" + s
	f, err := os.Open(n)
	if err != nil {
		return nil, n + ": open failed"
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, n + ": stat failed"
	}
	m := fi.Mode()
	if !m.IsRegular() {
		return nil, n + ": not a file"
	}
	b := make([]byte, z, z)
	sz, err := f.Read(b)
	if err != nil {
		return nil, n + ": read failed"
	}
	return b[:sz], ""
}

func readpstr(p string, s string) (string, string) {
	b, err := readpbytes(p, s, Bsize)
	if err != "" {
		return "", err
	}
	return string(b), ""
}
