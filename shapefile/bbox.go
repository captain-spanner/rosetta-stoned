package shapefile

import (
	"fmt"
	"io"
)

type bbox struct {
	xmin float64
	ymin float64
	xmax float64
	ymax float64
}

func makebbox(b []byte, bb *bbox) {
	bb.xmin = fl64(b[0:])
	bb.ymin = fl64(b[8:])
	bb.xmax = fl64(b[16:])
	bb.ymax = fl64(b[24:])
}

func (b *bbox) print(out io.Writer) {
	fmt.Fprintf(out, "[(%f, %f) (%f %f)]\n", b.xmin, b.ymin, b.xmax, b.ymax)
}