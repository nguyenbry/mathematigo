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
	} else if curr.Type == Ident && curr.Text.equals(True) {
		p.advance()
		out := (Expr)(Literal{
			literal: True,
		})

		return out, nil
	} else if v := SmartRune("null"); curr.Type == Ident && curr.Text.equals(SmartRune("null")) {
		p.advance()

		return Literal{
			literal: v,
		}, nil
	} else if curr.Type == Number {
		p.advance()

		return Literal{
			literal: curr.Literal,
		}, nil
	} else if curr.Type == String {
		p.advance()

		return Literal{
			literal: curr.Literal,
		}, nil
	} else if curr.Type == OpenParen {
		p.advance()

		e, err := p.expression()

		if err != nil {
			return nil, err
		}

		if next, ok := p.peek(); ok && next.Type == CloseParen {
			p.advance()

			return Grouping{content: e}, nil
		} else {
			return nil, errors.New("Expect ')' after expression.")
		}

	}

	return nil, errors.New("TODO")
}

func (p *Parser) expression() (Expr, error) {
	return p.primary()
}

type Expr interface{}

type Literal struct {
	literal []rune
}

type Grouping struct {
	// non-nil
	content Expr
}

var _ Expr = Literal{}
