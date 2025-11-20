package mathematigo

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

type OperatorFnName string

const (
	OperatorFnBitOr      OperatorFnName = "bitOr"
	OperatorFnBitAnd     OperatorFnName = "bitAnd"
	OperatorFnAdd        OperatorFnName = "add"
	OperatorFnMinus      OperatorFnName = "minus"
	OperatorFnMultiply   OperatorFnName = "multiply"
	OperatorFnDivide     OperatorFnName = "divide"
	OperatorFnUnequal    OperatorFnName = "unequal"
	OperatorFnEqual      OperatorFnName = "equal"
	OperatorFnGt         OperatorFnName = "larger"
	OperatorFnGteq       OperatorFnName = "largerEq"
	OperatorFnLt         OperatorFnName = "smaller"
	OperatorFnLteq       OperatorFnName = "smallerEq"
	OperatorFnFactorial  OperatorFnName = "factorial"
	OperatorFnUnaryMinus OperatorFnName = "unaryMinus"
	OperatorFnMod        OperatorFnName = "mod"
	OperatorFnPower      OperatorFnName = "pow"
)

var operatorFnsMap = map[OperatorFnName]struct{}{
	OperatorFnBitOr:      {},
	OperatorFnBitAnd:     {},
	OperatorFnAdd:        {},
	OperatorFnMinus:      {},
	OperatorFnMultiply:   {},
	OperatorFnDivide:     {},
	OperatorFnUnequal:    {},
	OperatorFnEqual:      {},
	OperatorFnGt:         {},
	OperatorFnGteq:       {},
	OperatorFnLt:         {},
	OperatorFnLteq:       {},
	OperatorFnFactorial:  {},
	OperatorFnUnaryMinus: {},
	OperatorFnMod:        {},
	OperatorFnPower:      {},
}

func (o OperatorFnName) Valid() bool {
	_, ok := operatorFnsMap[o]
	return ok
}

type MathNode interface {
	String() string
	ForEach(func(MathNode))
	Equal(other MathNode) bool
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

func (s SymbolNode) Equal(other MathNode) bool {
	otherSym, ok := other.(SymbolNode)
	return ok && s.Name == otherSym.Name
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

func (f FunctionNode) Equal(other MathNode) bool {
	otherFunc, ok := other.(FunctionNode)
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

func (p ParenthesisNode) Equal(other MathNode) bool {
	otherPar, ok := other.(ParenthesisNode)
	if !ok {
		return false
	}

	return p.Content.Equal(otherPar.Content)
}

var _ MathNode = ParenthesisNode{}

type OperatorNode struct {
	Args []MathNode
	Op   string
	Fn   OperatorFnName
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

func (o OperatorNode) Equal(other MathNode) bool {
	otherOp, ok := other.(OperatorNode)
	if !ok {
		return false
	}

	if o.Op != otherOp.Op || o.Fn != otherOp.Fn || len(o.Args) != len(otherOp.Args) {
		return false
	}

	for i := range o.Args {
		if !o.Args[i].Equal(otherOp.Args[i]) {
			return false
		}
	}

	return true
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

func (b BlockNode) Equal(other MathNode) bool {
	otherBlock, ok := other.(BlockNode)
	if !ok {
		return false
	}

	if len(b.Blocks) != len(otherBlock.Blocks) {
		return false
	}

	for i := range b.Blocks {
		if !b.Blocks[i].Equal(otherBlock.Blocks[i]) {
			return false
		}
	}

	return true
}

var _ MathNode = BlockNode{}

type BooleanNode bool

func (b BooleanNode) String() string {
	return strconv.FormatBool(bool(b))
}

func (b BooleanNode) ForEach(cb func(MathNode)) {
	cb(b)
}

func (b BooleanNode) Equal(other MathNode) bool {
	otherBool, ok := other.(BooleanNode)
	return ok && b == otherBool
}

var _ MathNode = BooleanNode(true)

type NullNode struct{}

func (n NullNode) String() string {
	return "null"
}

func (n NullNode) ForEach(cb func(MathNode)) {
	cb(n)
}

func (n NullNode) Equal(other MathNode) bool {
	_, ok := other.(NullNode)
	return ok
}

var _ MathNode = NullNode{}

type FloatNode float64

func (f FloatNode) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f FloatNode) ForEach(cb func(MathNode)) {
	cb(f)
}

func (f FloatNode) Equal(other MathNode) bool {
	otherFloat, ok := other.(FloatNode)
	return ok && f == otherFloat
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

func (i IntNode) Equal(other MathNode) bool {
	otherInt, ok := other.(IntNode)
	return ok && i == otherInt
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

func (c ConstantNode) Equal(other MathNode) bool {
	otherConst, ok := other.(ConstantNode)
	return ok && c == otherConst
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
