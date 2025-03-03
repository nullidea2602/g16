package assembler

import (
	"fmt"
)

func tokenTypeToOperandType(tok TokenType) (OperandType, error) {
	switch tok {
	case TokenRegDirect:
		return OperandRegDirect, nil
	case TokenRegIndirect:
		return OperandRegIndirect, nil
	case TokenImmDec:
		return OperandImmDec, nil
	case TokenImmHex:
		return OperandImmHex, nil
	case TokenImmAscii:
		return OperandImmAscii, nil
	case TokenLabelImm:
		return OperandLabelImm, nil
	case TokenLabel:
		return OperandLabel, nil
	default:
		return 0, fmt.Errorf("unsupported token type: %s\n", tok)
	}
}

// DataItem represents a data section item.
type DataItem struct {
	Label   string
	Data    string
	Address int
}

type Instruction struct {
	Address  int
	Opcode   string
	Operands []Operand
}

// Parser holds the list of tokens and a pointer to the current position.
type Parser struct {
	tokens        []Token
	pos           int
	instructions  []Instruction
	dataItems     []DataItem
	symbolTable   map[string]int
	addrCounter   int
	inDataSection bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:      tokens,
		symbolTable: make(map[string]int),
	}
}

// Parse performs the two-pass parsing.
func (p *Parser) Parse() ([]Instruction, []DataItem, error) {
	// First pass: build instructions and record label addresses.
	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		// If we're in the data section, only accept data tokens.
		if p.inDataSection {
			switch token.Type {
			case TokenImmAscii:
				// There should be a preceding data label.
				if len(p.dataItems) == 0 {
					return nil, nil, fmt.Errorf("data token found without a preceding data label at position %d", p.pos)
				}
				// Append the data to the most recent DataItem.
				lastIndex := len(p.dataItems) - 1
				p.dataItems[lastIndex].Data += token.Value
				// Optionally, update the address counter based on data length.
				p.addrCounter += len(token.Value)
				p.pos++
			default:
				return nil, nil, fmt.Errorf("unexpected token in data section %v at position %d", token, p.pos)
			}
			continue
		}

		// Not in data section.
		switch token.Type {
		case TokenLabel:
			// Look ahead to see if this label is followed by a data token.
			if p.pos+1 < len(p.tokens) && p.tokens[p.pos+1].Type == TokenImmAscii {
				// This marks the start of a data section.
				p.symbolTable[token.Value] = p.addrCounter
				p.dataItems = append(p.dataItems, DataItem{
					Label:   token.Value,
					Data:    "", // will be filled by following data tokens.
					Address: p.addrCounter,
				})
				p.inDataSection = true // switch to data mode.
				p.pos++                // consume the label token.
			} else {
				// Otherwise, it's a normal label for instructions.
				p.symbolTable[token.Value] = p.addrCounter
				p.pos++
			}
		case TokenOpcode:
			inst, err := p.parseInstruction()
			if err != nil {
				return nil, nil, err
			}
			inst.Address = p.addrCounter
			p.addrCounter++ // For simplicity, assume each instruction is 1 word.
			p.instructions = append(p.instructions, inst)
		default:
			return nil, nil, fmt.Errorf("unexpected token %v at position %d", token, p.pos)
		}
	}

	// Second pass: resolve label references in instruction operands.
	for i, inst := range p.instructions {
		for j, op := range inst.Operands {
			// Check for label operands that need resolution.
			labelImm, _ := tokenTypeToOperandType(TokenLabelImm)
			label, _ := tokenTypeToOperandType(TokenLabel)
			if op.Type == labelImm || op.Type == label {
				addr, ok := p.symbolTable[op.Value]
				if !ok {
					return nil, nil, fmt.Errorf("undefined label: %s", op.Value)
				}
				p.instructions[i].Operands[j].Value = fmt.Sprintf("%d", addr)
			}
		}
	}
	return p.instructions, p.dataItems, nil
}

// parseInstruction builds an instruction from an opcode and its operands.
func (p *Parser) parseInstruction() (Instruction, error) {
	inst := Instruction{
		Operands: []Operand{},
	}

	// Expect an opcode.
	opToken := p.tokens[p.pos]
	inst.Opcode = opToken.Value
	p.pos++

	// Read operands until we reach an opcode or a label token.
	for p.pos < len(p.tokens) {
		tok := p.tokens[p.pos]
		if tok.Type == TokenOpcode || tok.Type == TokenLabel {
			break
		}
		tokType, _ := tokenTypeToOperandType(tok.Type)
		operand := Operand{
			// For simplicity, we reuse the token type.
			// In a more refined parser, you might map token types to OperandType explicitly.
			Type:  tokType,
			Value: tok.Value,
		}
		inst.Operands = append(inst.Operands, operand)
		p.pos++
	}
	return inst, nil
}
