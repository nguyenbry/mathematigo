package mathematigo

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScannerSource(t *testing.T) {
	s := NewScanner("a")

	assert.Equal(t, []rune("a"), s.source)
	assert.Equal(t, 1, len(s.source))
}

func TestAssumptionAboutEndLine(t *testing.T) {
	s := "a\nb"

	assert.Equal(t, 3, len(s))
}

func TestAssumptionAboutEmoji(t *testing.T) {
	s := ")+-*&^%$#"
	r := []rune(s)

	for _, x := range r {
		require.True(t, unicode.IsGraphic(x))
	}
}

func TestSingleScanToken(t *testing.T) {
	s := NewScanner(")")

	s.scanToken()

	assert.Equal(t, []Token{{
		Type: CloseParen,
		Text: []rune(")"),
		Line: 0,
	}}, s.tokens)
}

func TestScanAllTokensSimple(t *testing.T) {
	s := NewScanner(")+")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: CloseParen,
			Text: []rune(")"),
			Line: 0,
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
	}, tokens)
}

func TestScanTokensTwiceDoesNothing(t *testing.T) {
	s := NewScanner(")+")

	resA, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, resA)

	resB, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, resB)

	assert.Equal(t, []Token{
		{
			Type: CloseParen,
			Text: []rune(")"),
			Line: 0,
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
	}, resA)

	assert.Equal(t, resA, resB)
}

func TestScanBangEq(t *testing.T) {
	s := NewScanner("!!=")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
		{
			Type: BangEq,
			Text: []rune("!="),
			Line: 0,
		},
	}, tokens)
}

func TestLongerScanTokens(t *testing.T) {
	s := NewScanner("!!=)()<<=<>")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
		{
			Type: BangEq,
			Text: []rune("!="),
			Line: 0,
		},
		{
			Type: CloseParen,
			Text: []rune(")"),
			Line: 0,
		},
		{
			Type: OpenParen,
			Text: []rune("("),
			Line: 0,
		},
		{
			Type: CloseParen,
			Text: []rune(")"),
			Line: 0,
		},
		{
			Type: Lt,
			Text: []rune("<"),
			Line: 0,
		},
		{
			Type: Lteq,
			Text: []rune("<="),
			Line: 0,
		},
		{
			Type: Lt,
			Text: []rune("<"),
			Line: 0,
		},
		{
			Type: Gt,
			Text: []rune(">"),
			Line: 0,
		},
	}, tokens)
}

func TestWhitespaceLineNumDoesntIncrement(t *testing.T) {
	s := NewScanner(".  !\t")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
	}, tokens)

	assert.Equal(t, 0, s.line)
}

func TestWhitespaceLineNumIncrements(t *testing.T) {
	s := NewScanner(".  !\n")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
		{
			Type: NewLine,
		},
	}, tokens)

	assert.Equal(t, 1, s.line)
}

func TestTokenAfterWhitespaceHasCorrectLineNum(t *testing.T) {
	s := NewScanner(".  !\n.  <")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
		{
			Type: NewLine,
		},
		{
			Type: Dot,
			Text: []rune("."),
			Line: 1,
		},
		{
			Type: Lt,
			Text: []rune("<"),
			Line: 1,
		},
	}, tokens)

	assert.Equal(t, 1, s.line)
}

func TestStringToken(t *testing.T) {
	s := NewScanner("!\"hey\"")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
		{
			Type:    String,
			Text:    []rune("\"hey\""),
			Line:    0,
			Literal: []rune("hey"),
		},
	}, tokens)
}

func TestIsASCIIDigit(t *testing.T) {
	rs := []rune("0123456789")

	for _, x := range rs {
		assert.True(t, isASCIIDigit(x))
	}
}

func TestNumberStartsWithDot(t *testing.T) {
	s := NewScanner(" .1234 ")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune(".1234"),
			Line: 0,
		},
	}, tokens)
}

func TestNumberBangNumber(t *testing.T) {
	s := NewScanner("1!1")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("1"),
			Line: 0,
		},
		{
			Type: Bang,
			Text: []rune("!"),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune("1"),
			Line: 0,
		},
	}, tokens)
}

func TestScientificNumberLookalike(t *testing.T) {
	s := NewScanner("9e ")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9"),
			Line: 0,
		},
		{
			Type: Ident,
			Text: []rune("e"),
		},
	}, tokens)
}

func TestScientifics(t *testing.T) {
	tokens, err := NewScanner("9e1 ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e1"),
			Line: 0,
		},
	}, tokens)

	tokens, err = NewScanner("9e+1 ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e+1"),
			Line: 0,
		},
	}, tokens)

	tokens, err = NewScanner("9e-1 ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e-1"),
			Line: 0,
		},
	}, tokens)

	tokens, err = NewScanner("9e-02 ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e-02"),
			Line: 0,
		},
	}, tokens)

	tokens, err = NewScanner("+ 9e-02+ ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune("9e-02"),
			Line: 0,
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
	}, tokens)
}

func TestScientificNumberWithDotBeforeE(t *testing.T) {
	tokens, err := NewScanner("+ 92.e-02+ ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune("92.e-02"),
			Line: 0,
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
	}, tokens)
}

func TestDecimalScientific(t *testing.T) {
	tokens, err := NewScanner("+ .92e-02 + ").scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)
	assert.Equal(t, []Token{
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune(".92e-02"),
			Line: 0,
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
	}, tokens)

	// assert.Equal(t, []Token{
	// 	{
	// 		Type: Slash,
	// 		Text: []rune("/"),
	// 		Line: 0,
	// 	},
	// 	{
	// 		Type: Number,
	// 		Text: []rune(".92e-02"),
	// 		Line: 0,
	// 	},
	// 	{
	// 		Type: Plus,
	// 		Text: []rune("+"),
	// 		Line: 0,
	// 	},
	// }, NewScanner(" / .92e-02 + ").scanTokens())
}

func TestInteger(t *testing.T) {
	s := NewScanner("73824")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("73824"),
			Line: 0,
		},
	}, tokens)
}

func TestSpaceAfterDotDoesNotResultInNumberToken(t *testing.T) {
	s := NewScanner(" . 1234 ")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune("1234"),
			Line: 0,
		},
	}, tokens)
}

func TestDecimal(t *testing.T) {
	s := NewScanner("73824.2")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("73824.2"),
			Line: 0,
		},
	}, tokens)
}

func TestDecimalWithTwoDots(t *testing.T) {
	s := NewScanner("73824..2")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("73824."),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune(".2"),
			Line: 0,
		},
	}, tokens)
}

func TestIntegerWithTwoDots(t *testing.T) {
	s := NewScanner("1234..")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("1234."),
			Line: 0,
		},
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
	}, tokens)
}

func TestIntegerWithDotsAndWhitespace(t *testing.T) {
	s := NewScanner("1234.. .")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("1234."),
			Line: 0,
		},
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
		{
			Type: Dot,
			Text: []rune("."),
			Line: 0,
		},
	}, tokens)
}

func TestDecimalWithOtherSymbols(t *testing.T) {
	s := NewScanner("728.3 < 391 + \n 'hi'")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("728.3"),
			Line: 0,
		},
		{
			Type: Lt,
			Text: []rune("<"),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune("391"),
			Line: 0,
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
		{
			Type: NewLine,
			Line: 0,
		},
		{
			Type:    String,
			Text:    []rune("'hi'"),
			Literal: []rune("hi"),
			Line:    1,
		},
	}, tokens)
}

func TestStringPreservesSingleOrDoubleQuotes(t *testing.T) {
	s := NewScanner("'1'")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type:    String,
			Text:    []rune("'1'"),
			Line:    0,
			Literal: []rune("1"),
		},
	}, tokens)

	s = NewScanner("\"2\"")

	tokens, err = s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type:    String,
			Text:    []rune("\"2\""),
			Line:    0,
			Literal: []rune("2"),
		},
	}, tokens)
}

func TestIdent(t *testing.T) {
	s := NewScanner("'1' + 1abc \n")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type:    String,
			Text:    []rune("'1'"),
			Line:    0,
			Literal: []rune("1"),
		},
		{
			Type: Plus,
			Text: []rune("+"),
			Line: 0,
		},
		{
			Type: Number,
			Text: []rune("1"),
			Line: 0,
		},
		{
			Type: Ident,
			Text: []rune("abc"),
			Line: 0,
		},
		{
			Type: NewLine,
			Line: 0,
		},
	}, tokens)
}

func TestIdentWithUnderscoreSimple(t *testing.T) {
	s := NewScanner("_")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_"),
		},
	}, tokens)
}

func TestIdentWithUnderscoreSimple2(t *testing.T) {
	s := NewScanner("_928u")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_928u"),
		},
	}, tokens)
}

func TestIdentWithUnderscoreAndWhitespace(t *testing.T) {
	s := NewScanner("_928u\t\r")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_928u"),
		},
	}, tokens)
}

func TestMultipleNewLinesHaveCorrectLineNum(t *testing.T) {
	s := NewScanner("_928u\n\n\t\r\n. ")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_928u"),
		},
		{
			Type: NewLine,
			Line: 0,
		},
		{
			Type: NewLine,
			Line: 1,
		},
		{
			Type: NewLine,
			Line: 2,
		},
		{
			Type: Dot,
			Text: []rune("."),
			Line: 3,
		},
	}, tokens)
}

func TestSingleScanTokenBinaryNumberRequiresDigitAfterB(t *testing.T) {
	s := NewScanner("0b ")
	assert.NotNil(t, s.scanToken())
	// assert.Equal(t, []Token{}, s.scanTokens())
}

func TestSingleScanTokenBinaryNumberRequiresDigitAfterB2(t *testing.T) {
	s := NewScanner("0b.")
	assert.NotNil(t, s.scanToken())
	// assert.Equal(t, []Token{}, s.scanTokens())
}

func TestSingleScanTokenBinaryNumberSimple(t *testing.T) {
	s := NewScanner("0b0")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("0b0"),
		},
	}, tokens)
}

func TestSingleScanTokenBinaryNumberIncludesSpace(t *testing.T) {
	s := NewScanner("0b0.")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("0b0."),
		},
	}, tokens)
}

func TestSingleScanTokenBinaryNumberLookalike(t *testing.T) {
	s := NewScanner("0b2")

	assert.NotNil(t, s.scanToken())
}

func TestIntegerBeforeBinary(t *testing.T) {
	s := NewScanner("20b1") // looks like 2(0b1)

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("20"),
		},
		{
			Type: Ident,
			Text: []rune("b1"),
		},
	}, tokens)

}

func TestScanScientificNumber(t *testing.T) {
	t.Parallel()
	s := NewScanner("9e+10")

	tokens, err := s.scanTokens()
	require.NoError(t, err)
	require.NotNil(t, tokens)

	require.Equal(t, []Token{
		{Type: Number,
			Text: []rune("9e+10"),
			Line: 0,
		},
	}, tokens)
}

func TestScanUnterminatedStringErrors(t *testing.T) {
	t.Parallel()
	s := NewScanner(`"50 < 20`)

	toks, err := s.scanTokens()
	require.Error(t, err)
	require.Nil(t, toks)
}
