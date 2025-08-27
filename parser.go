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

	// handle ( ( "-" | "+" ) factor )*
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

				fb := newFunctionNodeBuilder().withFn(string(curr.Text))

				if next, ok := p.peek(); ok && next.Type == CloseParen {
					// simple case myFunc()
					p.advance()
					return fb.build(), nil
				}

				for next, ok = p.peek(); ok; {
					arg, err := p.expression()

					if err != nil {
						return nil, err
					}

					if b, ok := arg.(Block); ok {
						// it is a block
						if len(b.parts) != 1 {
							// TODO
							return nil, errors.New("todo special block arg case")
						} else {
							fb = fb.withArg(b.parts[0])
						}
					} else {
						fb = fb.withArg(arg)
					}

					next, ok = p.peek() // next iter here because I need value

					if !ok {
						return nil, ErrTodoUnendedFunction
					}

					// we have either have expressionComma | expressionClose
					switch next.Type {
					case Comma:
						p.advance()
					case CloseParen:
						p.advance()
						return fb.build(), nil
					default:
						// must be comma or close
						return nil, ErrTodoUnendedFunction
					}

				}

				// if we get here, then the fn hasn't ended
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

		return Unary{content: content, op: next}, nil
	} else {
		return p.postfix()
	}
}

func (p *Parser) postfix() (Expr, error) {
	e, err := p.implicit()

	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && next.Type == Bang; next, ok = p.peek() {
		p.advance()

		e = Unary{op: NewToken(Bang, []rune("!"), -101, nil), content: e}
	}

	return e, nil
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

var _ Expr = Literal{}
var _ Expr = Grouping{}
var _ Expr = Unary{}
var _ Expr = Binary{}
var _ Expr = Symbol{}
var _ Expr = Block{}

// ( "e" ("+" | "-")? digit+ )?
// .12
// 0.12
// .12e
