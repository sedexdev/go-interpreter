package lexer

import (
	"regexp"

	"github.com/sedexdev/go-interpreter/internal/token"
)

/*
=====================
Define the Lexer type
=====================
*/

// Lexer definition
type Lexer struct {
	program      string
	currentIndex int
	currentChar  byte
}

// CreateLexer creates a new Lexer object
func CreateLexer(program string) *Lexer {
	lexer := &Lexer{program: program}
	lexer.currentChar = program[lexer.currentIndex]
	return lexer
}

/*
===================
Main Lexer Function
===================
*/

// ReadNextToken is a method on the Lexer struct for
// inspecting the current character and then making
// new tokens based on it's value
func (lexer *Lexer) ReadNextToken() token.Token {

	lexer.skipWhitespace()
	advance := true

	var newToken token.Token

	switch lexer.currentChar {
	case '=':
		if lexer.peek() == '=' {
			tok := string(lexer.currentChar) + string(lexer.peek())
			lexer.advance()
			newToken = makeToken("EQUAL", tok)
		} else {
			newToken = makeToken("ASSIGNMENT", string(lexer.currentChar))
		}
	case '!':
		if lexer.peek() == '=' {
			tok := string(lexer.currentChar) + string(lexer.peek())
			lexer.advance()
			newToken = makeToken("NOTEQUAL", tok)
		} else {
			newToken = makeToken("INVALID", string(lexer.currentChar))
		}
	case '<':
		if lexer.peek() == '=' {
			tok := string(lexer.currentChar) + string(lexer.peek())
			lexer.advance()
			newToken = makeToken("LESSTHANEQUAL", tok)
		} else {
			newToken = makeToken("LESSTHAN", string(lexer.currentChar))
		}
	case '>':
		if lexer.peek() == '=' {
			tok := string(lexer.currentChar) + string(lexer.peek())
			lexer.advance()
			newToken = makeToken("GREATERTHANEQUAL", tok)
		} else {
			newToken = makeToken("GREATERTHAN", string(lexer.currentChar))
		}
	case '&':
		if lexer.peek() == '&' {
			tok := string(lexer.currentChar) + string(lexer.peek())
			lexer.advance()
			newToken = makeToken("AND", tok)
		} else {
			newToken = makeToken("INVALID", string(lexer.currentChar))
		}
	case '|':
		if lexer.peek() == '|' {
			tok := string(lexer.currentChar) + string(lexer.peek())
			lexer.advance()
			newToken = makeToken("OR", tok)
		} else {
			newToken = makeToken("INVALID", string(lexer.currentChar))
		}
	case '+':
		newToken = makeToken("PLUS", string(lexer.currentChar))
	case '-':
		newToken = makeToken("MINUS", string(lexer.currentChar))
	case '*':
		newToken = makeToken("MULTIPLY", string(lexer.currentChar))
	case '/':
		newToken = makeToken("DIVIDE", string(lexer.currentChar))
	case '%':
		newToken = makeToken("MODULO", string(lexer.currentChar))
	case '(':
		newToken = makeToken("LEFTPARENTHESES", string(lexer.currentChar))
	case ')':
		newToken = makeToken("RIGHTPARENTHESES", string(lexer.currentChar))
	case '{':
		newToken = makeToken("LEFTCURLYBRACE", string(lexer.currentChar))
	case '}':
		newToken = makeToken("RIGHTCURLYBRACE", string(lexer.currentChar))
	case ',':
		newToken = makeToken("COMMA", string(lexer.currentChar))
	case 0:
		newToken = makeToken("END", "nil")
	default:
		if digit(lexer.currentChar) {
			num := lexer.createNumber()
			newToken = makeToken("INTEGER", num)
			advance = false
		} else if letter(lexer.currentChar) {
			val := lexer.createIdentifier()
			t := token.IsKeyword(val)
			newToken = makeToken(t, val)
			advance = false
		} else {
			newToken = makeToken("INVALID", string(lexer.currentChar))
		}
	}

	if advance {
		lexer.advance()
	}
	return newToken
}

/*
======================================
Matching functions for creating tokens
======================================
*/

// Match whitespace character
func whitespace(char byte) bool {
	spaceChar := regexp.MustCompile(`[\s]`)
	return spaceChar.MatchString(string(char))
}

// Match a number
func digit(char byte) bool {
	number := regexp.MustCompile(`[0-9]`)
	return number.MatchString(string(char))
}

// Match a letter
func letter(char byte) bool {
	letter := regexp.MustCompile(`[a-zA-Z]`)
	return letter.MatchString(string(char))
}

/*
======================
Lexer helper functions
======================
*/

// move on to the next character
func (lexer *Lexer) advance() {
	if lexer.currentIndex >= len(lexer.program)-1 {
		lexer.currentChar = 0
	} else {
		lexer.currentIndex++
		lexer.currentChar = lexer.program[lexer.currentIndex]
	}
}

// Look at the next character in the input
func (lexer *Lexer) peek() byte {
	if lexer.currentIndex+1 >= len(lexer.program) {
		return 0
	}
	return lexer.program[lexer.currentIndex+1]
}

// Skip over whitespace
func (lexer *Lexer) skipWhitespace() {
	for whitespace(lexer.currentChar) {
		lexer.advance()
	}
}

// Create an identifier by analysing the current character
func (lexer *Lexer) createIdentifier() string {
	word := ""
	for letter(lexer.currentChar) {
		word += string(lexer.currentChar)
		lexer.advance()
	}
	return word
}

// Create a number by analysing the current character
func (lexer *Lexer) createNumber() string {
	number := ""
	for digit(lexer.currentChar) {
		number += string(lexer.currentChar)
		lexer.advance()
	}
	return number
}

// Create a new token from the token.Token struct
func makeToken(tokenType string, tokenValue string) token.Token {
	return token.Token{Type: tokenType, Value: tokenValue}
}
