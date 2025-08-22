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
)

type Token struct {
	Type TokenType
	Text []rune
	Line int
}

func NewToken(t TokenType, text []rune, line int) Token {
	return Token{
		Type: t,
		Text: text,
		Line: line,
	}
}
