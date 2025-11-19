package mathematigo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimaryFalse(t *testing.T) {
	s := NewScanner(" false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BooleanNode(false), ex)
}

func TestPrimaryConsumes(t *testing.T) {
	s := NewScanner(" false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, 1, p.current)
}

func TestMultiplePrimaryCalls(t *testing.T) {
	s := NewScanner(" false true null\n")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

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

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BooleanNode(false), ex)
}

func TestGrouping(t *testing.T) {
	s := NewScanner("(false)")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.Equal(t, ParenthesisNode{
		Content: BooleanNode(false),
	}, ex)
}

func TestGroupingErrors(t *testing.T) {
	s := NewScanner("(false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

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
	ex, err := Parse("myFunc()")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "myFunc",
		},
		Args: nil,
	}, ex)
}

func TestParseFunction1Arg(t *testing.T) {
	ex, err := Parse("myFunc(2)")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "myFunc",
		},
		Args: []MathNode{FloatNode(float64(2))},
	}, ex)
}

func TestImplicitMult(t *testing.T) {
	ex, err := Parse("2 a")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			FloatNode(float64(2)),
			SymbolNode{Name: "a"},
		},
		Op: "*",
		Fn: OperatorFnMultiply,
	},
		ex)
}

func TestImplicitMult2(t *testing.T) {
	ex, err := Parse("1a 2")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(1)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
				Fn: OperatorFnMultiply,
			},
			FloatNode(float64(2)),
		},
		Op: "*",
		Fn: OperatorFnMultiply,
	},
		ex)
}

func TestBlockSimple(t *testing.T) {
	ex, err := Parse("2 \n a")

	require.NoError(t, err)
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
	ex, err := Parse("2 a \n\n\n")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			FloatNode(float64(2)),
			SymbolNode{Name: "a"},
		},
		Op: "*",
		Fn: OperatorFnMultiply,
	},
		ex)
}

func TestLeadingNewLinesProducesBlock(t *testing.T) {
	ex, err := Parse("\n2 a")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(2)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
				Fn: OperatorFnMultiply,
			},
		},
	},
		ex)
}

func TestMultipleBlocksWithFunctionCall(t *testing.T) {
	ex, err := Parse("\n2 a\nmyFunc(2)")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(2)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
				Fn: OperatorFnMultiply,
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
	ex, err := Parse("\n2 a\nmyFunc(2) * 2")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, BlockNode{
		Blocks: []MathNode{
			OperatorNode{
				Args: []MathNode{
					FloatNode(float64(2)),
					SymbolNode{Name: "a"},
				},
				Op: "*",
				Fn: OperatorFnMultiply,
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
				Fn: OperatorFnMultiply,
			},
		},
	},
		ex)
}

func TestNewLineInFunctionArgs(t *testing.T) {
	ex, err := Parse("myFunc(2 \n 3)")

	assert.NotNil(t, err)
	assert.Nil(t, ex)
}

func TestFactorial(t *testing.T) {
	ex, err := Parse("a!")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "a"},
		},
		Op: "!",
		Fn: OperatorFnFactorial,
	}, ex)
}

func TestFactorialAndUnaryMinusPrecedence(t *testing.T) {
	ex, err := Parse("-a!")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "a"},
				},
				Op: "!",
				Fn: OperatorFnFactorial,
			},
		},
		Op: "-",
		Fn: OperatorFnUnaryMinus,
	}, ex)
}

func TestParseFunctionMultipleArgs(t *testing.T) {
	ex, err := Parse("myFunc(2, 3, x)")

	require.NoError(t, err)
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
	ex, err := Parse("a | b & c")

	require.NoError(t, err)
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
				Fn: OperatorFnBitAnd,
			},
		},
		Op: "|",
		Fn: OperatorFnBitOr,
	}, ex)

	ex, err = Parse("a & b | c")

	require.NoError(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "a"},
					SymbolNode{Name: "b"},
				},
				Op: "&",
				Fn: OperatorFnBitAnd,
			},
			SymbolNode{Name: "c"},
		},
		Op: "|",
		Fn: OperatorFnBitOr,
	}, ex)
}

func TestParsePipe(t *testing.T) {
	ex, err := Parse("(\"X\" | \"y\")")

	require.NoError(t, err)
	require.NotNil(t, ex)

	require.Equal(t, ParenthesisNode{
		Content: OperatorNode{
			Args: []MathNode{
				ConstantNode("X"),
				ConstantNode("y"),
			},
			Op: "|",
			Fn: OperatorFnBitOr,
		},
	}, ex)

}

func TestParseDoubleEqualsTighterThanAmpersand(t *testing.T) {
	ex, err := Parse("a == b & c")

	require.NoError(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			OperatorNode{
				Args: []MathNode{
					SymbolNode{Name: "a"},
					SymbolNode{Name: "b"},
				},
				Op: "==",
				Fn: OperatorFnEqual,
			},
			SymbolNode{Name: "c"},
		},
		Op: "&",
		Fn: OperatorFnBitAnd,
	}, ex)

	ex, err = Parse("a & b == c")

	require.NoError(t, err)
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
				Fn: OperatorFnEqual,
			},
		},
		Op: "&",
		Fn: OperatorFnBitAnd,
	}, ex)
}

func TestImplicitMultiplicationCannotHappenForConstantNode(t *testing.T) {
	ex, err := Parse(`x "abc"`)

	if !assert.Error(t, err) {
		fmt.Printf("error was nil, ex was: %+v\n", ex)
		fmt.Printf("ex type: %T\n", ex)
		t.FailNow()
	}
	require.Error(t, err)
	require.Nil(t, ex)

	ex, err = Parse(`x"abc"`)

	require.Error(t, err)
	require.Nil(t, ex)

	ex, err = Parse(`x "abc" * 2`)

	require.Error(t, err)
	require.Nil(t, ex)

	ex, err = Parse(`x ("abc")`)

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
	ex, err := Parse("(1+2)(3+4)")

	require.NoError(t, err)
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
					Fn: OperatorFnAdd,
				},
			},
			ParenthesisNode{
				Content: OperatorNode{
					Args: []MathNode{
						FloatNode(3),
						FloatNode(4),
					},
					Op: "+",
					Fn: OperatorFnAdd,
				},
			},
		},
		Op: "*",
		Fn: OperatorFnMultiply,
	}, ex)
}

func TestParseComplexFunction(t *testing.T) {
	ex, err := Parse("pattern_match(\"x\", \"2023-12-23 15:41\", \"2024-02-21 23:59\") >= 1")

	require.NoError(t, err)
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
		Fn: OperatorFnGteq,
	}, ex)
}

func TestParseComplexNestedFunction(t *testing.T) {
	ex, err := Parse("funky(\"y\", concat(\"2023-12-23 \", \"15:41\"), concat(\"2024-02-21 \", \"23:59\")) >= 1")

	require.NoError(t, err)
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
		Fn: OperatorFnGteq,
	}, ex)
}

func TestParseScientificNumber(t *testing.T) {
	ex, err := Parse("9e+10")

	require.NoError(t, err)
	require.NotNil(t, ex)

	require.Equal(t, FloatNode(9e10), ex)
}

func TestParseBang(t *testing.T) {
	ex, err := Parse("a!")

	require.Nil(t, err)
	require.NotNil(t, ex)

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "a"},
		},
		Op: "!",
		Fn: OperatorFnFactorial,
	}, ex)

	ex, err = Parse("!a")
	require.Error(t, err)
	require.Nil(t, ex)
}

func TestParsePowerOperator(t *testing.T) {
	ex, err := Parse("a ^ b ^ c")

	require.NoError(t, err)

	aRaisedTo := OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "b"},
			SymbolNode{Name: "c"},
		},
		Op: "^",
		Fn: OperatorFnPower,
	}

	require.Equal(t, OperatorNode{
		Args: []MathNode{
			SymbolNode{Name: "a"},
			aRaisedTo,
		},
		Op: "^",
		Fn: OperatorFnPower,
	}, ex)
}
