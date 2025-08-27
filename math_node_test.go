package main

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
