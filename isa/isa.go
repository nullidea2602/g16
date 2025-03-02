package isa

const ( // 3-bit RR memory modes
	DD  uint16 = iota // Direct <- Direct: $rx, $ry
	DLI               // Direct(Lower) <- Indirect: $rx, @ry
	DUI               // Direct(Upper) <- Indirect: %rx, @ry
	DWI               // Direct(Word) <- Indirect&Indirect+1: &rx, @ry
	II                // Indirect <- Indirect: @rx, @ry
	IDL               // Indirect <- Direct(Lower): @rx, $ry
	IDU               // Indirect <- Direct(Upper): @rx, %ry
	IDW               // Indirect&Indirect+1 <- Direct(Word): @rx, &ry
)

const (
	RR  uint16 = iota // _rx, _ry
	RLI               // $rl, #i
	RUI               // $rl, ^i
	RLA               // $rl, =i
)

const ( // 5-bit opcode
	HALT uint16 = iota

	MOV   // RX RY flag
	MOVI  // RL(L) <- IMM/ZP
	MOVIU // RL(U) <- IMM
	MOVIO // RL <- RPC + Address Offset

	INC
	DEC
	ADD
	ADDI
	SUB
	SUBI
	MUL
	DIV

	AND
	OR
	XOR
	NOT
	SHL
	SHR

	JMP
	JE
	JZ
	JNZ
	JC
	JNC
	CALL
	RET

	PUSH
	POP

	NOP
)
