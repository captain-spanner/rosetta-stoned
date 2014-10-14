package shapefile

import (
	"fmt"
	"io"
	"os"
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

func (b *bbox) prints(out io.Writer, s string) {
	fmt.Fprintf(out, "%s: ", s)
	b.print(out)
}

func (b *bbox) full(o *bbox, eps2 float64) bool {
	return b.area()-o.area() <= eps2
}

func (b *bbox) inside(o *bbox) bool {
	return b.xmin >= o.xmin && b.ymin >= o.ymin && b.xmax <= o.xmax && b.ymax <= o.ymax
}

func (b *bbox) area() float64 {
	return (b.xmax - b.xmin) * (b.ymax - b.ymin)
}

func (b *bbox) normal() bool {
	return b.xmax > b.xmin && b.ymax > b.ymin
}

func (b *bbox) divide() []*bbox {
	d := make([]*bbox, 4, 4)
	mx := (b.xmin + b.xmax) / 2.
	my := (b.ymin + b.ymax) / 2.
	if Qdebug2 {
		fmt.Printf("divide: mx = %f, my = %f\n", mx, my)
	}
	for i := 0; i < 4; i++ {
		n := new(bbox)
		switch i {
		case 0:
			n.xmin = mx
			n.xmax = b.xmax
			n.ymin = b.ymin
			n.ymax = my
		case 1:
			n.xmin = mx
			n.xmax = b.xmax
			n.ymin = my
			n.ymax = b.ymax
		case 2:
			n.xmin = b.xmin
			n.xmax = mx
			n.ymin = b.ymin
			n.ymax = my
		case 3:
			n.xmin = b.xmin
			n.xmax = mx
			n.ymin = my
			n.ymax = b.ymax
		}
		d[i] = n
		if Qdebug2 {
			n.print(os.Stdout)
		}
	}
	return d
}

func max(a float64, b float64) float64 {
	if a > b {
		return a
	} else {
		return b
	}
}

func min(a float64, b float64) float64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func (b *bbox) intersection(a *bbox) *bbox {
	xmin := max(a.xmin, b.xmin)
	ymin := max(a.ymin, b.ymin)
	xmax := min(a.xmax, b.xmax)
	ymax := min(a.ymax, b.ymax)
	if xmax < xmin || ymax < ymin {
		return nil
	}
	r := new(bbox)
	r.xmin = xmin
	r.ymin = ymin
	r.xmax = xmax
	r.ymax = ymax
	return r
}

func (b *bbox) enclosed(p *point) bool {
	return p.x >= b.xmin && p.x <= b.xmax && p.y >= b.ymin && p.y <= b.ymax
}

func (b *bbox) encloseds(pt *point) string {
	if b.enclosed(pt) {
		return "yes"
	} else {
		return "no"
	}
}
