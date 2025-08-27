package main

type MathNode interface{}

type SymbolNode struct {
	Name string
}

func (f SymbolNode) Valid() bool {
	return true
}

var _ MathNode = SymbolNode{}

type FunctionNode struct {
	Fn   SymbolNode
	Args []MathNode
}

// func (f FunctionNode) Valid() bool {
// 	for _, argNode := range f.Args {
// 		if !argNode.Valid() {
// 			return false
// 		}
// 	}

// 	return f.Fn.Valid()
// }

var _ MathNode = FunctionNode{}

type functionNodeBuilder struct {
	fNode FunctionNode
}

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
