package main

import (
	"fmt"
	"log"
)

type CPU struct {
	ram  [1 << 16]byte
	reg  [16]uint16
	op   uint16
	up   uint64
	halt bool
}

func (cpu *CPU) reset() {
	for i := range cpu.reg {
		cpu.reg[i] = 0
	}
	cpu.reg[RSP] = 0x01FF // Initialize the Stack Pointer
	cpu.reg[RPC] = 0x0200 // Start execution at address 0x0200
	cpu.halt = false
}

func (cpu *CPU) init() {
	log.Printf("RPC: %04X\n[RPC]: %02X\n[RPC+1]: %02X\n",
		cpu.reg[RPC],
		cpu.ram[cpu.reg[RPC]],
		cpu.ram[cpu.reg[RPC]+1])
	cpu.reg[RINS] = uint16(cpu.ram[cpu.reg[RPC]]) |
		(uint16(cpu.ram[cpu.reg[RPC]+1]) << 8)
}

func (cpu *CPU) step() {
	cpu.up++
	log.Printf("Cycle: %d\n", cpu.up)
	cpu.reg[RPC] += 2
	log.Printf("INS: %04X\n", cpu.reg[RINS])
	cpu.decode()
}

func (cpu *CPU) decode() {
	cpu.op = cpu.reg[RINS] >> 11
	log.Printf("OP: %04X\n", cpu.op)

	switch cpu.op {
	case HALT:
		cpu.halt = true
		fmt.Printf("Halted after %d cycles.\n", cpu.up)
	case MOVRR: // RX <- RY
		rx := cpu.reg[RINS] >> 4
	default:
		cpu.halt = true
		fmt.Printf("Halted after %d cycles due to unrecognized OP: %02X\n", cpu.up, cpu.op)
	}
	cpu.reg[RINS] = uint16(cpu.ram[RPC])<<8 | uint16(cpu.ram[RPC+1])
}
