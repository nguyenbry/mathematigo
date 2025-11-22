package mathematigo

type MathNode interface {
	String() string
	ForEach(func(MathNode))
	Equal(other MathNode) bool
	Transform(func(MathNode) MathNode) MathNode
}
