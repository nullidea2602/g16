package main

import (
	"strings"
)

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

// Tokenize takes an ARM-like assembly source string and returns a slice of tokens.
func Tokenize(source string) []Token {
	var tokens []Token
	// Split the input into lines.
	lines := strings.Split(source, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// Skip empty lines or full-line comments.
		if trimmed == "" || strings.HasPrefix(trimmed, ";") {
			continue
		}
		// Remove inline comments (anything after ';').
		if idx := strings.Index(trimmed, ";"); idx != -1 {
			trimmed = strings.TrimSpace(trimmed[:idx])
		}
		// Replace commas with spaces to separate tokens like "$r0,".
		trimmed = strings.ReplaceAll(trimmed, ",", " ")
		// Split the line by whitespace.
		parts := strings.Fields(trimmed)
		for _, part := range parts {
			token := classifyToken(part)
			tokens = append(tokens, token)
		}
	}
	return tokens
}

// classifyToken examines a token string and determines its type.
func classifyToken(tok string) Token {
	// If the token ends with a colon, it's a label definition.
	if strings.HasSuffix(tok, ":") {
		return Token{Type: TokenLabel, Value: strings.TrimSuffix(tok, ":")}
	}
	// Register direct: starts with '$'
	if strings.HasPrefix(tok, "$") {
		return Token{Type: TokenRegDirect, Value: tok[1:]}
	}
	// Register indirect: starts with '@'
	if strings.HasPrefix(tok, "@") {
		return Token{Type: TokenRegIndirect, Value: tok[1:]}
	}
	// Literal immediates: start with '#' followed by a mode letter.
	if strings.HasPrefix(tok, "#") {
		if len(tok) < 2 {
			return Token{Type: TokenIdentifier, Value: tok}
		}
		mode := tok[1]
		literal := tok[2:]
		switch mode {
		case 'd':
			return Token{Type: TokenImmDec, Value: literal}
		case 'x':
			return Token{Type: TokenImmHex, Value: literal}
		case 'a':
			return Token{Type: TokenImmAscii, Value: literal}
		default:
			return Token{Type: TokenIdentifier, Value: tok}
		}
	}
	// Label address immediates: start with '='.
	if strings.HasPrefix(tok, "=") {
		return Token{Type: TokenLabelImm, Value: tok[1:]}
	}
	// Opcodes: for our example, we support "MOV".
	switch tok {
	case "halt",
		"mov", "movi", "moviu", "movio",
		"inc", "dec", "add", "addi", "sub", "subi", "mul", "div",
		"and", "or", "xor", "not", "shl", "shr",
		"jmp", "je", "jz", "jnz", "jc", "jnc", "call", "ret",
		"push", "pop",
		"nop":
		return Token{Type: TokenOpcode, Value: tok}
	}

	// Otherwise, treat it as an identifier.
	return Token{Type: TokenIdentifier, Value: tok}
}
