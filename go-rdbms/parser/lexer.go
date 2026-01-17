package parser

import (
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	// Keywords
	TOKEN_SELECT TokenType = iota
	TOKEN_INSERT
	TOKEN_UPDATE
	TOKEN_DELETE
	TOKEN_CREATE
	TOKEN_TABLE
	TOKEN_FROM
	TOKEN_WHERE
	TOKEN_VALUES
	TOKEN_SET
	TOKEN_INTO
	TOKEN_JOIN
	TOKEN_ON
	TOKEN_PRIMARY
	TOKEN_KEY
	TOKEN_UNIQUE

	// Literals
	TOKEN_IDENTIFIER
	TOKEN_STRING
	TOKEN_NUMBER
	TOKEN_TRUE
	TOKEN_FALSE

	// Operators
	TOKEN_EQUALS
	TOKEN_NOT_EQUALS
	TOKEN_GREATER
	TOKEN_LESS
	TOKEN_GREATER_EQUALS
	TOKEN_LESS_EQUALS
	TOKEN_COMMA
	TOKEN_SEMICOLON
	TOKEN_LEFT_PAREN
	TOKEN_RIGHT_PAREN
	TOKEN_STAR
	TOKEN_DOT

	// End of input
	TOKEN_EOF
)

// Token represents a lexical token
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer performs lexical analysis on SQL input
type Lexer struct {
	input   string
	pos     int
	current byte
}

// NewLexer creates a new lexer for the given input
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// NextToken returns the next token from the input
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.current {
	case '=':
		tok = Token{Type: TOKEN_EQUALS, Literal: "="}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_NOT_EQUALS, Literal: "!="}
		} else {
			tok = Token{Type: TOKEN_IDENTIFIER, Literal: "!"}
		}
	case '>':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_GREATER_EQUALS, Literal: ">="}
		} else {
			tok = Token{Type: TOKEN_GREATER, Literal: ">"}
		}
	case '<':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: TOKEN_LESS_EQUALS, Literal: "<="}
		} else {
			tok = Token{Type: TOKEN_LESS, Literal: "<"}
		}
	case ',':
		tok = Token{Type: TOKEN_COMMA, Literal: ","}
	case ';':
		tok = Token{Type: TOKEN_SEMICOLON, Literal: ";"}
	case '(':
		tok = Token{Type: TOKEN_LEFT_PAREN, Literal: "("}
	case ')':
		tok = Token{Type: TOKEN_RIGHT_PAREN, Literal: ")"}
	case '*':
		tok = Token{Type: TOKEN_STAR, Literal: "*"}
	case '.':
		tok = Token{Type: TOKEN_DOT, Literal: "."}
	case 0:
		tok = Token{Type: TOKEN_EOF, Literal: ""}
	default:
		if isLetter(l.current) {
			tok.Literal = l.readIdentifier()
			tok.Type = l.lookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.current) {
			tok.Type = TOKEN_NUMBER
			tok.Literal = l.readNumber()
			return tok
		} else if l.current == '\'' {
			tok.Type = TOKEN_STRING
			tok.Literal = l.readString()
			return tok
		} else {
			tok = Token{Type: TOKEN_IDENTIFIER, Literal: string(l.current)}
		}
	}

	l.readChar()
	return tok
}

// readChar advances the lexer to the next character
func (l *Lexer) readChar() {
	if l.pos >= len(l.input) {
		l.current = 0
	} else {
		l.current = l.input[l.pos]
	}
	l.pos++
}

// peekChar looks at the next character without advancing
func (l *Lexer) peekChar() byte {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

// skipWhitespace skips over whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.current == ' ' || l.current == '\t' || l.current == '\n' || l.current == '\r' {
		l.readChar()
	}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	pos := l.pos - 1
	for isLetter(l.current) || isDigit(l.current) || l.current == '_' {
		l.readChar()
	}
	return l.input[pos : l.pos-1]
}

// readNumber reads a numeric literal
func (l *Lexer) readNumber() string {
	pos := l.pos - 1
	for isDigit(l.current) {
		l.readChar()
	}
	return l.input[pos : l.pos-1]
}

// readString reads a string literal
func (l *Lexer) readString() string {
	l.readChar() // skip opening quote
	pos := l.pos - 1
	for l.current != '\'' && l.current != 0 {
		l.readChar()
	}
	if l.current == '\'' {
		result := l.input[pos : l.pos-1]
		l.readChar() // skip closing quote
		return result
	}
	return l.input[pos : l.pos-1]
}

// lookupIdent maps keywords to token types
func (l *Lexer) lookupIdent(ident string) TokenType {
	switch strings.ToUpper(ident) {
	case "SELECT":
		return TOKEN_SELECT
	case "INSERT":
		return TOKEN_INSERT
	case "UPDATE":
		return TOKEN_UPDATE
	case "DELETE":
		return TOKEN_DELETE
	case "CREATE":
		return TOKEN_CREATE
	case "TABLE":
		return TOKEN_TABLE
	case "FROM":
		return TOKEN_FROM
	case "WHERE":
		return TOKEN_WHERE
	case "VALUES":
		return TOKEN_VALUES
	case "SET":
		return TOKEN_SET
	case "INTO":
		return TOKEN_INTO
	case "JOIN":
		return TOKEN_JOIN
	case "ON":
		return TOKEN_ON
	case "PRIMARY":
		return TOKEN_PRIMARY
	case "KEY":
		return TOKEN_KEY
	case "UNIQUE":
		return TOKEN_UNIQUE
	case "TRUE":
		return TOKEN_TRUE
	case "FALSE":
		return TOKEN_FALSE
	default:
		return TOKEN_IDENTIFIER
	}
}

// Helper functions
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}
