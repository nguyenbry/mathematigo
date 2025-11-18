package mathematigo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimaryFalse(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BooleanNode(false), ex)
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

func TestMultiplePrimaryCalls(t *testing.T) {
	s := NewScanner(" false true null\n")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	require.Equal(t, BooleanNode(false), ex)

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	require.Equal(t, BooleanNode(true), ex)

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	require.Equal(t, NullNode{}, ex)
}

func TestExpression(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BooleanNode(false), ex)
}

func TestGrouping(t *testing.T) {
	s := NewScanner("(false)")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.Equal(t, ParenthesisNode{
		Content: BooleanNode(false),
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

	require.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "myFunc",
		},
		Args: []MathNode{FloatNode(float64(2))},
	}, ex)
}

func TestImplicitMult(t *testing.T) {
	s := NewScanner("2 a")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			FloatNode(float64(2)),
			SymbolNode{Name: "a"},
		},
		Op: "*",
	},
		ex)
}

func TestImplicitMult2(t *testing.T) {
	s := NewScanner("1a 2")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(1)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
			},
			FloatNode(float64(2)),
		},
		Op: "*",
	},
		ex)
}

func TestBlockSimple(t *testing.T) {
	s := NewScanner("2 \n a")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			FloatNode(float64(2)),
			SymbolNode{Name: "a"},
		},
	},
		ex)
}

func TestTrailingNewLinesDoesNotProduceBlock(t *testing.T) {
	s := NewScanner("2 a \n\n\n")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			FloatNode(float64(2)),
			SymbolNode{Name: "a"},
		},
		Op: "*",
	},
		ex)
}

func TestLeadingNewLinesProducesBlock(t *testing.T) {
	s := NewScanner("\n2 a")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(2)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
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

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(2)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
			},
			FunctionNode{Fn: SymbolNode{
				Name: "myFunc",
			},
				Args: []MathNode{FloatNode(float64(2))},
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

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(2)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
			},
			OperatorNode{
				Args: []MathNode{
					FunctionNode{Fn: SymbolNode{
						Name: "myFunc",
					},
						Args: []MathNode{FloatNode(float64(2))},
					},
					FloatNode(float64(2)),
				},
				Op: "*",
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

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "a"},
		},
		Op: "!",
	}, ex)
}

func TestFactorialAndUnaryMinusPrecedence(t *testing.T) {
	s := NewScanner("-a!")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "a"},
				},
				Op: "!",
			},
		},
		Op: "-",
	}, ex)
}

func TestParseFunctionMultipleArgs(t *testing.T) {
	s := NewScanner("myFunc(2, 3, x)")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "myFunc",
		},
		Args: []MathNode{
			FloatNode(float64(2)),
			FloatNode(float64(3)),
			SymbolNode{Name: "x"},
		},
	}, ex)
}

func TestParseAmpersandBindsTighterThanPipe(t *testing.T) {
	s := NewScanner("a | b & c")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "a"},
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "b"},
					SymbolNode{Name: "c"},
				},
				Op: "&",
			},
		},
		Op: "|",
	}, ex)

	s = NewScanner("a & b | c")

	p = NewParser(s.scanTokens())

	ex, err = p.expression()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "a"},
					SymbolNode{Name: "b"},
				},
				Op: "&",
			},
			SymbolNode{Name: "c"},
		},
		Op: "|",
	}, ex)
}

func TestParsePipe(t *testing.T) {
	s := NewScanner("(\"X\" | \"y\")")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, ParenthesisNode{
		Content: OperatorNode{
			Args: []MathNode{
				ConstantNode("X"),
				ConstantNode("y"),
			},
			Op: "|",
		},
	}, ex)

}

func TestParseDoubleEqualsTighterThanAmpersand(t *testing.T) {
	s := NewScanner("a == b & c")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "a"},
					SymbolNode{Name: "b"},
				},
				Op: "==",
			},
			SymbolNode{Name: "c"},
		},
		Op: "&",
	}, ex)

	s = NewScanner("a & b == c")

	p = NewParser(s.scanTokens())

	ex, err = p.expression()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "a"},
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "b"},
					SymbolNode{Name: "c"},
				},
				Op: "==",
			},
		},
		Op: "&",
	}, ex)
}

func TestImplicitMultiplicationCannotHappenForConstantNode(t *testing.T) {
	s := NewScanner(`x "abc"`)

	p := NewParser(s.scanTokens())

	ex, err := p.Parse()

	if !assert.Error(t, err) {
		fmt.Printf("error was nil, ex was: %+v\n", ex)
		fmt.Printf("ex type: %T\n", ex)
		t.FailNow()
	}
	require.Error(t, err)
	require.Nil(t, ex)

	s = NewScanner(`x"abc"`)
	p = NewParser(s.scanTokens())
	ex, err = p.Parse()

	require.Error(t, err)
	require.Nil(t, ex)

	s = NewScanner(`x "abc" * 2`)
	p = NewParser(s.scanTokens())
	ex, err = p.Parse()

	require.Error(t, err)
	require.Nil(t, ex)

	s = NewScanner(`x ("abc")`)

	p = NewParser(s.scanTokens())

	ex, err = p.Parse()

	require.NoError(t, err)
	require.NotNil(t, ex)

	require.Equal(t, FunctionNode{
		Fn: SymbolNode{Name: "x"},
		Args: []MathNode{
			ConstantNode("abc"),
		},
	}, ex)
}

func TestParseImplicitMult(t *testing.T) {
	s := NewScanner("(1+2)(3+4)")

	p := NewParser(s.scanTokens())

	ex, err := p.Parse()

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			ParenthesisNode{
				Content: OperatorNode{
					Args: []MathNode{
						FloatNode(1),
						FloatNode(2),
					},
					Op: "+",
				},
			},
			ParenthesisNode{
				Content: OperatorNode{
					Args: []MathNode{
						FloatNode(3),
						FloatNode(4),
					},
					Op: "+",
				},
			},
		},
		Op: "*",
	}, ex)
}

func TestParseComplexFunction(t *testing.T) {
	s := NewScanner("pattern_match(\"x\", \"2023-12-23 15:41\", \"2024-02-21 23:59\") >= 1")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()
	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			FunctionNode{
				Fn: SymbolNode{Name: "pattern_match"},
				Args: []MathNode{
					ConstantNode("x"),
					ConstantNode("2023-12-23 15:41"),
					ConstantNode("2024-02-21 23:59"),
				},
			},
			FloatNode(1),
		},
		Op: ">=",
	}, ex)
}

func TestParseComplexNestedFunction(t *testing.T) {
	s := NewScanner("funky(\"y\", concat(\"2023-12-23 \", \"15:41\"), concat(\"2024-02-21 \", \"23:59\")) >= 1")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()
	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			FunctionNode{
				Fn: SymbolNode{Name: "funky"},
				Args: []MathNode{
					ConstantNode("y"),
					FunctionNode{
						Fn: SymbolNode{Name: "concat"},
						Args: []MathNode{
							ConstantNode("2023-12-23 "),
							ConstantNode("15:41"),
						},
					},
					FunctionNode{
						Fn: SymbolNode{Name: "concat"},
						Args: []MathNode{
							ConstantNode("2024-02-21 "),
							ConstantNode("23:59"),
						},
					},
				},
			},
			FloatNode(1),
		},
		Op: ">=",
	}, ex)
}

func TestParseScientificNumber(t *testing.T) {
	s := NewScanner("9e+10")

	p := NewParser(s.scanTokens())

	ex, err := p.expression()
	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, FloatNode(9e10), ex)
}
