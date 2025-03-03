package cpu

import (
	. "code/g16/isa"
	"code/g16/pins"
	"fmt"
	"log"
)

type CPU struct {
	Pins      *pins.Pins
	State     CPUState
	reg       [REGISTER_COUNT]uint16
	op        uint16
	f         uint16
	rx        uint16
	ry        uint16
	i         uint16
	iiPending bool
	up        uint64
	Halt      bool
}

func (cpu *CPU) SetupCycle() {
	log.Printf("Cycle (setup): %d\n", cpu.up)
	cpu.up++

	switch cpu.State {
	case FetchInstruction:
		log.Printf("State: Fetch. Requesting instruction from bus (RW, Valid = true) at %04X\n", cpu.reg[RPC])
		cpu.Pins.Address = cpu.reg[RPC]
		cpu.Pins.RW = true // Read
		cpu.Pins.Valid = true
	case ExecuteInstruction: // Memory not ready
		log.Printf("State: Execute, switching on OP: %04X\n", cpu.op)
		switch cpu.op {
		case HALT:
			cpu.Halt = true
			fmt.Printf("Halted after %d cycles at %04X.\n", cpu.up, cpu.reg[RPC])
		case MOV:
			switch cpu.f {
			case DD: // $RX <- $RY
				log.Printf("Executing MOV DD $%d <- $%d (%04X)", cpu.rx, cpu.ry, cpu.reg[cpu.ry])
				cpu.reg[cpu.rx] = cpu.reg[cpu.ry]
			case DLI, DUI, DWI: // $RX(L/U/W) <- [$RY]
				log.Printf("Executing MOV DXI $%d <- $%d (@%04X)", cpu.rx, cpu.ry, cpu.reg[cpu.ry])
				cpu.Pins.Address = cpu.reg[cpu.ry]
				cpu.Pins.RW = true // Read
				cpu.Pins.Valid = true
			case II: // [$RX] <- [$RY]
				if !cpu.iiPending {
					// First cycle: read from memory at address in RY.
					log.Printf("Executing MOV II (read) $%d <- $%d (@%04X <- @%04X)", cpu.rx, cpu.ry, cpu.reg[cpu.rx], cpu.reg[cpu.ry])
					log.Printf("MOV II: Moving address in RY (%04X) to Bus", cpu.reg[cpu.ry])
					cpu.Pins.Address = cpu.reg[cpu.ry]
					cpu.Pins.RW = true // Read
					cpu.Pins.Valid = true
				} else {
					// Second cycle: write the previously read value to address in RX.
					log.Printf("MOV II: Moving address in RX (%04X) to Bus", cpu.reg[cpu.rx])
					cpu.Pins.Address = cpu.reg[cpu.rx]
					log.Printf("MOV II: Moving data in RTEMP (%04X) to Bus", cpu.reg[RTEMP])
					cpu.Pins.Data = cpu.reg[RTEMP]
					cpu.Pins.RW = false // Write
					cpu.Pins.Valid = true
				}
			default:
				cpu.Halt = true
				fmt.Printf("Panic during MOV execute after %d cycles due to unrecognized FLAG: %02X\n", cpu.up, cpu.f)
			}
		case MOVI:
			log.Printf("Executing MOVI $%d <- #%d", cpu.rx, cpu.i)
			cpu.reg[cpu.rx] = cpu.i
		case MOVIU:
			log.Printf("Executing MOVIU $%d(U) <- #%d", cpu.rx, cpu.i)
			cpu.reg[cpu.rx] = (cpu.reg[cpu.rx] & 0x00FF) | cpu.i<<8
		case MOVIO:
			log.Printf("Executing MOVIO $%d <- #(%04X + %d)", cpu.rx, cpu.reg[RPC], cpu.i)
			cpu.reg[cpu.rx] = cpu.reg[RPC] + cpu.i - 2
		case INC:
			log.Printf("Executing INC $%d", cpu.rx)
			cpu.reg[cpu.rx]++
		case DEC:
			log.Printf("Executing DEC $%d", cpu.rx)
			cpu.reg[cpu.rx]--
		case JZ:
			log.Printf("Executing JZ @%d (@%04X), $%d (%04X)", cpu.rx, cpu.reg[cpu.rx], cpu.ry, cpu.reg[cpu.ry])
			if cpu.reg[cpu.ry] == 0 {
				log.Printf("Jumping")
				cpu.reg[RPC] = cpu.reg[cpu.rx]
			} else {
				log.Printf("Not jumping")
			}
		case JNZ:
			log.Printf("Executing JNZ @%d (@%04X), $%d (%04X)", cpu.rx, cpu.reg[cpu.rx], cpu.ry, cpu.reg[cpu.ry])
			if cpu.reg[cpu.ry] != 0 {
				log.Printf("Jumping")
				cpu.reg[RPC] = cpu.reg[cpu.rx]
			} else {
				log.Printf("Not jumping")
			}
		default:
			cpu.Halt = true
			log.Printf("Panic during execute after %d cycles due to unrecognized OP: %02X\n", cpu.up, cpu.op)
		}

	}
}

func (cpu *CPU) CompleteCycle() {
	log.Printf("Cycle (complete): %d\n", cpu.up)
	cpu.up++

	if cpu.Pins.Valid && cpu.Pins.RW { // Read Operation
		switch cpu.State {
		case FetchInstruction:
			log.Printf("State: Fetch, Reading instruction %04X from bus at %04X\n", cpu.Pins.Data, cpu.reg[RPC])
			cpu.reg[RINS] = cpu.Pins.Data
			cpu.Pins.Valid = false
			cpu.reg[RPC] += BYTES_PER_WORD
			log.Printf("Done fetching, incremented RPC, changing state to Decode\n")
			cpu.State = DecodeInstruction
		case ExecuteInstruction: // Memory ready
			log.Printf("State: Execute (Valid & RW = true), switching on OP: %04X\n", cpu.op)
			switch cpu.op {
			case MOV:
				switch cpu.f {
				case DD:
					log.Printf("Executing MOV DD, Nothing to do... should we be here? Why is Valid true?")
				case DLI:
					log.Printf("Executing MOV DLI $%d(L) <- $%04X", cpu.rx, cpu.Pins.Data)
					cpu.reg[cpu.rx] = (cpu.reg[cpu.rx] & 0xFF00) | cpu.Pins.Data
					cpu.Pins.Valid = false
				case DUI:
					log.Printf("Executing MOV DUI $%d(U) <- $%04X", cpu.rx, cpu.Pins.Data)
					cpu.reg[cpu.rx] = (cpu.reg[cpu.rx] & 0x00FF) | (cpu.Pins.Data << 8)
					cpu.Pins.Valid = false
				case DWI:
					log.Printf("Executing MOV DWI $%d(U) <- $%04X", cpu.rx, cpu.Pins.Data)
					cpu.reg[cpu.rx] = cpu.Pins.Data
					cpu.Pins.Valid = false
				case II:
					if !cpu.iiPending {
						// First cycle (read) complete: store data and mark pending.
						log.Printf("MOV II: Moving data (%04X) from Bus to RTEMP", cpu.Pins.Data)
						cpu.reg[RTEMP] = cpu.Pins.Data
						cpu.Pins.Valid = false
						cpu.iiPending = true
						// Do not transition state; remain in ExecuteInstruction for the write.
					}
				default:
					cpu.Halt = true
					fmt.Printf("Panic during MOV execute after %d cycles due to unrecognized FLAG: %02X\n", cpu.up, cpu.f)
				}
			default:
				fmt.Printf("Warning, didn't implement opcode to handle memory read, %04X", cpu.op)
			}
		}
	} else if !cpu.Pins.Valid && !cpu.Pins.RW {
		switch cpu.op {
		case MOV:
			switch cpu.f {
			case II:
				if cpu.iiPending {
					log.Printf("MOV II: Operation complete")
					cpu.iiPending = false
				}
			}
		}
	}

	// Transition State Machine
	switch cpu.State {
	case DecodeInstruction:
		log.Printf("State: Decode, decoding INS: %04X", cpu.reg[RINS])
		cpu.op = extract(cpu.reg[RINS], OPCODE_OFFSET, OPCODE_WIDTH)
		log.Printf("Decoded OP: %04X\n", cpu.op)

		switch cpu.op {
		case HALT:
			// nothing to do
		case MOV, JZ, JNZ: // $RX $RY
			log.Printf("Decoding MOV/JNZ OP (0x0001/21/22)")
			cpu.f = extract(cpu.reg[RINS], FLAG_OFFSET, FLAG_WIDTH)
			cpu.rx = extract(cpu.reg[RINS], RX_OFFSET, R_WIDTH)
			cpu.ry = extract(cpu.reg[RINS], RY_OFFSET, R_WIDTH)
		case MOVI, MOVIU, MOVIO:
			log.Printf("Decoding MOVI/IU/IO OP (0x0002/3/4)")
			cpu.rx = extract(cpu.reg[RINS], RL_OFFSET, RL_WIDTH)
			cpu.i = extract(cpu.reg[RINS], I_OFFSET, I_WIDTH)
		case INC, DEC:
			log.Printf("Decoding INC/DEC (0x0005/6)")
			cpu.rx = extract(cpu.reg[RINS], RX_OFFSET, R_WIDTH)
		default:
			log.Printf("Panic during decode after %d cycles due to unrecognized OP: %04X\n", cpu.up, cpu.op)
			cpu.op = HALT
		}
		log.Printf("Done decoding, changing state to Execute")
		cpu.State = ExecuteInstruction
	case ExecuteInstruction:
		if cpu.op == MOV && cpu.f == II && cpu.iiPending {
			// Remain in ExecuteInstruction to process the write cycle.
			log.Printf("State: Execute, MOV II: Still in indirect operation; remaining in ExecuteInstruction state.")
		} else {
			log.Printf("Done executing, changing state to Fetch")
			cpu.State = FetchInstruction
		}
	}
}

func (cpu *CPU) Reset() {
	for i := range cpu.reg {
		cpu.reg[i] = 0
	}
	cpu.reg[RSP] = STACK_TOP     // Initialize the Stack Pointer
	cpu.reg[RPC] = PROGRAM_START // Start execution at address 0x0200
	cpu.Halt = false
}

func (cpu *CPU) DumpReg() {
	for i, v := range cpu.reg {
		log.Printf("R%d: %04X\t", i, v)
		if (i+1)%4 == 0 {
			log.Printf("\n")
		}
	}
}
