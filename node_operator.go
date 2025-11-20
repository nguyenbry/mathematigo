package mathematigo

import "fmt"

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
