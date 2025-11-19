package mathematigo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type operatorFnName string

const (
	OperatorFnBitOr      operatorFnName = "bitOr"
	OperatorFnBitAnd     operatorFnName = "bitAnd"
	OperatorFnAdd        operatorFnName = "add"
	OperatorFnMinus      operatorFnName = "minus"
	OperatorFnMultiply   operatorFnName = "multiply"
	OperatorFnDivide     operatorFnName = "divide"
	OperatorFnUnequal    operatorFnName = "unequal"
	OperatorFnEqual      operatorFnName = "equal"
	OperatorFnGt         operatorFnName = "larger"
	OperatorFnGteq       operatorFnName = "largerEq"
	OperatorFnLt         operatorFnName = "smaller"
	OperatorFnLteq       operatorFnName = "smallerEq"
	OperatorFnFactorial  operatorFnName = "factorial"
	OperatorFnUnaryMinus operatorFnName = "unaryMinus"
	OperatorFnMod        operatorFnName = "mod"
	OperatorFnPower      operatorFnName = "pow"
)

type MathNode interface {
	String() string
	ForEach(func(MathNode))
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

func (s SymbolNode) ForEach(cb func(MathNode)) {
	cb(s)
}

var _ MathNode = SymbolNode{}

type FunctionNode struct {
	Fn   SymbolNode
	Args []MathNode
}

func (f FunctionNode) ForEach(cb func(MathNode)) {
	cb(f)

	for _, arg := range f.Args {
		arg.ForEach(cb) // recursively traverse children
	}
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

func (p ParenthesisNode) ForEach(cb func(MathNode)) {
	cb(p)
	p.Content.ForEach(cb) // recursively traverse content
}

var _ MathNode = ParenthesisNode{}

type OperatorNode struct {
	Args []MathNode
	Op   string
	Fn   operatorFnName
}

func (o OperatorNode) String() string {
	switch len(o.Args) {
	case 1:
		switch o.Op {
		case "!":
			return fmt.Sprintf("%s%s", o.Args[0].String(), o.Op)
		default:
			return fmt.Sprintf("%s%s", o.Op, o.Args[0].String())
		}
	case 2:
		// binary
		return fmt.Sprintf("%s %s %s", o.Args[0].String(), o.Op, o.Args[1].String())
	}
	panic("todo String() OperatorNode")
}

func (o OperatorNode) ForEach(cb func(MathNode)) {
	cb(o)

	for _, arg := range o.Args {
		arg.ForEach(cb) // recursively traverse children
	}
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

func (b BlockNode) ForEach(cb func(MathNode)) {
	cb(b)

	for _, block := range b.Blocks {
		block.ForEach(cb) // recursively traverse children
	}
}

var _ MathNode = BlockNode{}

type BooleanNode bool

func (b BooleanNode) String() string {
	return strconv.FormatBool(bool(b))
}

func (b BooleanNode) ForEach(cb func(MathNode)) {
	cb(b)
}

var _ MathNode = BooleanNode(true)

type NullNode struct{}

func (n NullNode) String() string {
	return "null"
}

func (n NullNode) ForEach(cb func(MathNode)) {
	cb(n)
}

var _ MathNode = NullNode{}

type FloatNode float64

func (f FloatNode) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f FloatNode) ForEach(cb func(MathNode)) {
	cb(f)
}

// IsInt checks if the FloatNode represents an integer value
func (f FloatNode) IsInt() bool {
	val := float64(f)
	return val == math.Trunc(val)
}

// AsInt converts the FloatNode to an int64 if it represents an integer
// Returns the integer value and true if successful, 0 and false otherwise
func (f FloatNode) AsInt() (int64, bool) {
	val := float64(f)
	if val == math.Trunc(val) {
		return int64(val), true
	}
	return 0, false
}

// ToIntNode converts the FloatNode to an IntNode if it represents an integer
// Returns the IntNode and true if successful, FloatNode and false otherwise
func (f FloatNode) ToIntNode() (IntNode, bool) {
	if intVal, ok := f.AsInt(); ok {
		return IntNode(intVal), true
	}
	return 0, false
}

var _ MathNode = FloatNode(0)

type IntNode int64

func (i IntNode) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i IntNode) ForEach(cb func(MathNode)) {
	cb(i)
}

var _ MathNode = IntNode(0)

type functionNodeBuilder struct {
	fNode FunctionNode
}

type ConstantNode string

func (c ConstantNode) String() string {
	return `"` + string(c) + `"`
}

func (c ConstantNode) ForEach(cb func(MathNode)) {
	cb(c)
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
