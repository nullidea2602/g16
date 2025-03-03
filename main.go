package main

import (
	"code/g16/assembler"
	"code/g16/bus"
	"code/g16/console"
	"code/g16/cpu"
	. "code/g16/isa"
	"code/g16/pins"
	"code/g16/ram"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {

	file, err := os.Create("debug.log")
	if err != nil {
		log.Fatalf("failed to create log file: %v", err)
	}
	defer file.Close()

	log.SetOutput(file)

	source := `
	mov $r1, =data ; 00: set r1 to address of data (PC:00 + OFFSET:18)
	mov $r2, #d13  ; 02: set r2 to length of data
	mov $r3, =loop ; 04: set r3 to address of loop (PC:04 + OFFSET:02)
	loop:          ; 06 (not an instruction)
	mov @r0, @r1   ; 06 copy character to stdout @0x0000
	inc $r1        ; 08 advance r1 to address of next character
	inc $r1        ; 10 +2 for next word
	dec $r2        ; 12 count down
	jnz $r3, @r2   ; 14 goto loop
	halt           ; 16 halt
	data:          ; 18 (not an instruction)
	#'Hello World!\n'
	`
	log.Println(source)

	tokens := assembler.Tokenize(source)

	for _, t := range tokens {
		log.Printf("%s: %s\n", t.Type, t.Value)
	}

	intermediate := []uint16{
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

	log.Println(intermediate)

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

	bus := bus.Bus{}
	cpu := cpu.CPU{}
	ram := ram.RAM{}
	console := console.Console{}

	cpu_pins := &pins.Pins{}
	ram_pins := &pins.Pins{}
	console_pins := &pins.Pins{}

	cpu.Reset()
	cpu.Pins = cpu_pins
	ram.Init(program)
	ram.Pins = ram_pins
	console.Pins = console_pins

	bus.CPU_Pins = cpu_pins
	bus.RAM_Pins = ram_pins
	bus.CONSOLE_Pins = console_pins

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
			console.ProcessCycle()
			bus.ReturnCycle()
			cpu.CompleteCycle()
		}
	}

	cpu.DumpReg()
}
