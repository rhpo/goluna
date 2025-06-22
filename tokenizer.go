package main

import (
	"fmt"
	"strings"
	"unicode"
)

type TokenType int

const (
	// Literals
	IDENTIFIER TokenType = iota
	STRING
	INT
	FLOAT
	BOOLEAN
	UNDEFINED

	// Keywords
	FN
	LAMBDA
	IF
	ELSE
	RETURN
	TYPEOF
	FOR
	WHILE
	DEBUG
	USE
	OUT

	// Operators
	BINARY_OPERATOR
	EQUALS
	PLUS_EQ
	MINUS_EQ
	EQUALITY_OP
	INEQUALITY_OP
	SMALLER_THAN
	GREATER_THAN
	SMALLER_OR_EQUAL
	GREATER_OR_EQUAL
	AND
	OR
	NEGATION_OP
	INCREMENT
	DECREMENT

	// Punctuation
	COMMA
	DOT
	COLON
	SEMICOLON
	OPEN_PAREN
	CLOSE_PAREN
	OPEN_BRACE
	CLOSE_BRACE
	OPEN_BRACKET
	CLOSE_BRACKET
	TERNARY

	// Special
	NEWLINE
	EOF
)

var keywords = map[string]TokenType{
	"fn":     FN,
	"lambda": LAMBDA,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"typeof": TYPEOF,
	"for":    FOR,
	"while":  WHILE,
	"debug":  DEBUG,
	"use":    USE,
	"out":    OUT,
	"true":   BOOLEAN,
	"false":  BOOLEAN,
	"undef":  UNDEFINED,
}

type Position struct {
	Line   int
	Column int
	Index  int
}

type Token struct {
	Type     TokenType
	Value    string
	Position Position
}

type Tokenizer struct {
	input    []rune
	position int
	line     int
	index    int
}

func NewTokenizer(input string) *Tokenizer {
	return &Tokenizer{
		input:    []rune(input),
		position: 0,
		line:     0,
		index:    0,
	}
}

func (t *Tokenizer) Tokenize() ([]Token, error) {
	var tokens []Token

	for t.position < len(t.input) {
		char := t.current()

		switch {
		case char == '\n':
			tokens = append(tokens, Token{NEWLINE, string(char), Position{t.line, t.index, t.position}})
			t.line++
			t.index = 0
			t.advance()

		case unicode.IsSpace(char):
			t.advance()

		case char == '#':
			// Skip comments
			for t.position < len(t.input) && t.current() != '\n' {
				t.advance()
			}

		case char == '"' || char == '\'':
			startPos := Position{t.line, t.index, t.position}
			str, err := t.readString(char)
			if err != nil {
				return nil, err
			}
			tokens = append(tokens, Token{STRING, str, startPos})

		case unicode.IsDigit(char):
			startPos := Position{t.line, t.index, t.position}
			num, isFloat := t.readNumber()
			tokenType := INT
			if isFloat {
				tokenType = FLOAT
			}
			tokens = append(tokens, Token{tokenType, num, startPos})

		case unicode.IsLetter(char) || char == '_':
			startPos := Position{t.line, t.index, t.position}
			identifier := t.readIdentifier()
			tokenType := IDENTIFIER
			if kw, exists := keywords[identifier]; exists {
				tokenType = kw
			}
			tokens = append(tokens, Token{tokenType, identifier, startPos})

		case char == '(':
			tokens = append(tokens, Token{OPEN_PAREN, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == ')':
			tokens = append(tokens, Token{CLOSE_PAREN, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == '{':
			tokens = append(tokens, Token{OPEN_BRACE, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == '}':
			tokens = append(tokens, Token{CLOSE_BRACE, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == '[':
			tokens = append(tokens, Token{OPEN_BRACKET, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == ']':
			tokens = append(tokens, Token{CLOSE_BRACKET, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == ',':
			tokens = append(tokens, Token{COMMA, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == '.':
			tokens = append(tokens, Token{DOT, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == ':':
			tokens = append(tokens, Token{COLON, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == ';':
			tokens = append(tokens, Token{SEMICOLON, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		case char == '?':
			tokens = append(tokens, Token{TERNARY, string(char), Position{t.line, t.index, t.position}})
			t.advance()

		default:
			if t.isOperator(char) {
				startPos := Position{t.line, t.index, t.position}
				op := t.readOperator()
				tokens = append(tokens, Token{t.getOperatorType(op), op, startPos})
			} else {
				return nil, fmt.Errorf("unexpected character: %c at line %d, column %d", char, t.line, t.index)
			}
		}
	}

	tokens = append(tokens, Token{EOF, "", Position{t.line, t.index, t.position}})
	return tokens, nil
}

func (t *Tokenizer) current() rune {
	if t.position >= len(t.input) {
		return 0
	}
	return t.input[t.position]
}

func (t *Tokenizer) peek() rune {
	if t.position+1 >= len(t.input) {
		return 0
	}
	return t.input[t.position+1]
}

func (t *Tokenizer) advance() {
	if t.position < len(t.input) {
		t.position++
		t.index++
	}
}

func (t *Tokenizer) readString(quote rune) (string, error) {
	t.advance() // Skip opening quote
	var result strings.Builder
	escaped := false

	for t.position < len(t.input) {
		char := t.current()

		if escaped {
			switch char {
			case 'n':
				result.WriteRune('\n')
			case 't':
				result.WriteRune('\t')
			case 'r':
				result.WriteRune('\r')
			case '\\':
				result.WriteRune('\\')
			case '"':
				result.WriteRune('"')
			case '\'':
				result.WriteRune('\'')
			default:
				result.WriteRune(char)
			}
			escaped = false
		} else if char == '\\' {
			escaped = true
		} else if char == quote {
			t.advance() // Skip closing quote
			return result.String(), nil
		} else {
			result.WriteRune(char)
		}

		t.advance()
	}

	return "", fmt.Errorf("unterminated string")
}

func (t *Tokenizer) readNumber() (string, bool) {
	var result strings.Builder
	isFloat := false

	for t.position < len(t.input) && (unicode.IsDigit(t.current()) || t.current() == '.') {
		if t.current() == '.' {
			if isFloat {
				break // Second dot, stop
			}
			isFloat = true
		}
		result.WriteRune(t.current())
		t.advance()
	}

	return result.String(), isFloat
}

func (t *Tokenizer) readIdentifier() string {
	var result strings.Builder

	for t.position < len(t.input) && (unicode.IsLetter(t.current()) || unicode.IsDigit(t.current()) || t.current() == '_') {
		result.WriteRune(t.current())
		t.advance()
	}

	return result.String()
}

func (t *Tokenizer) isOperator(char rune) bool {
	operators := "+-*/%=<>!&|^"
	return strings.ContainsRune(operators, char)
}

func (t *Tokenizer) readOperator() string {
	var result strings.Builder

	for t.position < len(t.input) && t.isOperator(t.current()) {
		result.WriteRune(t.current())
		t.advance()

		// Check for multi-character operators
		op := result.String()
		if len(op) >= 2 {
			switch op {
			case "==", "!=", "<=", ">=", "&&", "||", "++", "--", "+=", "-=", "*=", "/=", "**":
				return op
			}
		}
	}

	return result.String()
}

func (t *Tokenizer) getOperatorType(op string) TokenType {
	switch op {
	case "=":
		return EQUALS
	case "==":
		return EQUALITY_OP
	case "!=":
		return INEQUALITY_OP
	case "<":
		return SMALLER_THAN
	case ">":
		return GREATER_THAN
	case "<=":
		return SMALLER_OR_EQUAL
	case ">=":
		return GREATER_OR_EQUAL
	case "&&":
		return AND
	case "||":
		return OR
	case "!":
		return NEGATION_OP
	case "++":
		return INCREMENT
	case "--":
		return DECREMENT
	case "+=":
		return PLUS_EQ
	case "-=":
		return MINUS_EQ
	default:
		return BINARY_OPERATOR
	}
}
