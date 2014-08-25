package rose

func services() {
	addixq = make(chan *index)
	addcorq = make(chan *corpus)
	pcacheq	= make(chan *pcreq)
	petalq = make(chan *petalreq)
	runq = make(chan *runreq)
	go addixsrv()
	go addcorsrv()
	go pcachesrv()
	go petalsrv()
	go runsrv()
}
