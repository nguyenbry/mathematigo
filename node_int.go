package mathematigo

import "strconv"

func NewIntNode(v int64) *IntNode { i := IntNode(v); return &i }

type IntNode int64

func (i *IntNode) String() string {
	return strconv.FormatInt(int64(*i), 10)
}

func (i *IntNode) ForEach(cb func(MathNode)) {
	cb(i)
}

func (i *IntNode) Equal(other MathNode) bool {
	otherInt, ok := other.(*IntNode)
	return ok && *i == *otherInt
}

func (i *IntNode) Transform(fn func(MathNode) MathNode) MathNode { return fn(i) }

var _ MathNode = (*IntNode)(nil)
