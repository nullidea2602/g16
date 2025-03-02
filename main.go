package main

import (
	"code/g16/assembler"
	"code/g16/console"
	"code/g16/cpu"
	. "code/g16/isa"
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
	; below is interpreted as movio, meaning mov immediate offset address
	; intermediate is movio, $rl, 0x00
	; cpu executes as mov $rl, ProgramCounter+0x00
	; this keeps opcode, register, and address within 16 bits
	; mov $rx, =label

	mov $r1, =data ; set r1 to address of first character
	mov $r2, #d12 ; number of characters
	mov $r3, =loop ; set r3 to address of loop
	loop: ; address 6
	mov @r0, @r1 ; copy character to stdout @0x0000
	inc $r1 ; advance r1 to address of next character
	dec $r2 ; count down
	jnz $r3, @r2 ; goto loop
	halt
	data: ; address 16
	#aHello_World!
	`
	log.Println(source)

	tokens := assembler.Tokenize(source)

	for _, t := range tokens {
		fmt.Printf("%s: %s\n", t.Type, t.Value)
	}

	assembly := []uint16{
		MOVIO, cpu.R1, 16, // =data
		MOVI, cpu.R2, 12, // number of characters
		MOVIO, cpu.R3, 6, // =loop,
		//loop:
		MOV, II, cpu.R0, cpu.R1,
		INC, cpu.R1,
		DEC, cpu.R2,
		JNZ, cpu.R3, cpu.R2,
		HALT,
		//data:
		//#aHello_World!
		uint16('H'), uint16('e'), uint16('l'), uint16('l'), uint16('o'), uint16('_'),
		uint16('W'), uint16('o'), uint16('r'), uint16('l'), uint16('d'), uint16('!'),
	}

	log.Println(assembly)

	program := []byte{ // Little-endian
		16, byte(MOVIO<<3) | byte(cpu.R1),
		12, byte(MOVI<<3) | byte(cpu.R2),
		6, byte(MOVIO<<3) | byte(cpu.R3),
		byte(cpu.R0<<4) | byte(cpu.R1), byte(MOV<<3) | byte(II),
		byte(cpu.R1), byte(INC << 3),
		byte(cpu.R2), byte(DEC << 3),
		byte(cpu.R3<<4) | byte(cpu.R2), byte(JNZ << 3),
		0, byte(HALT << 3),
		byte('H'), byte('e'), byte('l'), byte('l'), byte('o'), byte('_'),
		byte('W'), byte('o'), byte('r'), byte('l'), byte('d'), byte('!'),
	}

	log.Printf("Program: %04X\n", program)

	cpu := cpu.CPU{}
	cpu.Reset()
	cpu.Load(program)

	console := console.Console{}

	var hertz uint16 = 1
	for !cpu.Halt {
		time.Sleep(time.Second / time.Duration(hertz))
		cpu.Step()
		console.Step(cpu.RAM)
	}

	cpu.DumpReg()
}
