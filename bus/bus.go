package bus

import (
	"code/g16/console"
	"code/g16/pins"
)

type Bus struct {
	CPU_Pins     *pins.Pins
	RAM_Pins     *pins.Pins
	CONSOLE_Pins *pins.Pins
}

func (bus *Bus) PropagateCycle() {
	if bus.CPU_Pins.Valid {
		bus.RAM_Pins.RW = bus.CPU_Pins.RW // even if writing to console, CPUs read/write intent must be updated
		if bus.CPU_Pins.Address == console.CONSOLE_ADDRESS {
			bus.CONSOLE_Pins.Data = bus.CPU_Pins.Data
			bus.CONSOLE_Pins.Valid = true
		} else {
			bus.RAM_Pins.Address = bus.CPU_Pins.Address
			bus.RAM_Pins.Data = bus.CPU_Pins.Data
			bus.RAM_Pins.Valid = true
		}
	} else {
		bus.RAM_Pins.Valid = false
		bus.CONSOLE_Pins.Valid = false
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
