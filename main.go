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
		0x00, byte(MOVRR << 3), // MOVRR = 0x01, 0x01<<3 = 8
	}

	log.Printf("Program: %04X\n", program)

	cpu := CPU{}
	console := Console{}
	cpu.reset()
	var hertz uint16 = 1
	copy(cpu.ram[0x0200:], program)
	log.Printf("RAM 0x0200: %02X\n", cpu.ram[0x0200])
	log.Printf("RAM 0x0201: %02X\n", cpu.ram[0x0201])
	cpu.init()
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
