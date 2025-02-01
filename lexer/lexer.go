package lexer

import "monkey/token"

// define the Lexer struct (object)
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// Create a pointer to a Lexer instance to an input string and assign its input property a value
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// Create a method received by the pointer to a Lexer instance to read the current input char and increment the lexer position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// create a method received by the pointer to a Lexer instance to look at the next input char
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}

// a function to create a new Token struct from a char in the input string
// takes a tokenType and a char and returns a Token struct with (Type: TokenType Literal:string(ch))
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

// method on the pointer to the Lexer instance which if current char is a Letter
// continues reading and updating the reader as long as its a Letter "a-z" or "A-Z"
// if that value is a letter then it call readChar() to increment the postions being read
// returns a slice of the input from the current position to the new position
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// func which accepts a char and returns a bool indicating if the char is a letter
// NOTE: we also include '_' as a letter to make names like foo_bar
// NOTE: if we wanted to include other letters in our named identifiers/keywords we could add here ('?', '!', etc.)
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// function to strip white space from source code so that it isn't read as a char
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// Note: there are a lot of ways to improve this further
// Currently we only check for 0-9; but could extend this to include floats, numbers in hex notation, etc.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// function used by readNumber() and NextToken() to determine if char is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// Create a method received by the pointer to a Lexer instance to interpret the current char (byte) and return a Token struct (object) representing that char
// NOTE: token.Token is the syntax to refer to the Token type from the "token" package (i.e. "monkey/token")
// REMEMBER: Token is a struct which has properties: Type: TokenType and Literal: string

// NOTE: Refactor so that each "ch" is mapped to its type (i.e. "=" : "EQ") so the NextToken function can just lookup the type in the map without a switch statement
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			literal := string(ch) + string(l.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.STAR, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}
