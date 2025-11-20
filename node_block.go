package mathematigo

import "strings"

type BlockNode struct {
	Blocks []MathNode
}

func (b BlockNode) String() string {
	parts := make([]string, 0, len(b.Blocks))

	for _, x := range b.Blocks {
		parts = append(parts, x.String())
	}

	return strings.Join(parts, "\n")
}

func (b BlockNode) ForEach(cb func(MathNode)) {
	cb(b)

	for _, block := range b.Blocks {
		block.ForEach(cb) // recursively traverse children
	}
}

func (b BlockNode) Equal(other MathNode) bool {
	otherBlock, ok := other.(BlockNode)
	if !ok {
		return false
	}

	if len(b.Blocks) != len(otherBlock.Blocks) {
		return false
	}

	for i := range b.Blocks {
		if !b.Blocks[i].Equal(otherBlock.Blocks[i]) {
			return false
		}
	}

	return true
}

var _ MathNode = BlockNode{}
