package pins

type Pins struct {
	Address uint16
	Data    uint16
	RW      bool // True = Read, False = Write
	Valid   bool // Whether the bus cycle is active
}
