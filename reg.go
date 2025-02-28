package main

const (
	// Lower registers (3-bit)
	R0 uint16 = iota
	R1
	R2
	R3
	R4
	R5
	R6
	R7

	// Upper registers (4-bit)
	R8
	R9
	R10
	R11
	RINS
	RPC
	RSP
	RF
	REGISTER_COUNT
)
