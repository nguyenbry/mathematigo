package mathematigo

import "fmt"

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
