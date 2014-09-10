package shapefile

import (
	"fmt"
	"errors"
	"io"
	"os"
)

func ReadFile(n string) ([]byte, error) {
	f, err := os.Open(n)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	m := fi.Mode()
	if !m.IsRegular() {
		return nil, errors.New("not a file")
	}
	size := int(fi.Size())
	body := make([]byte, size, size)
	z, err := io.ReadFull(f, body)
	if err != nil {
		return nil, err
	}
	if z != size {
		mesg := fmt.Sprintf("read mismatch: size %d, %d read", size, z)
		return nil, errors.New(mesg)
	}
	return body, nil
}

func lencheck(n int, z int, s string) error {
	if n > z {
		mesg := fmt.Sprintf("need %d for %s, have %d", n, s, z)
		return errors.New(mesg)
	}
	return nil
}
