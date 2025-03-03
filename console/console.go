package console

import (
	"fmt"
)

const CONSOLE_ADDRESS = 0x00
const CONSOLE_BUFFER_SIZE = 0x3F

type Console struct {
	buffer [CONSOLE_BUFFER_SIZE]byte
	index  uint8
}

func (c *Console) Step(ram *[1 << 16]byte) {
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
