package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimaryFalse(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	if v, ok := ex.(Literal); ok {
		assert.Equal(t, []rune("false"), v.literal)
	} else {
		t.Fatal()
	}
}

func TestPrimaryConsumes(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, 1, p.current)
}

func asLiteral(t *testing.T, ex Expr, f func(Literal)) {
	if v, ok := ex.(Literal); ok {
		f(v)
	} else {
		t.Fatal("not Literal")
	}
}

func TestMultiplePrimaryCalls(t *testing.T) {
	s := NewScanner(" false true null\n")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	asLiteral(t, ex, func(v Literal) {
		assert.Equal(t, []rune("false"), v.literal)
	})

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	asLiteral(t, ex, func(v Literal) {
		assert.Equal(t, []rune("true"), v.literal)
	})

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	asLiteral(t, ex, func(v Literal) {
		assert.Equal(t, []rune("null"), v.literal)
	})
}

func TestExpression(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	if v, ok := ex.(Literal); ok {
		assert.Equal(t, []rune("false"), v.literal)
	} else {
		t.Fatal()
	}
}

func TestGrouping(t *testing.T) {
	s := NewScanner("(false)")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.Equal(t, Grouping{
		content: Literal{
			literal: []rune("false"),
		},
	}, ex)
}

func TestGroupingErrors(t *testing.T) {
	s := NewScanner("(false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, ex)
	assert.NotNil(t, err)
}

// func TestParseEmpty(t *testing.T) {
// 	s := NewScanner(" ")

// 	toks := s.scanTokens()

// 	p := NewParser(toks)

// 	ex, err := p.expression()

// 	assert.Nil(t, ex)
// 	assert.NotNil(t, err)
// }

func TestParseNoArgsFunction(t *testing.T) {
	s := NewScanner("myFunc()")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "myFunc",
		},
		Args: nil,
	}, ex)
}

func TestParseFunction1Arg(t *testing.T) {
	s := NewScanner("myFunc(2)")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "myFunc",
		},
		Args: []MathNode{Literal{
			literal: []rune("2"),
		}},
	}, ex)
}

func TestImplicitMult(t *testing.T) {
	s := NewScanner("2 a")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Binary{
		left: Literal{
			literal: []rune("2"),
		},
		op: NewToken(Star, []rune("*"), -100, nil),
		right: Symbol{
			name: []rune("a"),
		},
	},
		ex)
}

func TestImplicitMult2(t *testing.T) {
	s := NewScanner("1a 2")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Binary{
		left: Binary{
			left: Literal{
				literal: []rune("1"),
			},
			op: NewToken(Star, []rune("*"), -100, nil),
			right: Symbol{
				name: []rune("a"),
			},
		},
		op: NewToken(Star, []rune("*"), -100, nil),
		right: Literal{
			literal: []rune("2"),
		},
	},
		ex)
}

func TestBlockSimple(t *testing.T) {
	s := NewScanner("2 \n a")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Block{
		parts: []Expr{Literal{
			literal: []rune("2"),
		},
			Symbol{
				name: []rune("a"),
			}},
	},
		ex)
}

func TestTrailingNewLinesDoesNotProduceBlock(t *testing.T) {
	s := NewScanner("2 a \n\n\n")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Binary{
		left: Literal{
			literal: []rune("2"),
		},
		op: NewToken(Star, []rune("*"), -100, nil),
		right: Symbol{
			name: []rune("a"),
		},
	},
		ex)
}

func TestLeadingNewLinesProducesBlock(t *testing.T) {
	s := NewScanner("\n2 a")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Block{
		parts: []Expr{
			Binary{
				left: Literal{
					literal: []rune("2"),
				},
				op: NewToken(Star, []rune("*"), -100, nil),
				right: Symbol{
					name: []rune("a"),
				},
			},
		},
	},
		ex)
}

func TestMultipleBlocksWithFunctionCall(t *testing.T) {
	s := NewScanner("\n2 a\nmyFunc(2)")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Block{
		parts: []Expr{
			Binary{
				left: Literal{
					literal: []rune("2"),
				},
				op: NewToken(Star, []rune("*"), -100, nil),
				right: Symbol{
					name: []rune("a"),
				},
			},
			FunctionNode{Fn: SymbolNode{
				Name: ("myFunc"),
			},
				Args: []MathNode{Literal{
					literal: []rune("2"),
				}},
			},
		},
	},
		ex)
}

func TestMultipleBlocksWithFunctionCallAndAddition(t *testing.T) {
	s := NewScanner("\n2 a\nmyFunc(2) * 2")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Block{
		parts: []Expr{
			Binary{
				left: Literal{
					literal: []rune("2"),
				},
				op: NewToken(Star, []rune("*"), -100, nil),
				right: Symbol{
					name: []rune("a"),
				},
			},
			Binary{
				left: FunctionNode{Fn: SymbolNode{
					Name: "myFunc",
				},
					Args: []MathNode{Literal{
						literal: []rune("2"),
					}},
				},
				op: NewToken(Star, []rune("*"), 2, nil),
				right: Literal{
					literal: []rune("2"),
				},
			},
		},
	},
		ex)
}

func TestNewLineInFunctionArgs(t *testing.T) {
	s := NewScanner("myFunc(2 \n 3)")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.NotNil(t, err)
	assert.Nil(t, ex)
}

func TestFactorial(t *testing.T) {
	s := NewScanner("a!")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Unary{
		op: Token{
			Type: Bang,
			Text: []rune("!"),
			Line: -101,
		},
		content: Symbol{
			name: []rune("a"),
		},
	}, ex)
}

func TestFactorialAndUnaryMinusPrecedence(t *testing.T) {
	s := NewScanner("-a!")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, Unary{
		op: Token{
			Type: Minus,
			Text: []rune("-"),
			Line: 0,
		},
		content: Unary{
			op: Token{
				Type: Bang,
				Text: []rune("!"),
				Line: -101,
			},
			content: Symbol{
				name: []rune("a"),
			},
		},
	}, ex)
}
