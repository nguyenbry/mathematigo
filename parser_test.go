package mathematigo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrimaryFalse(t *testing.T) {
	s := NewScanner(" false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, NewBooleanNode(false), ex)
}

func TestPrimaryConsumes(t *testing.T) {
	s := NewScanner(" false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, 1, p.current)
}

func TestMultiplePrimaryCalls(t *testing.T) {
	s := NewScanner(" false true null\n")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	require.Equal(t, NewBooleanNode(false), ex)

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	require.Equal(t, NewBooleanNode(true), ex)

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	require.Equal(t, &NullNode{}, ex)
}

func TestExpression(t *testing.T) {
	s := NewScanner(" false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, NewBooleanNode(false), ex)
}

func TestGrouping(t *testing.T) {
	s := NewScanner("(false)")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.Equal(t, &ParenthesisNode{Content: NewBooleanNode(false)}, ex)
}

func TestGroupingErrors(t *testing.T) {
	s := NewScanner("(false")

	toks, err := s.scanTokens()
	require.NoError(t, err)

	p := newParser(toks)

	ex, err := p.expression()

	assert.Nil(t, ex)
	assert.NotNil(t, err)
}

// func TestParseEmpty(t *testing.T) {
// 	s := NewScanner(" ")

// 	toks := s.scanTokens()

// 	p := NewParser(toks)

// 	ex, err := p.expression()

// 	assert.Nil(t, ex)
// 	assert.NotNil(t, err)
// }

func TestParseNoArgsFunction(t *testing.T) {
	ex, err := Parse("myFunc()")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &FunctionNode{Fn: NewSymbolNode("myFunc"), Args: nil}, ex)
}

func TestParseFunction1Arg(t *testing.T) {
	ex, err := Parse("myFunc(2)")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &FunctionNode{Fn: NewSymbolNode("myFunc"), Args: []MathNode{NewFloatNode(2)}}, ex)
}

func TestImplicitMult(t *testing.T) {
	ex, err := Parse("2 a")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewFloatNode(2), NewSymbolNode("a")}, Op: "*", Fn: OperatorFnMultiply}, ex)
}

func TestImplicitMult2(t *testing.T) {
	ex, err := Parse("1a 2")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&OperatorNode{Args: []MathNode{NewFloatNode(1), NewSymbolNode("a")}, Op: "*", Fn: OperatorFnMultiply}, NewFloatNode(2)}, Op: "*", Fn: OperatorFnMultiply}, ex)
}

func TestBlockSimple(t *testing.T) {
	ex, err := Parse("2 \n a")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &BlockNode{Blocks: []MathNode{NewFloatNode(2), NewSymbolNode("a")}}, ex)
}

func TestTrailingNewLinesDoesNotProduceBlock(t *testing.T) {
	ex, err := Parse("2 a \n\n\n")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewFloatNode(2), NewSymbolNode("a")}, Op: "*", Fn: OperatorFnMultiply}, ex)
}

func TestLeadingNewLinesProducesBlock(t *testing.T) {
	ex, err := Parse("\n2 a")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &BlockNode{Blocks: []MathNode{&OperatorNode{Args: []MathNode{NewFloatNode(2), NewSymbolNode("a")}, Op: "*", Fn: OperatorFnMultiply}}}, ex)
}

func TestMultipleBlocksWithFunctionCall(t *testing.T) {
	ex, err := Parse("\n2 a\nmyFunc(2)")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &BlockNode{Blocks: []MathNode{&OperatorNode{Args: []MathNode{NewFloatNode(2), NewSymbolNode("a")}, Op: "*", Fn: OperatorFnMultiply}, &FunctionNode{Fn: NewSymbolNode("myFunc"), Args: []MathNode{NewFloatNode(2)}}}}, ex)
}

func TestMultipleBlocksWithFunctionCallAndAddition(t *testing.T) {
	ex, err := Parse("\n2 a\nmyFunc(2) * 2")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &BlockNode{Blocks: []MathNode{&OperatorNode{Args: []MathNode{NewFloatNode(2), NewSymbolNode("a")}, Op: "*", Fn: OperatorFnMultiply}, &OperatorNode{Args: []MathNode{&FunctionNode{Fn: NewSymbolNode("myFunc"), Args: []MathNode{NewFloatNode(2)}}, NewFloatNode(2)}, Op: "*", Fn: OperatorFnMultiply}}}, ex)
}

func TestNewLineInFunctionArgs(t *testing.T) {
	ex, err := Parse("myFunc(2 \n 3)")

	assert.NotNil(t, err)
	assert.Nil(t, ex)
}

func TestFactorial(t *testing.T) {
	ex, err := Parse("a!")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewSymbolNode("a")}, Op: "!", Fn: OperatorFnFactorial}, ex)
}

func TestFactorialAndUnaryMinusPrecedence(t *testing.T) {
	ex, err := Parse("-a!")

	require.NoError(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&OperatorNode{Args: []MathNode{NewSymbolNode("a")}, Op: "!", Fn: OperatorFnFactorial}}, Op: "-", Fn: OperatorFnUnaryMinus}, ex)
}

func TestParseFunctionMultipleArgs(t *testing.T) {
	ex, err := Parse("myFunc(2, 3, x)")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &FunctionNode{Fn: NewSymbolNode("myFunc"), Args: []MathNode{NewFloatNode(2), NewFloatNode(3), NewSymbolNode("x")}}, ex)
}

func TestParseAmpersandBindsTighterThanPipe(t *testing.T) {
	ex, err := Parse("a | b & c")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewSymbolNode("a"), &OperatorNode{Args: []MathNode{NewSymbolNode("b"), NewSymbolNode("c")}, Op: "&", Fn: OperatorFnBitAnd}}, Op: "|", Fn: OperatorFnBitOr}, ex)

	ex, err = Parse("a & b | c")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&OperatorNode{Args: []MathNode{NewSymbolNode("a"), NewSymbolNode("b")}, Op: "&", Fn: OperatorFnBitAnd}, NewSymbolNode("c")}, Op: "|", Fn: OperatorFnBitOr}, ex)
}

func TestParsePipe(t *testing.T) {
	ex, err := Parse("(\"X\" | \"y\")")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &ParenthesisNode{Content: &OperatorNode{Args: []MathNode{NewConstantNode("X"), NewConstantNode("y")}, Op: "|", Fn: OperatorFnBitOr}}, ex)
}

func TestParseDoubleEqualsTighterThanAmpersand(t *testing.T) {
	ex, err := Parse("a == b & c")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&OperatorNode{Args: []MathNode{NewSymbolNode("a"), NewSymbolNode("b")}, Op: "==", Fn: OperatorFnEqual}, NewSymbolNode("c")}, Op: "&", Fn: OperatorFnBitAnd}, ex)

	ex, err = Parse("a & b == c")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewSymbolNode("a"), &OperatorNode{Args: []MathNode{NewSymbolNode("b"), NewSymbolNode("c")}, Op: "==", Fn: OperatorFnEqual}}, Op: "&", Fn: OperatorFnBitAnd}, ex)
}

func TestImplicitMultiplicationCannotHappenForConstantNode(t *testing.T) {
	ex, err := Parse(`x "abc"`)

	if !assert.Error(t, err) {
		fmt.Printf("error was nil, ex was: %+v\n", ex)
		fmt.Printf("ex type: %T\n", ex)
		t.FailNow()
	}
	require.Error(t, err)
	require.Nil(t, ex)

	ex, err = Parse(`x"abc"`)

	require.Error(t, err)
	require.Nil(t, ex)

	ex, err = Parse(`x "abc" * 2`)

	require.Error(t, err)
	require.Nil(t, ex)

	ex, err = Parse(`x ("abc")`)

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &FunctionNode{Fn: NewSymbolNode("x"), Args: []MathNode{NewConstantNode("abc")}}, ex)
}

func TestParseImplicitMult(t *testing.T) {
	ex, err := Parse("(1+2)(3+4)")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&ParenthesisNode{Content: &OperatorNode{Args: []MathNode{NewFloatNode(1), NewFloatNode(2)}, Op: "+", Fn: OperatorFnAdd}}, &ParenthesisNode{Content: &OperatorNode{Args: []MathNode{NewFloatNode(3), NewFloatNode(4)}, Op: "+", Fn: OperatorFnAdd}}}, Op: "*", Fn: OperatorFnMultiply}, ex)
}

func TestParseComplexFunction(t *testing.T) {
	ex, err := Parse("pattern_match(\"x\", \"2023-12-23 15:41\", \"2024-02-21 23:59\") >= 1")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&FunctionNode{Fn: NewSymbolNode("pattern_match"), Args: []MathNode{NewConstantNode("x"), NewConstantNode("2023-12-23 15:41"), NewConstantNode("2024-02-21 23:59")}}, NewFloatNode(1)}, Op: ">=", Fn: OperatorFnGteq}, ex)
}

func TestParseComplexNestedFunction(t *testing.T) {
	ex, err := Parse("funky(\"y\", concat(\"2023-12-23 \", \"15:41\"), concat(\"2024-02-21 \", \"23:59\")) >= 1")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{&FunctionNode{Fn: NewSymbolNode("funky"), Args: []MathNode{NewConstantNode("y"), &FunctionNode{Fn: NewSymbolNode("concat"), Args: []MathNode{NewConstantNode("2023-12-23 "), NewConstantNode("15:41")}}, &FunctionNode{Fn: NewSymbolNode("concat"), Args: []MathNode{NewConstantNode("2024-02-21 "), NewConstantNode("23:59")}}}}, NewFloatNode(1)}, Op: ">=", Fn: OperatorFnGteq}, ex)
}

func TestParseScientificNumber(t *testing.T) {
	ex, err := Parse("9e+10")

	require.NoError(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, NewFloatNode(9e10), ex)
}

func TestParseBang(t *testing.T) {
	ex, err := Parse("a!")

	require.Nil(t, err)
	require.NotNil(t, ex)

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewSymbolNode("a")}, Op: "!", Fn: OperatorFnFactorial}, ex)

	ex, err = Parse("!a")
	require.Error(t, err)
	require.Nil(t, ex)
}

func TestParsePowerOperator(t *testing.T) {
	ex, err := Parse("a ^ b ^ c")

	require.NoError(t, err)

	aRaisedTo := &OperatorNode{Args: []MathNode{NewSymbolNode("b"), NewSymbolNode("c")}, Op: "^", Fn: OperatorFnPower}

	assert.Equal(t, &OperatorNode{Args: []MathNode{NewSymbolNode("a"), aRaisedTo}, Op: "^", Fn: OperatorFnPower}, ex)
}
