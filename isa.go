package main

const ( // 5-bit opcode
	HALT uint16 = iota

	// ZeroPage ops?

	MOVRR  // RX <- RY
	MOVRM  // RX <- [RY]
	MOVRI  // RL(L) <- IMM
	MOVRIU // RL(U) <- IMM
	MOVMR  // [RX] <- RY(L)
	MOVMRU // [RX] <- RY(U)
	MOVMRW // [RX], [RX+1] <- RY
	MOVMI  // [RL] <- IMM

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
