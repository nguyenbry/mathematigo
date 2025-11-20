package mathematigo

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
