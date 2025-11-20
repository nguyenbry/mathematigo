package mathematigo

import (
	"math"
	"strconv"
)

type FloatNode float64

func (f FloatNode) String() string {
	return strconv.FormatFloat(float64(f), 'g', -1, 64)
}

func (f FloatNode) ForEach(cb func(MathNode)) {
	cb(f)
}

func (f FloatNode) Equal(other MathNode) bool {
	otherFloat, ok := other.(FloatNode)
	return ok && f == otherFloat
}

// IsInt checks if the FloatNode represents an integer value
func (f FloatNode) IsInt() bool {
	val := float64(f)
	return val == math.Trunc(val)
}

// AsInt converts the FloatNode to an int64 if it represents an integer
// Returns the integer value and true if successful, 0 and false otherwise
func (f FloatNode) AsInt() (int64, bool) {
	val := float64(f)
	if val == math.Trunc(val) {
		return int64(val), true
	}
	return 0, false
}

// ToIntNode converts the FloatNode to an IntNode if it represents an integer
// Returns the IntNode and true if successful, FloatNode and false otherwise
func (f FloatNode) ToIntNode() (IntNode, bool) {
	if intVal, ok := f.AsInt(); ok {
		return IntNode(intVal), true
	}
	return 0, false
}

var _ MathNode = FloatNode(0)
