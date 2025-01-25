package main

import (
	"fmt"
	"strings"
)

const END = 0

// --------------------------------------------------------------------------------------------------------------------
// Lexer
// --------------------------------------------------------------------------------------------------------------------

type Lexer struct {
	input                           string
	idx, peek, line, column, length int
	ch                              byte
}

// --------------------------------------------------------------------------------------------------------------------

func newLexer(input string) *Lexer {
	lexer := &Lexer{input: input, line: 1, length: len(input)}
	lexer.readChar()

	return lexer
}

// --------------------------------------------------------------------------------------------------------------------
// Token creation
// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) nextToken() Token {
	l.skipWhitespace()

	switch l.ch {
	case '=':
		return l.makeTwoCharToken(ASSIGN, EQ)
	case '!':
		return l.makeTwoCharToken(BANG, NOTEQ)
	case '>':
		return l.makeTwoCharToken(GT, GTEQ)
	case '<':
		return l.makeTwoCharToken(LT, LTEQ)
	case '+':
		return l.makeToken(PLUS)
	case '-':
		return l.makeToken(MINUS)
	case '*':
		return l.makeToken(ASTERIX)
	case '/':
		return l.makeToken(SLASH)
	case ';':
		return l.makeToken(SEMICOLON)
	case ',':
		return l.makeToken(COMMA)
	case '(':
		return l.makeToken(LPAREN)
	case ')':
		return l.makeToken(RPAREN)
	case '{':
		return l.makeToken(LBRACE)
	case '}':
		return l.makeToken(RBRACE)
	case '[':
		return l.makeToken(LBRACKET)
	case ']':
		return l.makeToken(RBRACKET)
	case '"':
		return l.lexString()
	case END:
		return l.makeToken(EOF)
	default:
		return l.lexOther()
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) lexOther() Token {
	if l.isDigit() {
		return l.lexNumber()
	}
	if l.isLetter() {
		return l.lexIdentKeyword()
	}

	return l.makeToken(ILLEGAL)
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) lexIdentKeyword() Token {
	line := l.line
	col := l.column
	literal := l.readLiteral(l.isLetter)

	return Token{tokenType: lookupIdent(literal), literal: literal, line: line, column: col}
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) lexNumber() Token {
	line := l.line
	col := l.column
	literal := l.readLiteral(l.isDigit)

	if strings.Contains(literal, ".") {
		return Token{tokenType: FLOAT, literal: literal, line: line, column: col}
	}

	return Token{tokenType: INT, literal: literal, line: line, column: col}
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) lexString() Token {
	line := l.line
	col := l.column
	literal := l.readString()
	tok := Token{tokenType: STRING, literal: literal, line: line, column: col}
	l.readChar()

	return tok
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) makeToken(tokenType TokenType) Token {
	tok := Token{tokenType: tokenType, literal: string(l.ch), line: l.line, column: l.column}
	l.readChar()

	return tok
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) makeTwoCharToken(singleType, secondType TokenType) Token {
	line := l.line
	col := l.column
	first := l.ch

	if l.peek < l.length && l.input[l.peek] == '=' {
		l.readChar()
		l.readChar()
		return Token{tokenType: secondType, literal: fmt.Sprintf("%v=", string(first)), line: line, column: col}
	}

	return l.makeToken(singleType)
}

// --------------------------------------------------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) isDigit() bool {
	return '0' <= l.ch && l.ch <= '9' || l.ch == '.'
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) isLetter() bool {
	return 'a' <= l.ch && l.ch <= 'z' || 'A' <= l.ch && l.ch <= 'Z' || l.ch == '_'
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) readLiteral(breakCond func() bool) string {
	start := l.idx
	for breakCond() && l.ch != END {
		l.readChar()
	}

	return l.input[start:l.idx]
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) readString() string {
	position := l.peek

	for {
		l.readChar()
		if l.ch == '"' || l.ch == END {
			break
		}
	}

	return l.input[position:l.idx]
}

// --------------------------------------------------------------------------------------------------------------------
// Advance tokens
// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) readChar() {
	if l.peek >= l.length {
		l.ch = END
	} else {
		l.ch = l.input[l.peek]
	}

	l.idx = l.peek
	l.peek += 1
	l.setLineCol()
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) setLineCol() {
	if l.ch == '\n' {
		l.line += 1
		l.column = 1
	} else {
		l.column += 1
	}
}

// --------------------------------------------------------------------------------------------------------------------

func (l *Lexer) skipWhitespace() {
	for {
		switch l.ch {
		case '\n', ' ', '\t', '\r':
			l.readChar()
		default:
			return
		}
	}
}

// --------------------------------------------------------------------------------------------------------------------
