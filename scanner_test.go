package main

import (
	"fmt"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
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
		fmt.Println(unicode.IsGraphic(x))

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
	}, s.scanTokens())
}

func TestScanTokensTwiceDoesNothing(t *testing.T) {
	s := NewScanner(")+")

	resA := s.scanTokens()
	resB := s.scanTokens()

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
	}, s.scanTokens())
}

func TestLongerScanTokens(t *testing.T) {
	s := NewScanner("!!=)()<<=<>")

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
	}, s.scanTokens())
}

func TestWhitespaceLineNumDoesntIncrement(t *testing.T) {
	s := NewScanner(".  !\t")

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
	}, s.scanTokens())

	assert.Equal(t, 0, s.line)
}

func TestWhitespaceLineNumIncrements(t *testing.T) {
	s := NewScanner(".  !\n")

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
	}, s.scanTokens())

	assert.Equal(t, 1, s.line)
}

func TestTokenAfterWhitespaceHasCorrectLineNum(t *testing.T) {
	s := NewScanner(".  !\n.  <")

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
	}, s.scanTokens())

	assert.Equal(t, 1, s.line)
}

func TestStringToken(t *testing.T) {
	s := NewScanner("!\"hey\"")

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
	}, s.scanTokens())
}

func TestIsASCIIDigit(t *testing.T) {
	rs := []rune("0123456789")

	for _, x := range rs {
		assert.True(t, isASCIIDigit(x))
	}
}

func TestNumberStartsWithDot(t *testing.T) {
	s := NewScanner(" .1234 ")

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune(".1234"),
			Line: 0,
		},
	}, s.scanTokens())
}

func TestNumberBangNumber(t *testing.T) {
	s := NewScanner("1!1")

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
	}, s.scanTokens())
}

func TestScientificNumberLookalike(t *testing.T) {
	s := NewScanner("9e ")

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
	}, s.scanTokens())
}

func TestScientifics(t *testing.T) {
	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e1"),
			Line: 0,
		},
	}, NewScanner("9e1 ").scanTokens())

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e+1"),
			Line: 0,
		},
	}, NewScanner("9e+1 ").scanTokens())

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e-1"),
			Line: 0,
		},
	}, NewScanner("9e-1 ").scanTokens())

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("9e-02"),
			Line: 0,
		},
	}, NewScanner("9e-02 ").scanTokens())

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
	}, NewScanner("+ 9e-02+ ").scanTokens())
}

func TestScientificNumberWithDotBeforeE(t *testing.T) {
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
	}, NewScanner("+ 92.e-02+ ").scanTokens())
}

func TestDecimalScientific(t *testing.T) {
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
	}, NewScanner("+ .92e-02 + ").scanTokens())

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

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("73824"),
			Line: 0,
		},
	}, s.scanTokens())
}

func TestSpaceAfterDotDoesNotResultInNumberToken(t *testing.T) {
	s := NewScanner(" . 1234 ")

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
	}, s.scanTokens())
}

func TestDecimal(t *testing.T) {
	s := NewScanner("73824.2")

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("73824.2"),
			Line: 0,
		},
	}, s.scanTokens())
}

func TestDecimalWithTwoDots(t *testing.T) {
	s := NewScanner("73824..2")

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
	}, s.scanTokens())
}

func TestIntegerWithTwoDots(t *testing.T) {
	s := NewScanner("1234..")

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
	}, s.scanTokens())
}

func TestIntegerWithDotsAndWhitespace(t *testing.T) {
	s := NewScanner("1234.. .")

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
	}, s.scanTokens())
}

func TestDecimalWithOtherSymbols(t *testing.T) {
	s := NewScanner("728.3 < 391 + \n 'hi'")

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
	}, s.scanTokens())
}

func TestStringPreservesSingleOrDoubleQuotes(t *testing.T) {
	s := NewScanner("'1'")

	assert.Equal(t, []Token{
		{
			Type:    String,
			Text:    []rune("'1'"),
			Line:    0,
			Literal: []rune("1"),
		},
	}, s.scanTokens())

	s = NewScanner("\"2\"")

	assert.Equal(t, []Token{
		{
			Type:    String,
			Text:    []rune("\"2\""),
			Line:    0,
			Literal: []rune("2"),
		},
	}, s.scanTokens())
}

func TestIdent(t *testing.T) {
	s := NewScanner("'1' + 1abc \n")

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
	}, s.scanTokens())
}

func TestIdentWithUnderscoreSimple(t *testing.T) {
	s := NewScanner("_")

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_"),
		},
	}, s.scanTokens())
}

func TestIdentWithUnderscoreSimple2(t *testing.T) {
	s := NewScanner("_928u")

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_928u"),
		},
	}, s.scanTokens())
}

func TestIdentWithUnderscoreAndWhitespace(t *testing.T) {
	s := NewScanner("_928u\t\r")

	assert.Equal(t, []Token{
		{
			Type: Ident,
			Text: []rune("_928u"),
		},
	}, s.scanTokens())
}

func TestMultipleNewLinesHaveCorrectLineNum(t *testing.T) {
	s := NewScanner("_928u\n\n\t\r\n. ")

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
	}, s.scanTokens())
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

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("0b0"),
		},
	}, s.scanTokens())
}

func TestSingleScanTokenBinaryNumberIncludesSpace(t *testing.T) {
	s := NewScanner("0b0.")

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("0b0."),
		},
	}, s.scanTokens())
}

func TestSingleScanTokenBinaryNumberLookalike(t *testing.T) {
	s := NewScanner("0b2")

	assert.NotNil(t, s.scanToken())
}

func TestIntegerBeforeBinary(t *testing.T) {
	s := NewScanner("20b1") // looks like 2(0b1)

	assert.Equal(t, []Token{
		{
			Type: Number,
			Text: []rune("20"),
		},
		{
			Type: Ident,
			Text: []rune("b1"),
		},
	}, s.scanTokens())

}
