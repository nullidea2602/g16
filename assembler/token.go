package assembler

// TokenType defines the type of token.
type TokenType string

const (
	TokenOpcode      TokenType = "OPCODE"
	TokenRegDirect   TokenType = "REG_DIRECT"
	TokenRegIndirect TokenType = "REG_INDIRECT"
	TokenImmDec      TokenType = "IMMEDIATE_DEC"
	TokenImmHex      TokenType = "IMMEDIATE_HEX"
	TokenImmAscii    TokenType = "IMMEDIATE_ASCII"
	TokenLabelImm    TokenType = "LABEL_IMMEDIATE"
	TokenLabel       TokenType = "LABEL"
	TokenIdentifier  TokenType = "IDENTIFIER"
)

// Token represents a lexical token.
type Token struct {
	Type  TokenType
	Value string
}
