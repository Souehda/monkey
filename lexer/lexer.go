package lexer

import "monkey/token"

type Lexer struct {
	input               string
	currentPosition     int
	currentReadPosition int // The read position is always one char after currentPosition
	currentChar         byte
}

func (l *Lexer) readChar() {
	l.currentChar = l.peekChar()

	l.currentPosition = l.currentReadPosition
	l.currentReadPosition += 1
}

func (l *Lexer) peekChar() byte {
	if l.currentReadPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.currentReadPosition]
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhiteSpace()

	switch l.currentChar {
	case '+':
		tok = newToken(token.PLUS, l.currentChar)
	case '=':
		if l.peekChar() == '=' {
			ch := l.currentChar
			l.readChar()
			literal := string(ch) + string(l.currentChar)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.currentChar)
		}
	case ',':
		tok = newToken(token.COMMA, l.currentChar)
	case ';':
		tok = newToken(token.SEMICOLON, l.currentChar)
	case '{':
		tok = newToken(token.LBRACE, l.currentChar)
	case '}':
		tok = newToken(token.RBRACE, l.currentChar)
	case '(':
		tok = newToken(token.LPAREN, l.currentChar)
	case ')':
		tok = newToken(token.RPAREN, l.currentChar)
	case '-':
		tok = newToken(token.MINUS, l.currentChar)
	case '/':
		tok = newToken(token.SLASH, l.currentChar)
	case '*':
		tok = newToken(token.ASTERISK, l.currentChar)
	case '<':
		tok = newToken(token.ST, l.currentChar)
	case '>':
		tok = newToken(token.GT, l.currentChar)
	case '!':
		if l.peekChar() == '=' {
			ch := l.currentChar
			l.readChar()
			literal := string(ch) + string(l.currentChar)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.currentChar)
		}
	case 0: // Note for myself. 0 != '0'
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.currentChar) {
			tok.Literal = l.readIdentifierOrNumber(isLetter)
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.currentChar) {
			tok.Literal = l.readIdentifierOrNumber(isDigit)
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.currentChar)
		}
	}
	l.readChar()
	return tok
}

func (l *Lexer) readIdentifierOrNumber(checkFunction func(ch byte) bool) string {
	position := l.currentPosition
	for checkFunction(l.currentChar) {
		l.readChar()
	}
	return l.input[position:l.currentPosition]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) skipWhiteSpace() {
	char := l.currentChar
	for char == ' ' || char == '\t' || char == '\n' || char == '\r' {
		l.readChar()
		char = l.currentChar
	}
}

func newToken(tokenType token.TokenType, character byte) token.Token {
	// For now, we don't mind casting the byte into a string because we only work with utf8.
	// To-do in the future : use runes to handle unicode
	return token.Token{Type: tokenType, Literal: string(character)}
}

func New(input string) *Lexer {
	// Here, currentPosition, currentReadPosition and currentChar will be equal to 0.
	l := &Lexer{input: input}
	l.readChar() // We read a char just after the creation of the Lexer so we already have the first char and the position set correctly
	return l
}
