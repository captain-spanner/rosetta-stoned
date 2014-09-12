package shapefile

import (
	"fmt"
	"io"
)

const (
	Hdrsize = 100
)

type Header struct {
	fid	int
	size	int
	version	int
	shape	int
	xybox	bbox
	zmbox	bbox
}

type bbox struct {
	xmin	float64
	ymin	float64
	xmax	float64
	ymax	float64
}

func MakeHeader(b []byte, out io.Writer) *Header {
	h := new(Header)
	h.fid = int(sb32(b[0:]))
	h.size = int(sb32(b[24:]))
	h.version = int(sb32(b[28:]))
	h.shape = int(sb32(b[32:]))
	makebbox(b[36:], &h.xybox)
	makebbox(b[68:], &h.zmbox)
	if out != nil {
		fmt.Fprintln(out, "header:")
		fmt.Fprintf(out, "fid\t%d\n", h.fid)
		fmt.Fprintf(out, "size\t%d\n", h.size)
		fmt.Fprintf(out, "version\t%d\n", h.version)
		fmt.Fprintf(out, "shape\t%d\n", h.shape)
	}
	return h
}

func makebbox(b []byte, bb *bbox) {
	bb.xmin = fl64(b[0:])
	bb.ymin = fl64(b[8:])
	bb.xmax = fl64(b[16:])
	bb.ymax = fl64(b[24:])
}
