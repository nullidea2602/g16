package main

import (
	"code/g16/assembler"
	"code/g16/console"
	"code/g16/cpu"
	. "code/g16/isa"
	"code/g16/pins"
	"encoding/binary"
	"log"
	"os"
	"time"
)

const RAM_SIZE = 1 << 16
const ROM_START = 0xF000 // ROM mapping starts here

type RAM struct {
	Pins   *pins.Pins
	memory [RAM_SIZE]byte
}

func (ram *RAM) init(program []byte) {
	log.Printf("Program: %04X\n", program)
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

type Bus struct {
	CPU_Pins *pins.Pins
	RAM_Pins *pins.Pins
}

func (bus *Bus) PropagateCycle() {
	if bus.CPU_Pins.Valid {
		bus.RAM_Pins.Address = bus.CPU_Pins.Address
		bus.RAM_Pins.Data = bus.CPU_Pins.Data
		bus.RAM_Pins.RW = bus.CPU_Pins.RW
		bus.RAM_Pins.Valid = true
	} else {
		bus.RAM_Pins.Valid = false
	}
}

func (bus *Bus) ReturnCycle() {
	// If RAM was in read mode and valid, return data to CPU
	if bus.RAM_Pins.Valid && bus.RAM_Pins.RW {
		bus.CPU_Pins.Data = bus.RAM_Pins.Data
		bus.CPU_Pins.Valid = true
	} else {
		bus.CPU_Pins.Valid = false
	}
}

func main() {

	file, err := os.Create("debug.log")
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)

	source := `
	; below is interpreted as movio, meaning mov immediate offset address
	; intermediate is movio, $rl, 0x00
	; cpu executes as mov $rl, ProgramCounter+0x00
	; this keeps opcode, register, and address within 16 bits
	; mov $rx, =label

	mov $r1, =data ; set r1 to address of first character
	mov $r2, #d13 ; number of characters
	mov $r3, =loop ; set r3 to address of loop
	loop: ; address 6
	mov @r0, @r1 ; copy character to stdout @0x0000
	inc $r1 ; advance r1 to address of next character
	inc $r1 ; +2 for next word
	dec $r2 ; count down
	jnz $r3, @r2 ; goto loop
	halt
	data: ; address 18
	#'Hello World!'
	`
	log.Println(source)

	tokens := assembler.Tokenize(source)

	for _, t := range tokens {
		log.Printf("%s: %s\n", t.Type, t.Value)
	}

	assembly := []uint16{
		MOVIO, cpu.R1, 18, // =data
		MOVI, cpu.R2, 13, // number of characters
		MOVIO, cpu.R3, 2, // =loop,
		//loop:
		MOV, II, cpu.R0, cpu.R1,
		INC, cpu.R1,
		INC, cpu.R1,
		DEC, cpu.R2,
		JNZ, cpu.R3, cpu.R2,
		HALT,
		//data:
		//#aHello_World!
		uint16('H'), 0, uint16('e'), 0, uint16('l'), 0, uint16('l'), 0, uint16('o'), 0, uint16(' '), 0,
		uint16('W'), 0, uint16('o'), 0, uint16('r'), 0, uint16('l'), 0, uint16('d'), 0, uint16('!'), 0,
		uint16('\n'),
	}

	log.Println(assembly)

	program := []byte{ // Little-endian
		18, byte(MOVIO<<3) | byte(cpu.R1),
		13, byte(MOVI<<3) | byte(cpu.R2),
		2, byte(MOVIO<<3) | byte(cpu.R3),
		byte(cpu.R0<<4) | byte(cpu.R1), byte(MOV<<3) | byte(II),
		byte(cpu.R1 << 4), byte(INC << 3),
		byte(cpu.R1 << 4), byte(INC << 3),
		byte(cpu.R2 << 4), byte(DEC << 3),
		byte(cpu.R3<<4) | byte(cpu.R2), byte(JNZ << 3),
		0, byte(HALT << 3),
		byte('H'), 0, byte('e'), 0, byte('l'), 0, byte('l'), 0, byte('o'), 0, byte(' '), 0,
		byte('W'), 0, byte('o'), 0, byte('r'), 0, byte('l'), 0, byte('d'), 0, byte('!'), 0,
		byte('\n'), 0,
	}

	cpu := cpu.CPU{}
	bus := Bus{}
	ram := RAM{}
	console := console.Console{}

	cpu_pins := &pins.Pins{}
	ram_pins := &pins.Pins{}

	cpu.Reset()
	cpu.Pins = cpu_pins
	ram.init(program)
	ram.Pins = ram_pins

	bus.CPU_Pins = cpu_pins
	bus.RAM_Pins = ram_pins

	var hertz uint16 = 100
	clk := false
	for !cpu.Halt {
		// TODO: add console output
		time.Sleep(time.Second / time.Duration(hertz))
		clk = !clk
		if clk {
			cpu.SetupCycle()
			bus.PropagateCycle()
		} else {
			ram.ProcessCycle()
			bus.ReturnCycle()
			cpu.CompleteCycle()
		}
		console.Step(&ram.memory)
	}

	cpu.DumpReg()
}
