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
	case '|':
		s.addToken(NewToken(Pipe, s.source[s.start:s.current], s.line, nil))
		return nil
	case '&':
		s.addToken(NewToken(Ampersand, s.source[s.start:s.current], s.line, nil))
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
		dotAt := s.advanceTilEndOfNumber(true)

		var numberPart []rune
		if dotAt == nil {
			numberPart = s.source[s.start:s.current]
		} else if s.current-1 == *dotAt {
			numberPart = s.source[s.start:(*dotAt - 1)]
		} else {
			numberPart = s.source[s.start:*dotAt]
		}

		// this will mutate s.current, so we save it above
		sci, err := s.scanScientific()

		if err != nil {
			return err
		}

		toParse := string(numberPart)
		if sci != nil {
			// append and parse as float
			toParse = fmt.Sprintf("%s%s", toParse, sci.string())
		}

		s.addToken(NewToken(Number, s.source[s.start:s.current], s.line, nil))

		// i, err := strconv.ParseInt(toParse, 10, 64)

		// if err != nil {
		// 	return err
		// }

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

func (s *Scanner) isAllZeroes(startInc, endExcl int) bool {
	for i := startInc; i < endExcl; i++ {
		if s.source[i] != '0' {
			return false
		}
	}

	return true
}

type Scientific struct {
	num      []rune
	positive bool
}

func (s Scientific) string() string {
	pref := ""
	if s.positive {
		pref = "e"
	} else {
		pref = "e-"
	}

	return fmt.Sprintf("%s%s", pref, string(s.num))
}

func (s *Scanner) scanScientific() (*Scientific, error) {
	next, ok := s.peek()

	if !ok || (next != 'e' && next != 'E') {
		return nil, nil
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

			out := &Scientific{
				positive: true,
			}

			startIdx := s.current
			// loop runs at least once
			digits := make([]rune, 1)
			for next, ok := s.peek(); ok && isASCIIDigit(next); next, ok = s.peek() {
				digits = append(digits, next)
				s.advance()
			}
			endIdx := s.current + 1
			out.num = s.source[startIdx:endIdx]

			return out, nil
		} else if afterE == '-' || afterE == '+' {
			s.advance() // consume 'e'
			s.advance() // consume +,-

			out := &Scientific{
				positive: afterE == '+',
			}

			mustDigi, ok := s.peek()

			if !ok || !isASCIIDigit(mustDigi) {
				panic("TODO, syntax error")
			}

			startIdx := s.current

			// loop runs at least once
			digits := make([]rune, 1)
			for next, ok := s.peek(); ok && isASCIIDigit(next); next, ok = s.peek() {
				digits = append(digits, next)
				s.advance()
			}

			endIdx := s.current + 1
			out.num = s.source[startIdx:endIdx]

			return out, nil
		} else {
			// not scientific
			return nil, nil
		}
	} else {
		return nil, nil
	}
}

func isIdentifierChar(r rune) bool {
	return r == '_' || isASCIIAlphanumeric(r)
}

// the return type is not useful if canDot=false
func (s *Scanner) advanceTilEndOfNumber(canDot bool) *int {
	lookedForDot := canDot
	dotAt := (*int)(nil)

	for next, ok := s.peek(); ok && (isASCIIDigit(next) || (canDot && next == '.')); next, ok = s.peek() {
		// if found dot, flip bool
		if next == '.' {
			at := s.current
			dotAt = &at
			canDot = false
		}

		// keep advancing until its no longer a digit
		s.advance()
	}

	foundDot := lookedForDot && !canDot

	if foundDot {
		return dotAt
	} else {
		return nil
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
