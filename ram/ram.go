package ram

import (
	"code/g16/pins"
	"encoding/binary"
	"log"
)

const RAM_SIZE = 1 << 16
const ROM_START = 0xF000 // ROM mapping starts here

type RAM struct {
	Pins   *pins.Pins
	memory [RAM_SIZE]byte
}

func (ram *RAM) Init(program []byte) {
	copy(ram.memory[ROM_START:], program)
	log.Printf("Loaded program into RAM at %04X\n", ROM_START)
}

func (ram *RAM) ProcessCycle() {
	if ram.Pins.Valid {
		addr := ram.Pins.Address

		if ram.Pins.RW { // Read
			ram.Pins.Data = binary.LittleEndian.Uint16(ram.memory[addr : addr+2])
		} else { // Write
			binary.LittleEndian.PutUint16(ram.memory[addr:addr+2], ram.Pins.Data)
		}
	}
}
