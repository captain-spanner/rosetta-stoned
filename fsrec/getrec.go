package fsrec

import (
	"log"
)

func (fs *Fsrec) GetRec(n int) []byte {
	z := fs.recsz
	b := make([]byte, z, z)
	n, err := fs.file.ReadAt(b, int64(n*z))
	if err != nil {
		log.Fatal(err)
	}
	return b
}
