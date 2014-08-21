package rose

import (
	"io"
)

type Petal struct {
	name	string
	rd	io.Reader
	wr	io.Writer
	ewr	io.Writer
}

type petalreq struct {
	rose	*Petal
	mesg	chan string
}

var (
	petalq	chan *petalreq
)

func MkPetal(name string, rd io.Reader, wr io.Writer, ewr io.Writer) *Petal {
	p := new(Petal)
	p.name = name
	p.rd = rd
	p.wr = wr
	p.ewr = ewr
	return p
}

func (p *Petal) XeqPetal() string {
	req := new(petalreq)
	req.rose = p
	req.mesg = make(chan string)
	petalq <- req
	return <- req.mesg
}

func petalsrv() {
	for {
		req := <- petalq
		go petalrun(req)
	}
}

func petalrun(req *petalreq) {
}
