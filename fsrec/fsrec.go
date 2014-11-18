package fsrec

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
)

const maxbsz = 8192

type Fsrec struct {
	name    string
	file    *os.File
	recsz   int
	keysz   int
	filesz  int
	nrecs   int
	blocksz int
}

func (f *Fsrec) Print() {
	fmt.Printf("File: %q\n", f.name)
	fmt.Printf("Record Size: %d\n", f.recsz)
	fmt.Printf("Key Size: %d\n", f.keysz)
	fmt.Printf("File Size: %d\n", f.filesz)
	fmt.Printf("Number of records: %d\n", f.nrecs)
	fmt.Printf("Block Size: %d\n", f.blocksz)
}

func MakeFsrec(n string, rs int, ks int) (*Fsrec, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	m := fi.Mode()
	if !m.IsRegular() {
		return nil, errors.New("not a file")
	}
	fr := new(Fsrec)
	fr.name = n
	fr.file = f
	fr.recsz = rs
	fr.keysz = ks
	fr.filesz = int(fi.Size())
	fr.nrecs = fr.filesz / fr.recsz
	if fr.recsz > maxbsz {
		fr.blocksz = fr.recsz
	} else {
		fr.blocksz = (maxbsz / fr.recsz) * fr.recsz
	}
	return fr, nil
}

func (fs *Fsrec) Search(k []byte) []byte {
	n := sort.Search(fs.nrecs, func(i int) bool { return fs.geq(i, k) })
	if n < 0 {
		return nil
	}
	return fs.getrec(n)
}

func (fs *Fsrec) geq(n int, k []byte) bool {
	return bytes.Compare(fs.getrec(n)[:fs.keysz], k) >= 0
}

func (fs *Fsrec) getrec(n int) []byte {
	z := fs.recsz
	b := make([]byte, z, z)
	n, err := fs.file.ReadAt(b, int64(n*z))
	if err != nil {
		log.Fatal(err)
	}
	return b
}
