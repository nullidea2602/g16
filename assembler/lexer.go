package assembler

import (
	"strings"
	"unicode"
)

// Tokenize takes an ARM-like assembly source string and returns a slice of tokens.
func Tokenize(source string) []Token {
	var tokens []Token
	// Split the input into lines.
	lines := strings.SplitSeq(source, "\n")
	for line := range lines {
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

		// Use a character-based scan to handle string literals that may include spaces.
		i := 0
		for i < len(trimmed) {
			// Skip whitespace.
			if unicode.IsSpace(rune(trimmed[i])) {
				i++
				continue
			}

			// Check if we have a string literal starting with "#'"
			if i+1 < len(trimmed) && trimmed[i] == '#' && trimmed[i+1] == '\'' {
				i += 2 // skip "#'"
				start := i
				// Read until the closing single quote.
				for i < len(trimmed) && trimmed[i] != '\'' {
					i++
				}
				if i < len(trimmed) && trimmed[i] == '\'' {
					// Extract the string literal value.
					tokenVal := trimmed[start:i]
					tokens = append(tokens, Token{Type: TokenImmAscii, Value: tokenVal})
					i++ // skip the closing quote
				} else {
					// Unterminated string literal: take rest of line.
					tokenVal := trimmed[start:]
					tokens = append(tokens, Token{Type: TokenImmAscii, Value: tokenVal})
					break
				}
			} else {
				// Otherwise, grab a non-string token until next whitespace.
				start := i
				for i < len(trimmed) && !unicode.IsSpace(rune(trimmed[i])) {
					i++
				}
				tokenStr := trimmed[start:i]
				tokens = append(tokens, classifyToken(tokenStr))
			}
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
	// Literal immediates that aren’t string literals.
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
		// Optionally remove the old string syntax if not used anymore.
		// case 'a':
		//	return Token{Type: TokenImmAscii, Value: literal}
		default:
			return Token{Type: TokenIdentifier, Value: tok}
		}
	}
	// Label address immediates: start with '='.
	if strings.HasPrefix(tok, "=") {
		return Token{Type: TokenLabelImm, Value: tok[1:]}
	}
	// Opcodes: for our example, we support "mov" and similar.
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
