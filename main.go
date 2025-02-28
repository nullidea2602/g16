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

	program := []byte{
		byte(R0<<4) | byte(RPC), byte(MOV<<3) | byte(DWI), // MOVRR $R0, [$RPC]
		0x33, 0x22,
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
