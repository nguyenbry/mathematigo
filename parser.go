package mathematigo

import (
	"errors"
	"fmt"
	"strconv"
)

var ErrEnd = errors.New("end of tokens")

type ParseErrType string

const (
	ParseErrEmpty           ParseErrType = "EMPTY"
	ParseErrEnd             ParseErrType = "END"
	ParseErrUnendedFunction ParseErrType = "UNENDED_FUNCTION"
	ParseErrUnexpected      ParseErrType = "UNEXPECTED"
)

type ParseErr struct {
	Type  ParseErrType
	chars []rune
}

var _ error = (*ParseErr)(nil)

func (pe *ParseErr) Error() string {
	switch pe.Type {
	case ParseErrEmpty:
		return "expression is empty"
	case ParseErrEnd:
		return "unexpected end of expression"
	case ParseErrUnexpected:
		return fmt.Sprintf("unexpected token: '%s'", string(pe.chars))
	default:
		return ""
	}
}

var (
	ErrEmptyExpression = &ParseErr{Type: ParseErrEmpty}
	ErrUnexpectedEnd   = &ParseErr{Type: ParseErrEmpty}
	ErrUnendedFunction = &ParseErr{Type: ParseErrUnendedFunction}
)

func newUnexpectedTokenErr(chars []rune) *ParseErr {
	return &ParseErr{
		Type:  ParseErrUnexpected,
		chars: chars,
	}
}

type parser struct {
	tokens  []Token
	current int
}

func newParser(tokens []Token) *parser {
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

func (p *parser) parse() (MathNode, error) {
	first, ok := p.peek()

	if !ok {
		// handles empty expressions
		return nil, ErrEmptyExpression
	}

	b, err := p.expression()

	if err != nil {
		return nil, err
	}

	// did not finish
	if !p.isAtEnd() {
		return nil, ErrUnexpectedEnd
	}

	switch len(b.Blocks) {
	case 0:
		return nil, ErrEmptyExpression
	case 1:
		isLeading := first.Type == NewLine
		isTrailing := p.tokens[p.current-1].Type == NewLine
		if isLeading || isTrailing {
			// keep block
			return b, nil
		} else {
			return b.Blocks[0], nil
		}
	default:
		return b, nil
	}
}

func (p *parser) expression() (*BlockNode, error) {
	// expression → NEWLINE* block (NEWLINE+ block)* NEWLINE*

	b := &BlockNode{}

	for !p.isAtEnd() {
		p.skipNewLines()

		if p.isAtEnd() {
			break
		}

		part, err := p.block()

		if err != nil {
			return nil, err
		}

		b.Blocks = append(b.Blocks, part)
	}

	return b, nil

}

func (p *parser) block() (MathNode, error) {
	return p.or()
}

func (p *parser) or() (MathNode, error) {
	// or → and ( "|" and )* ;

	curr, err := p.and()
	if err != nil {
		return nil, err
	}

	for next, ok := p.peek(); ok && next.Type == Pipe; next, ok = p.peek() {
		p.advance() // consume operator
		p.skipNewLines()

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
		p.skipNewLines()

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
		p.skipNewLines()

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
		p.skipNewLines()

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
			return nil, err
		}

		switch next.Type {
		case Plus:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnAdd}
		case Minus:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnSubtract}
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
			return nil, err
		}

		switch next.Type {
		case Star:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnMultiply}
		case Slash:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnDivide}
		case Mod:
			curr = &OperatorNode{Args: []MathNode{curr, right}, Op: string(next.Text), Fn: OperatorFnMod}
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

	for shouldImplicitMultiply, err := p.canImplicitMultiply(curr); ; shouldImplicitMultiply, err = p.canImplicitMultiply(curr) {
		if err != nil {
			return nil, err
		}
		if !shouldImplicitMultiply {
			break
		}

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
func (p *parser) canImplicitMultiply(leftNode MathNode) (bool, error) {
	right, ok := p.peek()

	if !ok {
		return false, nil
	}

	if right.Type == NewLine {
		return false, nil
	}

	// Strings cannot participate in implicit multiplication
	if right.Type == String {
		return false, errors.New("cannot implicit TODO")
	}

	// If left side is a ConstantNode (string), can't do implicit mult
	if _, ok := leftNode.(*ConstantNode); ok {
		return false, nil
	}

	// Only these tokens can start an implicit multiplication
	switch right.Type {
	case Ident, Number, OpenParen:
		return true, nil
	default:
		return false, nil
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
		if curr.Text.equals(RuneFalse) {
			b := BooleanNode(false)
			return &b, nil
		} else if curr.Text.equals(RuneTrue) {
			b := BooleanNode(true)
			return &b, nil
		} else if curr.Text.equals(RuneNull) {
			n := NullNode{}
			return &n, nil
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
					arg, err := p.block()

					if err != nil {
						return nil, err
					}

					fb = fb.withArg(arg)
					p.skipNewLines()

					next, ok = p.peek() // next iter here because I need value

					if !ok {
						return nil, ErrUnendedFunction
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
						return nil, ErrUnendedFunction
					}
				}

				// if we get here, then the fn hasn't ended
				return nil, ErrUnendedFunction

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
		p.skipNewLines()

		e, err := p.block()

		if err != nil {
			return nil, err
		}

		p.skipNewLines()

		if next, ok := p.peek(); ok && next.Type == CloseParen {
			p.advance()

			return &ParenthesisNode{Content: e}, nil
		} else {
			return nil, errors.New("Expect ')' after expression.")
		}
	default:
		return nil, newUnexpectedTokenErr((p.tokens[p.current].Text))
	}
}

func (p *parser) unary() (MathNode, error) {
	if next, ok := p.peek(); ok && next.Type == Minus {
		p.advance()
		p.skipNewLines()

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
		// do not skip new lines here
		e = &OperatorNode{Args: []MathNode{e}, Op: "!", Fn: OperatorFnFactorial}
	}

	return e, nil
}

func (p *parser) skipNewLines() {
	for next, ok := p.peek(); ok && next.Type == NewLine; next, ok = p.peek() {
		p.advance()
	}
}
