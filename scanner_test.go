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
