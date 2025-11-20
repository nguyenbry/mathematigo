package mathematigo

type SymbolNode struct {
	Name string
}

func NewSymbolNode(name string) *SymbolNode {
	return &SymbolNode{Name: name}
}

func (s *SymbolNode) Valid() bool {
	return true
}

func (s *SymbolNode) String() string {
	return s.Name
}

func (s *SymbolNode) ForEach(cb func(MathNode)) {
	cb(s)
}

func (s *SymbolNode) Equal(other MathNode) bool {
	otherSym, ok := other.(*SymbolNode)
	return ok && s.Name == otherSym.Name
}

var _ MathNode = (*SymbolNode)(nil)
