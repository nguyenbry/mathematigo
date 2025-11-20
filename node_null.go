package mathematigo

func NewNullNode() *NullNode { return &NullNode{} }

type NullNode struct{}

func (n *NullNode) String() string {
	return "null"
}

func (n *NullNode) ForEach(cb func(MathNode)) {
	cb(n)
}

func (n *NullNode) Equal(other MathNode) bool {
	_, ok := other.(*NullNode)
	return ok
}

var _ MathNode = (*NullNode)(nil)
