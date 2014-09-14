package shapefile

import (
	"fmt"
	"io"
)

const (
	Hdrsize = 100
)

type Header struct {
	version int
	shape   int
	xybox   bbox
	zmbox   bbox
}

func MakeHeader(b []byte, out io.Writer) *Header {
	h := new(Header)
	h.version = int(sl32(b[28:]))
	h.shape = int(sl32(b[32:]))
	makebbox(b[36:], &h.xybox)
	makebbox(b[68:], &h.zmbox)
	if out != nil {
		fmt.Fprintf(out, "version\t%d\n", h.version)
		fmt.Fprintf(out, "shape\t%d\n", h.shape)
	}
	return h
}
