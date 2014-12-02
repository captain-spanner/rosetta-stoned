package fsrec

import (
	"log"
)

type getreq struct {
	n    int
	resp chan []byte
}

func (fs *Fsrec) getsrv() {
	for {
		r := <- fs.getq
		r.resp <- fs.getrec(r.n)
	}
}

func (fs *Fsrec) GetRec(n int) []byte {
	req := new(getreq)
	req.n = n
	req.resp = make(chan []byte)
	fs.getq <- req
	b := <- req.resp
	return b
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
