package mathematigo

import "strconv"

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
