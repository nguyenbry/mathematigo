package mathematigo

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFunctionNodeBuilder(t *testing.T) {
	b := newFunctionNodeBuilder().withFn("max").build()

	assert.Equal(t, &FunctionNode{
		Fn:   NewSymbolNode("max"),
		Args: nil,
	}, b)

	b = newFunctionNodeBuilder().withFn("").withArg(nil).build()

	assert.Equal(t, &FunctionNode{
		Fn:   NewSymbolNode(""),
		Args: nil,
	}, b)

	b = newFunctionNodeBuilder().withFn("min").withArg(newFunctionNodeBuilder().withFn("innerFn").build()).build()

	assert.Equal(t, &FunctionNode{
		Fn: NewSymbolNode("min"),
		Args: []MathNode{&FunctionNode{
			Fn:   NewSymbolNode("innerFn"),
			Args: nil,
		}},
	}, b)

	b = newFunctionNodeBuilder().withFn("min").withArg(newFunctionNodeBuilder().withFn("innerFn").withArg(nil).build()).build()

	assert.Equal(t, &FunctionNode{
		Fn: NewSymbolNode("min"),
		Args: []MathNode{&FunctionNode{
			Fn:   NewSymbolNode("innerFn"),
			Args: nil,
		}},
	}, b)

	b = newFunctionNodeBuilder().withFn("min").withArg(newFunctionNodeBuilder().withFn("innerFn").withArg(nil).build()).build()

	assert.Equal(t, &FunctionNode{
		Fn: NewSymbolNode("min"),
		Args: []MathNode{&FunctionNode{
			Fn: NewSymbolNode("innerFn"),
		}},
	}, b)
}

func TestConstantNodeString(t *testing.T) {
	assert.Equal(t, "\"PI\"", NewConstantNode("PI").String())
}

func TestOperatorNodeString(t *testing.T) {
	op := OperatorNode{
		Op: "+",
		Args: []MathNode{
			NewFloatNode(1.0),
			NewFloatNode(2.0),
		},
	}

	assert.Equal(t, "1 + 2", op.String())

	op = OperatorNode{
		Op: "-",
		Args: []MathNode{
			NewFloatNode(3.0),
		},
	}

	assert.Equal(t, "-3", op.String())

	op = OperatorNode{
		Op: "!",
		Args: []MathNode{
			NewFloatNode(5.0),
		},
	}

	assert.Equal(t, "5!", op.String())

	op = OperatorNode{
		Op: "-",
		Args: []MathNode{
			&OperatorNode{
				Op: "+",
				Args: []MathNode{
					NewFloatNode(1.0),
					NewFloatNode(2.0),
				},
			},
			NewFloatNode(4.0),
		},
	}

	assert.Equal(t, "1 + 2 - 4", op.String())

	op = OperatorNode{
		Op: "+",
		Args: []MathNode{
			&OperatorNode{
				Op: "-",
				Args: []MathNode{
					NewFloatNode(1.0),
					NewFloatNode(2.0),
				},
			},
			NewFloatNode(4.0),
		},
	}

	assert.Equal(t, "1 - 2 + 4", op.String())
}

func TestTransformStopsRecursion(t *testing.T) {

	var others []string

	oldNode := NewOperatorNode(
		"+",
		"",
		NewFloatNode(1.0),
		NewFloatNode(2.0),
	)

	newNode := oldNode.
		Transform(func(n MathNode) MathNode {
			if _, ok := n.(*OperatorNode); ok {
				return NewFloatNode(42.0)
			}
			others = append(others, n.String())
			return n
		})

	require.Len(t, others, 0)
	require.Equal(t, NewFloatNode(42.0), newNode)
}

func TestAAAA(t *testing.T) {

	var m MathNode = NewOperatorNode(
		"+",
		"",
		NewFloatNode(1.0),
		NewFloatNode(2.0),
	)

	n, ok := m.(*OperatorNode)
	require.True(t, ok)

	require.True(t, n == m)

	var b MathNode = NewOperatorNode(
		"+",
		"",
		NewFloatNode(1.0),
		NewFloatNode(2.0),
	)

	require.False(t, m == b)
}

func TestTransformStopsRecursionEvenIfReturnsSameStructuralType(t *testing.T) {

	var others []string

	oldNode := NewOperatorNode(
		"+",
		"",
		NewFloatNode(1.0),
		NewFloatNode(2.0),
	)

	newNode := oldNode.
		Transform(func(n MathNode) MathNode {
			if _, ok := n.(*OperatorNode); ok {
				return NewOperatorNode(
					"+",
					"",
					NewFloatNode(1.0),
					NewFloatNode(2.0),
				)
			}

			// this should not be called
			others = append(others, n.String())
			return n
		})

	require.Len(t, others, 0)
	newNodeConcrete, ok := newNode.(*OperatorNode)
	require.True(t, ok)
	require.True(t, newNodeConcrete != oldNode)

	require.Equal(t, oldNode, newNode)
}

func TestTransformContinuesRecursion(t *testing.T) {

	var others []string
	old := NewOperatorNode(
		"+",
		"",
		NewFloatNode(1.0),
		NewFloatNode(2.0),
	)
	newNode := old.
		Transform(func(n MathNode) MathNode {
			others = append(others, n.String())
			return n
		})

	require.Equal(t, []string{"1 + 2", "1", "2"}, others)
	require.Equal(t, old, newNode)
}

func TestTransformStopsRecursionFunctionNode(t *testing.T) {

	var others []string
	newNode := newFunctionNodeBuilder().
		withFn("sum").
		withArg(NewFloatNode(1.0)).
		withArg(NewFloatNode(2.0)).
		build().
		Transform(func(n MathNode) MathNode {
			if _, ok := n.(*FunctionNode); ok {
				return NewFloatNode(42.0)
			}
			others = append(others, n.String())
			return n
		})

	require.Len(t, others, 0)
	require.Equal(t, NewFloatNode(42.0), newNode)
}

func TestComplexStopsRecursion(t *testing.T) {
	exp, err := Parse(`sum(1 + 2, 3 + 4, fn((1), "a", "b") * (5 + 6) )`)

	require.NoError(t, err)

	var other []string
	exp.Transform(func(n MathNode) MathNode {
		if f, ok := n.(*FunctionNode); ok && f.Fn.Name == "fn" {
			return NewFloatNode(42.0)
		} else if cn, ok := n.(*ConstantNode); ok {
			other = append(other, cn.String())
		}
		return n
	})

	require.Len(t, other, 0)

	now, err := Parse(`sum(1 + 2, 3 + 4, 42 * (5 + 6))`)

	require.NoError(t, err)
	require.Equal(t, now, exp)
}
