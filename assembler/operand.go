package assembler

// Define Operand types (you can expand these as needed).
type OperandType int

const (
	OperandRegDirect OperandType = iota
	OperandRegIndirect
	OperandImmDec
	OperandImmHex
	OperandImmAscii
	OperandLabel
	OperandLabelImm
)

type Operand struct {
	Type  OperandType
	Value string
}
