package rose

func services() {
	petalq = make(chan *petalreq)
	addixq = make(chan *index)
	runq = make(chan *runreq)
	go petalsrv()
	go addixsrv()
	go runsrv()
}
