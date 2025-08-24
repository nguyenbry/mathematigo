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
)

var ReservedIdentifiers map[string]struct{} = map[string]struct{}{
	"xor": {},
	"and": {},
	"or":  {},
	"not": {},
}

type Token struct {
	Type TokenType

	// how it appeared in source code
	Text []rune
	// the value itself. best way to think about this is a string without quotes, but .Text will have the quotes
	Literal []rune
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
