package main

import (
	"fmt"
	"fsrec"
)

const file = "/n/fossil/index/geo/noun.coords"

// const file = "/home/brucee/where/noun.coords"

// 52936: 099854173  37.4262 48.0673

func main() {
	fs, _ := fsrec.MakeFsrec(file, 32, 9)
	fs.Print()
	n := fs.Search([]byte("099854173"))
	r := fs.GetRec(n)
	fmt.Printf("%d: %s", n, string(r))
}
