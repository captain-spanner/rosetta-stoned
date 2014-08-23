package rose

import (
	"bufio"
	"fmt"
	"io"
)

type Petal struct {
	name		string
	rd		io.Reader
	wr		io.Writer
	ewr		io.Writer
	message		bool
	verbose		bool
	xeq		bool
	interactive	bool
	prompt		string
}

type petalreq struct {
	rose	*Petal
	mesg	chan string
}

var (
	petalq	chan *petalreq
)

func MkPetal(name string, rd io.Reader, wr io.Writer, ewr io.Writer, proto *Petal) *Petal {
	p := new(Petal)
	if proto != nil {
		*p = *proto
		if wr != nil {
			p.wr = wr
		}
		if ewr != nil {
			p.ewr = ewr
		}
	} else {
		p.wr = wr
		p.ewr = ewr
	}
	p.name = name
	p.rd = rd
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
	rose := req.rose
	r := bufio.NewReaderSize(rose.rd, 8192)
	for {
		if rose.interactive {
			fmt.Printf("%s", rose.prompt)
		}
		line, err := r.ReadString('\n')
		// fix
		if err != nil || (len(line) >=4 && line[0:4] == "quit") {
			break
		}
		run_cmd(line, rose)
	}
	req.mesg <- ""
}
