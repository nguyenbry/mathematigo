package main

import (
	"errors"
	"fmt"
)

var ErrInvalidSyntax = errors.New("invalid syntax")

type Scanner struct {
	source    []rune
	sourceLen int
	tokens    []Token

	current int
	start   int
	line    int
}

func NewScanner(source string) *Scanner {
	r := []rune(source)

	return &Scanner{
		source:    []rune(source),
		sourceLen: len(r),

		tokens: nil,

		current: 0,
		start:   0,
		line:    0,
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= s.sourceLen
}

// can only call when not at end
func (s *Scanner) advance() rune {
	idx := s.current
	s.current++
	return s.source[idx]
}

func (s *Scanner) peek() (rune, bool) {
	if s.isAtEnd() {
		return ' ', false
	}

	return s.source[s.current], true
}

func (s *Scanner) peekMany(num int) ([]rune, bool) {
	if num < 0 {
		panic("TODO")
	}

	if s.current+num > len(s.source) {
		return nil, false
	}

	return s.source[s.current : s.current+num], true
}

func (s *Scanner) string(opener rune) []rune {
	for next, isNotEnd := s.peek(); next != opener && isNotEnd; next, isNotEnd = s.peek() {
		if next == opener {
			// closer found
			break
		}
		if next == '\n' {
			s.line++
		}

		s.advance()
	}

	if s.isAtEnd() {
		panic("unterminated")
	}

	// the closer is at closerIdx, so we can slice using
	// this index since the slice is exclusive end
	closerIdx := s.current

	// this increments s.current
	s.advance()

	afterOpenerIdx := s.start + 1

	return s.source[afterOpenerIdx:closerIdx]

}

func isASCIIDigit(r rune) bool {
	return 48 <= r && r <= 57
}

func isASCIIAlpha(r rune) bool {
	return (65 <= r && r <= 90) || (97 <= r && r <= 122)
}

func isASCIIAlphanumeric(r rune) bool {
	return isASCIIAlpha(r) || isASCIIDigit(r)
}

func (s *Scanner) scanToken() error {
	r := s.advance()

	switch r {
	case ')':
		s.addToken(NewToken(CloseParen, s.source[s.start:s.current], s.line, nil))
		return nil
	case '(':
		s.addToken(NewToken(OpenParen, s.source[s.start:s.current], s.line, nil))
		return nil
	case '+':
		s.addToken(NewToken(Plus, s.source[s.start:s.current], s.line, nil))
		return nil
	case '-':
		s.addToken(NewToken(Minus, s.source[s.start:s.current], s.line, nil))
		return nil
	case '.':
		s.advanceTilEndOfNumber(false)

		isDot := s.current == s.start+1 // that is, s.source[s.start:s.current] == []rune(".")

		if isDot {
			s.addToken(NewToken(Dot, s.source[s.start:s.current], s.line, nil))
		} else {
			s.scanScientific()
			s.addToken(NewToken(Number, s.source[s.start:s.current], s.line, nil))
		}
		return nil
	case '*':
		s.addToken(NewToken(Star, s.source[s.start:s.current], s.line, nil))
		return nil
	case ';':
		s.addToken(NewToken(Semi, s.source[s.start:s.current], s.line, nil))
		return nil
	case '/':
		s.addToken(NewToken(Slash, s.source[s.start:s.current], s.line, nil))
		return nil
	case ',':
		s.addToken(NewToken(Comma, s.source[s.start:s.current], s.line, nil))
		return nil
	case '<':
		if s.matchNext('=') {
			s.addToken(NewToken(Lteq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Lt, s.source[s.start:s.current], s.line, nil))
		}
		return nil
	case '>':
		if s.matchNext('=') {
			s.addToken(NewToken(Gteq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Gt, s.source[s.start:s.current], s.line, nil))
		}
		return nil
	case '=':
		if s.matchNext('=') {
			s.addToken(NewToken(EqEq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Eq, s.source[s.start:s.current], s.line, nil))
		}

		return nil
	case '!':
		if s.matchNext('=') {
			s.addToken(NewToken(BangEq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Bang, s.source[s.start:s.current], s.line, nil))
		}

		return nil
	case '\'', '"':
		val := s.string(r)
		s.addToken(NewToken(String, s.source[s.start:s.current], s.line, val))

		return nil
	case ' ', '\t', '\r':
		// do nothing
		return nil
	case '\n':
		s.addToken(NewToken(NewLine, nil, s.line, nil))
		s.line++

		return nil
	case 48:
		s.advanceTilEndOfNumber(true)

		// advanceTilEndOfNumber will go and produce up til: 0 (d)* (. (d)*)?
		// that

		justZero := s.current == s.start+1

		fmt.Println("just", justZero)

		if justZero {
			next, ok := s.peek()

			if ok && next == 'b' {
				// must be binary
				s.advance() // consume 'b'

				next, ok := s.peek()

				if !ok || (next != '1' && next != '0') {
					return ErrInvalidSyntax
				}

				s.advance() // at least one num
				s.advanceTilEndOfNumberBinary()
				s.addToken(NewToken(Number, s.source[s.start:s.current], s.line, nil))
			} else {
				// default case
				s.scanScientific()
				s.addToken(NewToken(Number, s.source[s.start:s.current], s.line, nil))

			}
		} else {
			// default case
			s.scanScientific()
			s.addToken(NewToken(Number, s.source[s.start:s.current], s.line, nil))
		}

		return nil
	case 49, 50, 51, 52, 53, 54, 55, 56, 57:
		s.advanceTilEndOfNumber(true)
		s.scanScientific()
		s.addToken(NewToken(Number, s.source[s.start:s.current], s.line, nil))

		return nil
	default:
		if isIdentifierChar(r) {
			for next, ok := s.peek(); ok && isIdentifierChar(next); next, ok = s.peek() {
				s.advance()
			}

			s.addToken(NewToken(Ident, s.source[s.start:s.current], s.line, nil))

			return nil
		} else {
			return errors.New("TODO")
		}
	}
}

func (s *Scanner) scanScientific() {
	next, ok := s.peek()

	if !ok || next != 'e' {
		return
	}

	// handle scientific
	// if +,i comes next, MUST have digits after
	// else if digits comes next, then scientific
	// else default case
	if nexts, ok := s.peekMany(2); ok {
		afterE := nexts[1]

		if isASCIIDigit(afterE) {
			// must be scientific, parse til end of digits
			s.advance() // consume 'e'

			// loop runs at least once
			for next, ok := s.peek(); ok && isASCIIDigit(next); next, ok = s.peek() {
				s.advance()
			}

		} else if afterE == '-' || afterE == '+' {
			s.advance() // consume 'e'
			s.advance() // consume +,-

			mustDigi, ok := s.peek()

			if !ok || !isASCIIDigit(mustDigi) {
				panic("TODO, syntax error")
			}

			// loop runs at least once
			for next, ok := s.peek(); ok && isASCIIDigit(next); next, ok = s.peek() {
				s.advance()
			}
		}
	}

}

func isIdentifierChar(r rune) bool {
	return r == '_' || isASCIIAlphanumeric(r)
}

func (s *Scanner) advanceTilEndOfNumber(canDot bool) {
	for next, ok := s.peek(); ok && (isASCIIDigit(next) || (canDot && next == '.')); next, ok = s.peek() {
		// if found dot, flip bool
		if next == '.' {
			canDot = false
		}

		// keep advancing until its no longer a digit
		s.advance()
	}
}

func (s *Scanner) advanceTilEndOfNumberBinary() {
	canDot := true

	for next, ok := s.peek(); ok && ((next == '1' || next == '0') || (canDot && next == '.')); next, ok = s.peek() {
		// if found dot, flip bool
		if next == '.' {
			canDot = false
		}

		// keep advancing until its no longer a digit
		s.advance()
	}
}

func (s *Scanner) scanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	return s.tokens
}

func (s *Scanner) addToken(tok Token) {
	s.tokens = append(s.tokens, tok)
}

func (s *Scanner) matchNext(expected rune) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	// consumes value if true
	s.current++
	return true
}
