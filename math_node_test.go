package mathematigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFunctionNodeBuilder(t *testing.T) {
	b := newFunctionNodeBuilder().withFn("max").build()

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "max",
		},
		Args: nil,
	}, b)

	b = newFunctionNodeBuilder().withFn("").withArg(nil).build()

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "",
		},
		Args: []MathNode{nil},
	}, b)

	b = newFunctionNodeBuilder().withFn("min").withArg(newFunctionNodeBuilder().withFn("innerFn").build()).build()

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "min",
		},
		Args: []MathNode{FunctionNode{
			Fn: SymbolNode{
				Name: "innerFn",
			},
			Args: nil,
		}},
	}, b)

	b = newFunctionNodeBuilder().withFn("min").withArg(newFunctionNodeBuilder().withFn("innerFn").withArg(nil).build()).build()

	assert.Equal(t, FunctionNode{
		Fn: SymbolNode{
			Name: "min",
		},
		Args: []MathNode{FunctionNode{
			Fn: SymbolNode{
				Name: "innerFn",
			},
			Args: []MathNode{nil},
		}},
	}, b)
}

func TestConstantNodeString(t *testing.T) {
	assert.Equal(t, "\"PI\"", ConstantNode("PI").String())
}

func TestOperatorNodeString(t *testing.T) {
	op := OperatorNode{
		Op: "+",
		Args: []MathNode{
			FloatNode(1.0),
			FloatNode(2.0),
		},
	}

	assert.Equal(t, "1 + 2", op.String())

	op = OperatorNode{
		Op: "-",
		Args: []MathNode{
			FloatNode(3.0),
		},
	}

	assert.Equal(t, "-3", op.String())

	op = OperatorNode{
		Op: "!",
		Args: []MathNode{
			FloatNode(5.0),
		},
	}

	assert.Equal(t, "5!", op.String())

	op = OperatorNode{
		Op: "-",
		Args: []MathNode{
			OperatorNode{
				Op: "+",
				Args: []MathNode{
					FloatNode(1.0),
					FloatNode(2.0),
				},
			},
			FloatNode(4.0),
		},
	}

	assert.Equal(t, "1 + 2 - 4", op.String())

	op = OperatorNode{
		Op: "+",
		Args: []MathNode{
			OperatorNode{
				Op: "-",
				Args: []MathNode{
					FloatNode(1.0),
					FloatNode(2.0),
				},
			},
			FloatNode(4.0),
		},
	}

	assert.Equal(t, "1 - 2 + 4", op.String())
}
