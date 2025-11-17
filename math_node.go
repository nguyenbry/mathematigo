package main

import (
	"fmt"
	"strconv"
	"strings"
)

type MathNode interface {
	String() string
}

type SymbolNode struct {
	Name string
}

func (s SymbolNode) Valid() bool {
	return true
}

func (s SymbolNode) String() string {
	return s.Name
}

var _ MathNode = SymbolNode{}

type FunctionNode struct {
	Fn   SymbolNode
	Args []MathNode
}

func (f FunctionNode) String() string {
	s := fmt.Sprintf("%s(", f.Fn.String())

	for i, node := range f.Args {
		if i == len(f.Args)-1 {
			s += fmt.Sprintf("%s)", node.String())
		} else {
			s += fmt.Sprintf("%s,", node.String())
		}
	}

	return s
}

var _ MathNode = FunctionNode{}

type ParenthesisNode struct {
	Content MathNode // not nil
}

func (p ParenthesisNode) String() string {
	return fmt.Sprintf("(%s)", p.Content.String())
}

var _ MathNode = ParenthesisNode{}

type OperatorNode struct {
	Args []MathNode
	Op   string
}

func (o OperatorNode) String() string {
	panic("todo String() OperatorNode")
}

var _ MathNode = OperatorNode{}

type BlockNode struct {
	Blocks []MathNode
}

func (b BlockNode) String() string {
	parts := make([]string, 0, len(b.Blocks))

	for _, x := range b.Blocks {
		parts = append(parts, x.String())
	}

	return strings.Join(parts, "\n")
}

var _ MathNode = BlockNode{}

type BooleanNode bool

func (b BooleanNode) String() string {
	return strconv.FormatBool(bool(b))
}

var _ MathNode = BooleanNode(true)

type NullNode struct{}

func (n NullNode) String() string {
	return "null"
}

var _ MathNode = NullNode{}

type FloatNode float64

func (f FloatNode) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

var _ MathNode = FloatNode(0)

type IntNode int64

func (i IntNode) String() string {
	return strconv.FormatInt(int64(i), 10)
}

var _ MathNode = IntNode(0)

type functionNodeBuilder struct {
	fNode FunctionNode
}

type ConstantNode string

func (c ConstantNode) String() string {
	return string(c)
}

var _ MathNode = ConstantNode("")

func newFunctionNodeBuilder() functionNodeBuilder {
	return functionNodeBuilder{}
}

func (b functionNodeBuilder) withArg(arg MathNode) functionNodeBuilder {
	b.fNode.Args = append(b.fNode.Args, arg)
	return b
}

func (b functionNodeBuilder) withFn(name string) functionNodeBuilder {
	b.fNode.Fn.Name = name

	return b
}

func (b functionNodeBuilder) build() FunctionNode {
	return b.fNode
}
