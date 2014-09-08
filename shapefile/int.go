package shapefile

func ul16(b []byte) uint16 {
	return (uint16(b[1]) << 8) | uint16(b[0])
}

func ul32(b []byte) uint32 {
	return (uint32(b[3]) << 24) | (uint32(b[2]) << 16) |
		(uint32(b[1]) << 8) | uint32(b[0])
}

func sl16(b []byte) int16 {
	return int16(ul16(b))
}

func sl32(b []byte) int32 {
	return int32(ul32(b))
}

func ub16(b []byte) uint16 {
	return (uint16(b[0]) << 8) | uint16(b[1])
}

func ub32(b []byte) uint32 {
	return (uint32(b[0]) << 24) | (uint32(b[1]) << 16) |
		(uint32(b[2]) << 8) | uint32(b[3])
}

func sb16(b []byte) int16 {
	return int16(ul16(b))
}

func sb32(b []byte) int32 {
	return int32(ul32(b))
}
