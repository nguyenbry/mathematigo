package mathematigo

import "fmt"

type FunctionNode struct {
	Fn   *SymbolNode
	Args []MathNode
}

func (f *FunctionNode) ForEach(cb func(MathNode)) {
	cb(f)
	cb(f.Fn)

	for _, arg := range f.Args {
		arg.ForEach(cb) // recursively traverse children
	}
}

func (f *FunctionNode) String() string {
	s := fmt.Sprintf("%s(", f.Fn.String())

	for i, node := range f.Args {
		if i == len(f.Args)-1 {
			s += fmt.Sprintf("%s)", node.String())
		} else {
			s += fmt.Sprintf("%s, ", node.String())
		}
	}

	return s
}

func (f *FunctionNode) Equal(other MathNode) bool {
	otherFunc, ok := other.(*FunctionNode)
	if !ok {
		return false
	}

	if !f.Fn.Equal(otherFunc.Fn) {
		return false
	}

	if len(f.Args) != len(otherFunc.Args) {
		return false
	}

	for i := range f.Args {
		if !f.Args[i].Equal(otherFunc.Args[i]) {
			return false
		}
	}

	return true
}

var _ MathNode = (*FunctionNode)(nil)

type functionNodeBuilder struct {
	fNode *FunctionNode
}

func newFunctionNodeBuilder() functionNodeBuilder {
	return functionNodeBuilder{fNode: &FunctionNode{}}
}

func (b functionNodeBuilder) withArg(arg MathNode) functionNodeBuilder {
	if arg != nil {
		b.fNode.Args = append(b.fNode.Args, arg)
	}
	return b
}

func (b functionNodeBuilder) withFn(name string) functionNodeBuilder {
	b.fNode.Fn = NewSymbolNode(name)
	return b
}

func (b functionNodeBuilder) build() *FunctionNode {
	return b.fNode
}

func NewFunctionNode(name string, args ...MathNode) *FunctionNode {
	return &FunctionNode{Fn: NewSymbolNode(name), Args: args}
}
