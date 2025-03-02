package cpu

import (
	. "code/g16/isa"
	"fmt"
	"log"
)

const RAM_SIZE = 1 << 16
const STACK_TOP = 0x01FF
const PROGRAM_START = 0x0200
const BYTES_PER_WORD = 2
const BITS_PER_WORD = 16
const OPCODE_OFFSET = 11
const OPCODE_WIDTH = 5
const FLAG_OFFSET = 8
const FLAG_WIDTH = 3
const RX_OFFSET = 4
const RY_OFFSET = 0
const R_WIDTH = 4
const HIGHBYTE_OFFSET = 8

type CPU struct {
	RAM   [RAM_SIZE]byte
	reg   [REGISTER_COUNT]uint16
	op    uint16
	f     uint16
	rx    uint16
	ry    uint16
	mmode uint16
	maddr uint16
	mreg  uint16
	up    uint64
	Halt  bool
}

func (cpu *CPU) Reset() {
	for i := range cpu.reg {
		cpu.reg[i] = 0
	}
	cpu.reg[RSP] = STACK_TOP     // Initialize the Stack Pointer
	cpu.reg[RPC] = PROGRAM_START // Start execution at address 0x0200
	cpu.Halt = false
}

func (cpu *CPU) Load(program []byte) {
	copy(cpu.RAM[PROGRAM_START:], program)
}

func (cpu *CPU) Step() {
	cpu.up++
	log.Printf("Cycle: %d\n", cpu.up)
	cpu.fetch()
	cpu.decode()
	cpu.execute()
	cpu.memoryAccess()
	cpu.writeback()
}

func (cpu *CPU) DumpReg() {
	for i, v := range cpu.reg {
		log.Printf("R%d: %04X\t", i, v)
		if (i+1)%4 == 0 {
			log.Printf("\n")
		}
	}
}

func (cpu *CPU) fetch() {
	log.Printf("RPC: %04X\n[RPC]: %02X\n[RPC+1]: %02X\n",
		cpu.reg[RPC],
		cpu.RAM[cpu.reg[RPC]],
		cpu.RAM[cpu.reg[RPC]+1])
	cpu.reg[RINS] = loadWordLittleEndian(cpu.RAM, cpu.reg[RPC])
	log.Printf("INS: %04X\n", cpu.reg[RINS])
	cpu.reg[RPC] += BYTES_PER_WORD
}

func (cpu *CPU) decode() {
	cpu.op = extract(cpu.reg[RINS], OPCODE_OFFSET, OPCODE_WIDTH)
	log.Printf("OP: %04X\n", cpu.op)

	switch cpu.op {
	case HALT:
		// nothing to do
	case MOV: // $RX $RY
		cpu.f = extract(cpu.reg[RINS], FLAG_OFFSET, FLAG_WIDTH)
		cpu.rx = extract(cpu.reg[RINS], RX_OFFSET, R_WIDTH)
		cpu.ry = extract(cpu.reg[RINS], RY_OFFSET, R_WIDTH)
	default:
		cpu.Halt = true
		fmt.Printf("Halted during decode after %d cycles due to unrecognized OP: %02X\n", cpu.up, cpu.op)
	}
	cpu.reg[RINS] = loadWordLittleEndian(cpu.RAM, cpu.reg[RPC])
}

func (cpu *CPU) execute() {
	switch cpu.op {
	case HALT:
		cpu.Halt = true
		fmt.Printf("Halted after %d cycles.\n", cpu.up)
	case MOV:
		switch cpu.f {
		case DD: // $RX <- $RY
			cpu.reg[cpu.rx] = cpu.reg[cpu.ry]
		case DLI: // $RX(L) <- [$RY]
			cpu.mmode = READ
			cpu.mreg = cpu.rx
			cpu.maddr = cpu.reg[cpu.ry]
		case DUI: // $RX(U) <- [$RY]
			cpu.mmode = READUPPER
			cpu.mreg = cpu.rx
			cpu.maddr = cpu.reg[cpu.ry]
		case DWI: // $RX <- [$RY+1], [$RY]
			cpu.mmode = READWORD
			cpu.mreg = cpu.rx
			cpu.maddr = cpu.reg[cpu.ry]
		default:
			cpu.Halt = true
			fmt.Printf("Halted during execute after %d cycles due to unrecognized FLAG: %02X\n", cpu.up, cpu.f)
		}
	default:
		cpu.Halt = true
		fmt.Printf("Halted during execute after %d cycles due to unrecognized OP: %02X\n", cpu.up, cpu.op)
	}
}

const (
	READ uint16 = iota
	WRITE
	READUPPER
	WRITEUPPER
	READWORD
	WRITEWORD
	COPY
	COPYWORD
)

func (cpu *CPU) memoryAccess() {
	switch cpu.op {
	case MOV:
		switch cpu.mmode {
		case READ:
			cpu.reg[cpu.mreg] = (cpu.reg[cpu.mreg] & 0xFF00) | uint16(cpu.RAM[cpu.maddr])
		case READUPPER:
			cpu.reg[cpu.mreg] = (cpu.reg[cpu.mreg] & 0x00FF) | (uint16(cpu.RAM[cpu.maddr]) << 8)
		case READWORD:
			cpu.reg[cpu.mreg] = loadWordLittleEndian(cpu.RAM, cpu.maddr)
		}
	}

}

func (cpu *CPU) writeback() {

}
