package mathematigo

type ConstantNode string

func NewConstantNode(value string) *ConstantNode {
	c := ConstantNode(value)
	return &c
}

func (c *ConstantNode) String() string {
	return `"` + string(*c) + `"`
}

func (c *ConstantNode) ForEach(cb func(MathNode)) {
	cb(c)
}

func (c *ConstantNode) Equal(other MathNode) bool {
	otherConst, ok := other.(*ConstantNode)
	return ok && *c == *otherConst
}

func (c *ConstantNode) Transform(f func(MathNode) MathNode) MathNode { return f(c) }

var _ MathNode = (*ConstantNode)(nil)
