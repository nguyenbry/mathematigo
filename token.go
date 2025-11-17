package main

type TokenType int

const (
	OpenParen TokenType = iota
	CloseParen
	Plus
	Minus
	Dot
	Star
	Bang
	BangEq
	Eq
	EqEq
	Lteq
	Gteq
	Lt
	Gt
	Slash
	String
	NewLine
	Number
	Ident
	Semi
	Comma
	Pipe
	Ampersand
)

var ReservedIdentifiers map[string]struct{} = map[string]struct{}{
	"xor": {},
	"and": {},
	"or":  {},
	"not": {},
}

type SmartRune []rune

func (s SmartRune) equals(other SmartRune) bool {
	if len(s) != len(other) {
		return false
	}

	for i, val := range s {
		if val != other[i] {
			return false
		}
	}

	return true
}

var False SmartRune = []rune("false")
var True SmartRune = []rune("true")

type Token struct {
	Type TokenType

	// how it appeared in source code
	Text SmartRune
	// the value itself. best way to think about this is a string without quotes, but .Text will have the quotes
	Literal SmartRune
	Line    int
}

func NewToken(t TokenType, text []rune, line int, literal []rune) Token {
	return Token{
		Type:    t,
		Text:    text,
		Line:    line,
		Literal: literal,
	}
}
