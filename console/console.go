package console

import (
	"code/g16/cpu"
	"fmt"
)

const CONSOLE_ADDRESS = 0xFF
const CONSOLE_BUFFER_SIZE = 0x3F

type Console struct {
	buffer [CONSOLE_BUFFER_SIZE]byte
	index  uint8
}

func (c *Console) Step(ram [cpu.RAM_SIZE]byte) {
	if ram[CONSOLE_ADDRESS] == 0 {
		return
	}

	if ram[CONSOLE_ADDRESS] == 0x0A {
		fmt.Printf("Console output: %s\n", c.buffer[:c.index])
		ram[CONSOLE_ADDRESS] = 0
		c.index = 0
	} else {
		c.buffer[c.index] = ram[CONSOLE_ADDRESS]
		ram[CONSOLE_ADDRESS] = 0
		c.index++
	}
}
