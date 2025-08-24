package main

import "errors"

type Parser struct {
	tokens  []Token
	current int
}

func NewParser(tokens []Token) *Parser {
	if tokens == nil {
		panic("nil tokens")
	}
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *Parser) peek() (Token, bool) {
	if p.isAtEnd() {
		return Token{}, false
	}

	return p.tokens[p.current], true
}

func (p *Parser) advance() Token {
	out := p.tokens[p.current]
	p.current++
	return out
}

func (p *Parser) primary() (Expr, error) {
	curr, ok := p.peek()

	if !ok {
		panic("what to do")
	}

	if curr.Type == Ident && curr.Text.equals(False) {
		p.advance()
		out := (Expr)(Literal{
			literal: False,
		})

		return out, nil
	}

	return nil, errors.New("TODO")
}

type Expr interface{}

type Literal struct {
	literal []rune
}

var _ Expr = Literal{}
