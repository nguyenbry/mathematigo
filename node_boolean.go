package mathematigo

import "strconv"

type BooleanNode bool

func NewBooleanNode(v bool) *BooleanNode {
	b := BooleanNode(v)
	return &b
}

func (b *BooleanNode) String() string {
	return strconv.FormatBool(bool(*b))
}

func (b *BooleanNode) ForEach(cb func(MathNode)) {
	cb(b)
}

func (b *BooleanNode) Equal(other MathNode) bool {
	otherBool, ok := other.(*BooleanNode)
	return ok && *b == *otherBool
}

var _ MathNode = (*BooleanNode)(nil)
