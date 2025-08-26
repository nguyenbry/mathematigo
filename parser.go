package main

import (
	"errors"
	"fmt"
)

var End = errors.New("end of tokens")
var ErrTodoUnendedFunction = errors.New("unended function call")

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

func (p *Parser) term() (Expr, error) {
	// term → factor ( ( "-" | "+" ) factor )* ;

	// first time
	curr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && (next.Type == Minus || next.Type == Plus); next, ok = p.peek() {
		p.advance()

		right, err := p.factor()

		if err != nil {
			if errors.Is(err, End) {
				return nil, errors.New("TODO message: unended binary")
			} else {
				return nil, err
			}
		} else {
			curr = Binary{left: curr, op: next, right: right}
		}
	}
	return curr, nil
}

func (p *Parser) factor() (Expr, error) {
	// factor → unary ( ( "/" | "*" ) unary )* ;

	// first time
	curr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && (next.Type == Slash || next.Type == Star); next, ok = p.peek() {
		p.advance()

		right, err := p.unary()

		if err != nil {
			if errors.Is(err, End) {
				return nil, errors.New("TODO message: unended binary")
			} else {
				return nil, err
			}
		} else {
			curr = Binary{left: curr, op: next, right: right}
		}
	}
	return curr, nil
}

func (p *Parser) implicit() (Expr, error) {
	curr, err := p.primary()

	if err != nil {
		return nil, err
	}

	for ok := p.canPrimary(); ok; ok = p.canPrimary() {
		right, err := p.primary()

		if err != nil {
			return nil, err
		}
		curr = Binary{left: curr, op: NewToken(Star, []rune("*"), -100, nil), right: right}
	}

	return curr, nil
}

func (p *Parser) canPrimary() bool {
	curr, ok := p.peek()

	if !ok {
		return false
	}

	switch curr.Type {
	case Ident:
		fallthrough
	case Number:
		fallthrough
	case String:
		fallthrough
	case OpenParen:
		return true
	default:
		fmt.Println("canPrimary defaults")
		return false
	}
}

func (p *Parser) primary() (Expr, error) {
	curr, ok := p.peek()

	if !ok {
		return nil, End
	}

	switch curr.Type {
	case Ident:
		p.advance()

		if curr.Text.equals(False) || curr.Text.equals(True) || curr.Text.equals(SmartRune("null")) {
			out := Literal{
				literal: curr.Text,
			}

			return out, nil
		}

		// check if function call
		next, ok := p.peek()

		if ok {
			if next.Type == OpenParen {
				// is function call
				p.advance()

				f := Function{name: curr, args: nil}

				next, ok = p.peek()

				if !ok {
					// "func(" <- no ending
					return nil, ErrTodoUnendedFunction
				} else if next.Type == CloseParen {
					p.advance()
					return f, nil
				}

				for next, ok = p.peek(); ok; next, ok = p.peek() {
					e, err := p.expression()

					if err != nil {
						return nil, err
					}

					f.args = append(f.args, e)

					if commaOrClose, ok := p.peek(); ok {
						switch commaOrClose.Type {
						case Comma:
							p.advance()
						case CloseParen:
							p.advance()
							return f, nil
						}
					} else {
						return nil, ErrTodoUnendedFunction
					}
				}

				return nil, ErrTodoUnendedFunction

			} else {
				// is this right?
				return Symbol{name: curr.Text}, nil
			}

		} else {
			// at end, return Symbol?
			return Symbol{name: curr.Text}, nil
		}

	case Number:
		p.advance()

		return Literal{
			literal: curr.Text,
		}, nil
	case String:
		p.advance()

		return Literal{
			literal: curr.Literal,
		}, nil
	case OpenParen:
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
	default:
		return nil, errors.New("TODO")
	}
}

func (p *Parser) unary() (Expr, error) {
	if next, ok := p.peek(); ok && next.Type == Bang || next.Type == Minus {
		p.advance()

		content, err := p.unary()

		if err != nil {
			return nil, err
		}

		return Unary{content: content}, nil
	} else {
		return p.implicit()
	}
}

func (p *Parser) expression() (Expr, error) {
	// expression → NEWLINE* block (NEWLINE+ block)* NEWLINE*
	next, ok := p.peek()

	if !ok {
		return nil, End
	}

	isLeading := next.Type == NewLine

	if isLeading {
		p.advance()

		for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
			p.advance()
		}

		b := Block{}

		part, err := p.block()

		if err != nil {
			return nil, err
		}

		b.parts = append(b.parts, part)

		// handle (NEWLINE+ block)* NEWLINE*
		for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
			p.advance() // consume new line

			for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
				p.advance()
			}

			if p.isAtEnd() {
				return b, nil
			}

			part, err := p.block()

			if err != nil {
				return nil, err
			}

			b.parts = append(b.parts, part)
		}

		return b, nil

	} else {
		part, err := p.block()

		if err != nil {
			return nil, err
		}

		b := Block{parts: []Expr{part}}

		for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
			p.advance() // consume new line

			for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
				p.advance()
			}

			if p.isAtEnd() {
				if len(b.parts) == 1 {
					return part, nil
				} else {
					return b, nil
				}
			}

			part, err := p.block()

			if err != nil {
				return nil, err
			}

			b.parts = append(b.parts, part)
		}

		if len(b.parts) == 1 {
			return part, nil
		} else {
			return b, nil
		}
	}
}

func (p *Parser) block() (Expr, error) {

	return p.factor()
}

type Expr interface{}

type Literal struct {
	literal []rune
}

type Grouping struct {
	// non-nil
	content Expr
}

type Unary struct {
	op      Token
	content Expr
}

type Binary struct {
	left  Expr
	op    Token
	right Expr
}

type Symbol struct {
	name []rune
}

type Block struct {
	parts []Expr
}

type Function struct {
	name Token
	args []Expr
}

var _ Expr = Literal{}
var _ Expr = Grouping{}
var _ Expr = Unary{}
var _ Expr = Binary{}
var _ Expr = Symbol{}
var _ Expr = Block{}
