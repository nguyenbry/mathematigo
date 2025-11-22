package mathematigo

import (
	"errors"
	"strconv"
)

var ErrEnd = errors.New("end of tokens")
var ErrTodoUnendedFunction = errors.New("unended function call")

type parser struct {
	tokens  []Token
	current int
}

func newParser(tokens []Token) *parser {
	if tokens == nil {
		panic("nil tokens")
	}
	return &parser{
		tokens: tokens,
	}
}

func Parse(val string) (MathNode, error) {
	s := NewScanner(val)

	toks, err := s.scanTokens()
	if err != nil {
		return nil, err
	}

	p := newParser(toks)

	ex, err := p.parse()
	if err != nil {
		return nil, err
	}

	return ex, nil
}

func (p *parser) parse() (MathNode, error) {
	out, err := p.expression()

	if err != nil {
		return nil, err
	}

	if !p.isAtEnd() {
		return nil, errors.New("TODO message: unexpected tokens at end")
	}

	return out, nil
}

func (p *parser) isAtEnd() bool {
	return p.current >= len(p.tokens)
}

func (p *parser) peek() (Token, bool) {
	if p.isAtEnd() {
		return Token{}, false
	}

	return p.tokens[p.current], true
}

func (p *parser) advance() Token {
	out := p.tokens[p.current]
	p.current++
	return out
}

func (p *parser) or() (MathNode, error) {
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

		curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnBitOr}
	}

	return curr, nil
}

func (p *parser) and() (MathNode, error) {
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

		curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnBitAnd}
	}

	return curr, nil
}

func (p *parser) equality() (MathNode, error) {
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

		switch next.Type {
		case BangEq:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnUnequal}
		case EqEq:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnEqual}
		}
	}

	return curr, nil
}

func (p *parser) comparison() (MathNode, error) {
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

		switch next.Type {
		case Gt:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnGt}
		case Gteq:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnGteq}
		case Lt:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnLt}
		case Lteq:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnLteq}
		}
	}

	return curr, nil
}

func (p *parser) term() (MathNode, error) {
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
			if errors.Is(err, ErrEnd) {
				return nil, errors.New("TODO message: unended binary")
			} else {
				return nil, err
			}
		} else {

			switch next.Type {
			case Plus:
				curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnAdd}
			case Minus:
				curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnSubtract}
			}
		}
	}
	return curr, nil
}

func (p *parser) factor() (MathNode, error) {
	// factor → power ( ( "/" | "*" | "%" ) power )* ;

	// first time
	curr, err := p.power()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && (next.Type == Slash || next.Type == Star || next.Type == Mod); next, ok = p.peek() {
		p.advance()

		right, err := p.power()

		if err != nil {
			if errors.Is(err, ErrEnd) {
				return nil, errors.New("TODO message: unended binary")
			} else {
				return nil, err
			}
		} else {
			switch next.Type {
			case Star:
				curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnMultiply}
			case Slash:
				curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnDivide}
			case Mod:
				curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnMod}
			}
		}
	}
	return curr, nil
}

func (p *parser) power() (MathNode, error) {
	// power → unary ( "^" power )?

	curr, err := p.unary()
	if err != nil {
		return nil, err
	}

	if next, ok := p.peek(); ok && next.Type == Caret {
		p.advance()

		// RIGHT ASSOCIATIVE — recursive descent
		right, err := p.power()
		if err != nil {
			return nil, err
		}

		curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnPower}
	}

	return curr, nil
}

func (p *parser) implicit() (MathNode, error) {
	curr, err := p.primary()

	if err != nil {
		return nil, err
	}

	for p.canImplicitMultiply(curr) {
		right, err := p.primary()

		if err != nil {
			return nil, err
		}

		curr = &OperatorNode{Args: []MathNode{curr, right}, Op: "*", Fn: OperatorFnMultiply}
	}

	return curr, nil
}

// canImplicitMultiply checks if implicit multiplication can occur
// given the left operand and the next token
func (p *parser) canImplicitMultiply(leftNode MathNode) bool {
	curr, ok := p.peek()

	if !ok {
		return false
	}

	// Strings cannot participate in implicit multiplication
	if curr.Type == String {
		return false
	}

	// If left side is a ConstantNode (string), can't do implicit mult
	if _, ok := leftNode.(*ConstantNode); ok {
		return false
	}

	// Only these tokens can start an implicit multiplication
	switch curr.Type {
	case Ident, Number, OpenParen:
		return true
	default:
		return false
	}
}

func (p *parser) primary() (MathNode, error) {
	curr, ok := p.peek()

	if !ok {
		return nil, ErrEnd
	}

	switch curr.Type {
	case Ident:
		p.advance()
		if curr.Text.equals(False) { b := BooleanNode(false); return &b, nil } else if curr.Text.equals(True) { b := BooleanNode(true); return &b, nil } else if curr.Text.equals(SmartRune("null")) { n := NullNode{}; return &n, nil }
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

					if b, ok := arg.(*BlockNode); ok {
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
				return &SymbolNode{Name: string(curr.Text)}, nil
			}
		} else {
			// at end, return Symbol?
			return &SymbolNode{Name: string(curr.Text)}, nil
		}
	case Number:
		p.advance()

		toParse := curr.Text

		val, err := strconv.ParseFloat(string(toParse), 64)

		if err != nil {
			return nil, err
		}

		out := FloatNode(val)

		return &out, nil
	case String:
		p.advance()
		c := ConstantNode(string(curr.Literal))
		return &c, nil
	case OpenParen:
		p.advance()

		e, err := p.expression()

		if err != nil {
			return nil, err
		}

		if next, ok := p.peek(); ok && next.Type == CloseParen {
			p.advance()

			return &ParenthesisNode{Content: e}, nil
		} else {
			return nil, errors.New("Expect ')' after expression.")
		}
	default:
		return nil, errors.New("unexpected token: " + string(curr.Text))
	}
}

func (p *parser) unary() (MathNode, error) {
	if next, ok := p.peek(); ok && next.Type == Minus {
		p.advance()

		content, err := p.unary()

		if err != nil {
			return nil, err
		}

		return &OperatorNode{Args: []MathNode{content}, Op: string(next.Text), Fn: OperatorFnUnaryMinus}, nil
	} else {
		return p.postfix()
	}
}

func (p *parser) postfix() (MathNode, error) {
	e, err := p.implicit()

	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && next.Type == Bang; next, ok = p.peek() {
		p.advance()
		e = &OperatorNode{Args: []MathNode{e}, Op: "!", Fn: OperatorFnFactorial}
	}

	return e, nil
}

func (p *parser) expression() (MathNode, error) {
	// expression → NEWLINE* block (NEWLINE+ block)* NEWLINE*
	next, ok := p.peek()

	if !ok {
		return nil, ErrEnd
	}

	isLeading := next.Type == NewLine

	if isLeading {
		p.advance()

		for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
			p.advance()
		}

		b := &BlockNode{}

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

		b := &BlockNode{Blocks: []MathNode{part}}

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

func (p *parser) block() (MathNode, error) {
	return p.or()
}
