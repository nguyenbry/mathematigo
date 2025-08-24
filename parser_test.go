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
