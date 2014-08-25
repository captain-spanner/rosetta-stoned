package rose

func services() {
	addixq = make(chan *index)
	addcorq = make(chan *corpus)
	petalq = make(chan *petalreq)
	runq = make(chan *runreq)
	go addixsrv()
	go addcorsrv()
	go petalsrv()
	go runsrv()
}
