package cpu

const STACK_TOP = 0x01FF
const PROGRAM_START = 0xF000
const BYTES_PER_WORD = 2
const BITS_PER_WORD = 16
const OPCODE_OFFSET = 11
const OPCODE_WIDTH = 5
const FLAG_OFFSET = 8
const FLAG_WIDTH = 3
const RX_OFFSET = 4
const RY_OFFSET = 0
const RL_OFFSET = 8
const R_WIDTH = 4
const RL_WIDTH = 3
const I_OFFSET = 0
const I_WIDTH = 8
const HIGHBYTE_OFFSET = 8

type CPUState int

const (
	FetchInstruction CPUState = iota
	DecodeInstruction
	ExecuteInstruction
	Halted
)
