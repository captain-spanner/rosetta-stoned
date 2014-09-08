package shapefile

import (
	"encoding/binary"
	"math"
)

func fl64(bytes []byte) float64 {
    bits := binary.LittleEndian.Uint64(bytes)
    float := math.Float64frombits(bits)
    return float
}

func fb64(bytes []byte) float64 {
    bits := binary.BigEndian.Uint64(bytes)
    float := math.Float64frombits(bits)
    return float
}
