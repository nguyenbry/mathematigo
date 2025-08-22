package main

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

func (s *Scanner) scanToken() {
	r := s.advance()

	switch r {
	case ')':
		s.addToken(NewToken(CloseParen, s.source[s.start:s.current], s.line, nil))
	case '(':
		s.addToken(NewToken(OpenParen, s.source[s.start:s.current], s.line, nil))
	case '+':
		s.addToken(NewToken(Plus, s.source[s.start:s.current], s.line, nil))
	case '-':
		s.addToken(NewToken(Minus, s.source[s.start:s.current], s.line, nil))
	case '.':
		s.addToken(NewToken(Dot, s.source[s.start:s.current], s.line, nil))
	case '*':
		s.addToken(NewToken(Star, s.source[s.start:s.current], s.line, nil))
	case '/':
		s.addToken(NewToken(Slash, s.source[s.start:s.current], s.line, nil))
	case '<':
		if s.matchNext('=') {
			s.addToken(NewToken(Lteq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Lt, s.source[s.start:s.current], s.line, nil))
		}
	case '>':
		if s.matchNext('=') {
			s.addToken(NewToken(Gteq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Gt, s.source[s.start:s.current], s.line, nil))
		}
	case '=':
		if s.matchNext('=') {
			s.addToken(NewToken(EqEq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Eq, s.source[s.start:s.current], s.line, nil))
		}
	case '!':
		if s.matchNext('=') {
			s.addToken(NewToken(BangEq, s.source[s.start:s.current], s.line, nil))
		} else {
			s.addToken(NewToken(Bang, s.source[s.start:s.current], s.line, nil))
		}
	case '\'', '"':
		val := s.string(r)
		s.addToken(NewToken(String, s.source[s.start:s.current], s.line, val))
	case ' ', '\t', '\r':
		// do nothing
	case '\n':
		s.line++
	default:
		panic("TODO")
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
