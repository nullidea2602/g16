package main

const ( // 3-bit RR memory modes
	DD  uint16 = iota // Direct-Direct
	DLI               // Direct(Lower)-Indirect
	DUI               // Direct(Upper)-Indirect
	DWI               // Direct(Word)-Indirect&Indirect+1
)

const ( // 5-bit opcode
	HALT uint16 = iota

	MOV   // RX RY flag
	MOVI  // RL(L) <- IMM, ZP address
	MOVIU // RL(U) <- IMM
	MOVIM // [RL] <- IMM, ASCII

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
