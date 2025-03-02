package isa

const (
	FSIGN     uint16 = 1 << 15 // Sign flag
	FZERO     uint16 = 1 << 14 // Zero flag
	FCARRY    uint16 = 1 << 13 // Carry flag
	FOVERFLOW uint16 = 1 << 12 // Overflow flag
)
