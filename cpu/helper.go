package cpu

func extract(word uint16, offset uint16, width uint16) uint16 {
	if (offset + width) == BITS_PER_WORD {
		return word >> offset
	}
	var mask uint16 = ((1 << width) - 1)
	return (word >> offset) & mask
}
