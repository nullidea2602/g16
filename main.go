package main

import (
	"encoding/binary"
	"fmt"
)

func toByte(data []uint16) []byte {
	byteSlice := make([]byte, len(data)*2)
	for i, v := range data {
		binary.LittleEndian.PutUint16(byteSlice[i*2:], v) // Use BigEndian if needed
	}
	return byteSlice
}

func main() {
	fmt.Printf("Hello, world!\n")

	cpu := CPU{}
	cpu.reset()

	program := toByte([]uint16{
		0x0000,
	})

	copy(cpu.ram[0x0200:], program)
}

type CPU struct {
	ram  [1 << 16]byte
	reg  [16]uint16
	halt bool
}

func (cpu *CPU) reset() {
	cpu.reg[RSP] = 0x01FF // Initialize the Stack Pointer
	cpu.reg[RPC] = 0x0200 // Start execution at address 0x0200
	cpu.reg[RF] = 0x0000  // Clear all flags, TODO: Wipe ram and reg
	cpu.halt = false
}
