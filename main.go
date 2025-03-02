package main

import (
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
	mov r1, =data ; set r1 to address of first character
	mov r2, #d12 ; number of characters
	mov r3, =loop
	loop:
	mov @r0, @r1 ; copy character to stdout @0x0000
	inc r1 ; advance r1 to address of next character
	dec r2 ; count down
	jnz r3, @r2 ; goto loop
	halt
	data:
	#aHello_World!
	`
	log.Println(source)

	tokens := Tokenize(source)

	for _, t := range tokens {
		log.Printf("%s: %s\n", t.Type, t.Value)
	}

	assembly := []uint16{
		MOVIO, R1, 16, // =data
		MOVI, R2, 12, // number of characters
		MOVIO, R3, 6, // =loop,
		//loop:
		MOV, II, R0, R1, // II
		INC, R1,
		DEC, R2,
		JNZ, R3, R2,
		HALT,
		//data:
		//#aHello_World!
		uint16('H'), uint16('e'), uint16('l'), uint16('l'), uint16('o'), uint16('_'),
		uint16('W'), uint16('o'), uint16('r'), uint16('l'), uint16('d'), uint16('!'),
	}

	log.Println(assembly)

	program := []byte{ // Little-endian
		16, byte(MOVIO<<3) | byte(R1),
		12, byte(MOVI<<3) | byte(R2),
		6, byte(MOVIO<<3) | byte(R3),
		byte(R0<<4) | byte(R1), byte(MOV<<3) | byte(II),
		byte(R1), byte(INC << 3),
		byte(R2), byte(DEC << 3),
		byte(R3<<4) | byte(R2), byte(JNZ << 3),
		0, byte(HALT << 3),
		byte('H'), byte('e'), byte('l'), byte('l'), byte('o'), byte('_'),
		byte('W'), byte('o'), byte('r'), byte('l'), byte('d'), byte('!'),
	}

	log.Printf("Program: %04X\n", program)

	cpu := CPU{}
	console := Console{}
	cpu.reset()
	var hertz uint16 = 1
	copy(cpu.ram[PROGRAM_START:], program)
	for !cpu.halt {
		time.Sleep(time.Second / time.Duration(hertz))
		cpu.step()
		console.step(cpu.ram)
	}

	for i, v := range cpu.reg {
		log.Printf("R%d: %04X\t", i, v)
		if (i+1)%4 == 0 {
			log.Printf("\n")
		}
	}
}
