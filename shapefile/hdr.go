package shapefile

import (
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

func MakeHeader(b []byte) *Header {
	h := new(Header)
	h.fid = int(sb32(b[0:]))
	h.size = int(sb32(b[24:]))
	h.version = int(sb32(b[28:]))
	h.shape = int(sb32(b[32:]))
	makebbox(b[36:], &h.xybox)
	makebbox(b[68:], &h.zmbox)
	return h
}

func makebbox(b []byte, bb *bbox) {
	bb.xmin = fl64(b[0:])
	bb.ymin = fl64(b[8:])
	bb.xmax = fl64(b[16:])
	bb.ymax = fl64(b[24:])
}
