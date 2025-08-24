package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimaryFalse(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	if v, ok := ex.(Literal); ok {
		assert.Equal(t, []rune("false"), v.literal)
	} else {
		t.Fatal()
	}
}

func TestPrimaryConsumes(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	assert.Equal(t, 1, p.current)
}

func asLiteral(t *testing.T, ex Expr, f func(Literal)) {
	if v, ok := ex.(Literal); ok {
		f(v)
	} else {
		t.Fatal("not Literal")
	}
}

func TestMultiplePrimaryCalls(t *testing.T) {
	s := NewScanner(" false true null\n")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	asLiteral(t, ex, func(v Literal) {
		assert.Equal(t, []rune("false"), v.literal)
	})

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	asLiteral(t, ex, func(v Literal) {
		assert.Equal(t, []rune("true"), v.literal)
	})

	ex, err = p.primary()
	assert.Nil(t, err)
	assert.NotNil(t, ex)
	asLiteral(t, ex, func(v Literal) {
		assert.Equal(t, []rune("null"), v.literal)
	})

}

func TestExpression(t *testing.T) {
	s := NewScanner(" false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.NotNil(t, ex)

	if v, ok := ex.(Literal); ok {
		assert.Equal(t, []rune("false"), v.literal)
	} else {
		t.Fatal()
	}
}

func TestGrouping(t *testing.T) {
	s := NewScanner("(false)")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, err)
	assert.Equal(t, Grouping{
		content: Literal{
			literal: []rune("false"),
		},
	}, ex)
}

func TestGroupingErrors(t *testing.T) {
	s := NewScanner("(false")

	toks := s.scanTokens()

	p := NewParser(toks)

	ex, err := p.expression()

	assert.Nil(t, ex)
	assert.NotNil(t, err)
}
