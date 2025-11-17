package main

import (
	"errors"
	"fmt"
	"strconv"
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

func (p *Parser) or() (MathNode, error) {
	// or → and ( "|" and )* ;

	curr, err := p.and()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && next.Type == Pipe; next, ok = p.peek() {
		p.advance() // consume operator

		right, err := p.and()
		if err != nil {
			return nil, err
		}

		curr = OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text)}
	}

	return curr, nil
}

func (p *Parser) and() (MathNode, error) {
	// bitwiseAnd → comparison ( "&" comparison )* ;

	curr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && next.Type == Ampersand; next, ok = p.peek() {
		p.advance() // consume operator

		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		curr = OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text)}
	}

	return curr, nil
}

func (p *Parser) equality() (MathNode, error) {
	// equality → comparison ( ( "!=" | "==" ) comparison )* ;

	curr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && (next.Type == BangEq || next.Type == EqEq); next, ok = p.peek() {
		p.advance() // consume operator

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		curr = OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text)}
	}

	return curr, nil
}

func (p *Parser) comparison() (MathNode, error) {
	// comparison → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;

	curr, err := p.term()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && (next.Type == Gt || next.Type == Gteq || next.Type == Lt || next.Type == Lteq); next, ok = p.peek() {
		p.advance() // consume operator

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		curr = OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text)}
	}

	return curr, nil
}

func (p *Parser) term() (MathNode, error) {
	// term → factor ( ( "-" | "+" ) factor )* ;

	// first time
	curr, err := p.factor()
	if err != nil {
		return nil, err
	}

	// handle ( ( "-" | "+" ) factor )*
	for next, ok := p.peek(); ok && (next.Type == Minus || next.Type == Plus); next, ok = p.peek() {
		p.advance() // consume operator

		right, err := p.factor()

		if err != nil {
			if errors.Is(err, End) {
				return nil, errors.New("TODO message: unended binary")
			} else {
				return nil, err
			}
		} else {
			curr = OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text)}
		}
	}
	return curr, nil
}

func (p *Parser) factor() (MathNode, error) {
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
			curr = OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text)}
		}
	}
	return curr, nil
}

func (p *Parser) implicit() (MathNode, error) {
	curr, err := p.primary()

	if err != nil {
		return nil, err
	}

	for ok := p.canPrimary(); ok; ok = p.canPrimary() {
		right, err := p.primary()

		if err != nil {
			return nil, err
		}

		if c, ok := right.(ConstantNode); ok {
			// return nil, errors.New("cannot have constant node implicit multiplication")
			return nil, fmt.Errorf("unexpected token: \"%v\"", c.String())
		}

		curr = OperatorNode{Args: []MathNode{curr, right}, Op: "*"}
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
		// fmt.Println("canPrimary defaults")
		return false
	}
}

func (p *Parser) primary() (MathNode, error) {
	curr, ok := p.peek()

	if !ok {
		return nil, End
	}

	switch curr.Type {
	case Ident:
		p.advance()

		if curr.Text.equals(False) {
			return BooleanNode(false), nil
		} else if curr.Text.equals(True) {
			return BooleanNode(true), nil
		} else if curr.Text.equals(SmartRune("null")) {
			return NullNode{}, nil
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

					if b, ok := arg.(BlockNode); ok {
						// it is a block
						if len(b.Blocks) != 1 {
							// TODO
							return nil, errors.New("todo special block arg case")
						} else {
							fb = fb.withArg(b.Blocks[0])
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
				return SymbolNode{Name: string(curr.Text)}, nil
			}
		} else {
			// at end, return Symbol?
			return SymbolNode{Name: string(curr.Text)}, nil
		}
	case Number:
		p.advance()

		toParse := curr.Text

		val, err := strconv.ParseFloat(string(toParse), 64)

		if err != nil {
			return nil, err
		}

		out := FloatNode(val)

		return out, nil
	case String:
		p.advance()

		return ConstantNode(string(curr.Literal)), nil
	case OpenParen:
		p.advance()

		e, err := p.expression()

		if err != nil {
			return nil, err
		}

		if next, ok := p.peek(); ok && next.Type == CloseParen {
			p.advance()

			return ParenthesisNode{Content: e}, nil
		} else {
			return nil, errors.New("Expect ')' after expression.")
		}
	default:
		return nil, errors.New("TODO")
	}
}

func (p *Parser) unary() (MathNode, error) {
	if next, ok := p.peek(); ok && next.Type == Bang || next.Type == Minus {
		p.advance()

		content, err := p.unary()

		if err != nil {
			return nil, err
		}

		return OperatorNode{Args: []MathNode{content}, Op: string(next.Text)}, nil
	} else {
		return p.postfix()
	}
}

func (p *Parser) postfix() (MathNode, error) {
	e, err := p.implicit()

	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && next.Type == Bang; next, ok = p.peek() {
		p.advance()
		e = OperatorNode{Args: []MathNode{e}, Op: "!"}
	}

	return e, nil
}

func (p *Parser) Parse() (MathNode, error) {
	return p.expression()
}

func (p *Parser) expression() (MathNode, error) {
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

		b := BlockNode{}

		part, err := p.block()

		if err != nil {
			return nil, err
		}

		b.Blocks = append(b.Blocks, part)

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

			b.Blocks = append(b.Blocks, part)
		}

		return b, nil

	} else {
		part, err := p.block()

		if err != nil {
			return nil, err
		}

		b := BlockNode{Blocks: []MathNode{part}}

		for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
			p.advance() // consume new line

			for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
				p.advance()
			}

			if p.isAtEnd() {
				if len(b.Blocks) == 1 {
					return part, nil
				} else {
					return b, nil
				}
			}

			part, err := p.block()

			if err != nil {
				return nil, err
			}

			b.Blocks = append(b.Blocks, part)
		}

		if len(b.Blocks) == 1 {
			return part, nil
		} else {
			return b, nil
		}
	}
}

func (p *Parser) block() (MathNode, error) {
	return p.or()
}
