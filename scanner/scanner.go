package scanner

import (
	"fmt"
	loxerror "golox/error"
	tkn "golox/token"
	loxvalue "golox/value"
)

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []tkn.Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{source, 0, 0, 1, []tkn.Token{}}
}


func (s *Scanner) addToken(tokenType tkn.TokenType, literal loxvalue.LoxValue) {
	lexeme := s.source[s.start:s.current]
	token := tkn.NewToken(tokenType, lexeme, literal, s.line)
	s.tokens = append(s.tokens, token)
}

func (s *Scanner) advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) Scan() ([]tkn.Token, []error) {
	errors := []error{}
	for !s.isAtEnd() {
		err := s.scanToken()
		if err != nil {
			fmt.Println(err.Error())
			errors = append(errors, err)
		}
	}
	s.tokens = append(s.tokens, tkn.NewToken(tkn.EOF, "", nil, s.line))
	return s.tokens, errors
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) scanString() *loxerror.Error {

	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return loxerror.NewError(s.line, "", "unterminated string")
	}

	s.advance()
	text := s.source[s.start+1 : s.current-1]
	s.addToken(tkn.STRING, loxvalue.NewString(text))
	return nil
}

func (s *Scanner) scanToken() error {

	s.start = s.current
	c := s.advance()

	switch c {
	case ' ':
	case '\r':
	case '\t':
		break
	case '\n':
		s.line++
	case '(':
		s.addToken(tkn.LEFT_PAREN, nil)
	case ')':
		s.addToken(tkn.RIGHT_PAREN, nil)
	case '{':
		s.addToken(tkn.LEFT_BRACE, nil)
	case '}':
		s.addToken(tkn.RIGHT_BRACE, nil)
	case ',':
		s.addToken(tkn.COMMA, nil)
	case '.':
		s.addToken(tkn.DOT, nil)
	case '-':
		s.addToken(tkn.MINUS, nil)
	case '+':
		s.addToken(tkn.PLUS, nil)
	case ';':
		s.addToken(tkn.SEMICOLON, nil)
	case '*':
		s.addToken(tkn.STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(tkn.BANG_EQUAL, nil)
		} else {
			s.addToken(tkn.BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(tkn.EQUAL_EQUAL, nil)
		} else {
			s.addToken(tkn.EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(tkn.LESS_EQUAL, nil)
		} else {
			s.addToken(tkn.LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(tkn.GREATER_EQUAL, nil)
		} else {
			s.addToken(tkn.GREATER, nil)
		}
	case '/':
		if s.match('/') { //comment
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(tkn.SLASH, nil)
		}
	case '"':
		err := s.scanString()
		if err != nil {
			return err
		}
	default:
		if isDigit(c) {
			s.scanNumber()
		} else if isAlpha(c) {
			s.scanIdentifier()
		} else {
			return loxerror.NewError(s.line, "", "unexpected character")
		}
	}
	return nil
}

func (s *Scanner) scanNumber() error {
	for !s.isAtEnd() && isDigit(s.peek()) {
		s.advance()
	}

	if !s.isAtEnd() && s.peek() == '.' && isDigit(s.peekNext()) {
		s.advance()
		for !s.isAtEnd() && isDigit(s.peek()) {
			s.advance()
		}
	}

	loxNumber, _ := loxvalue.NewNumberFromText(s.source[s.start:s.current])
	s.addToken(tkn.NUMBER, loxNumber)
	return nil
}

func (s *Scanner) scanIdentifier() error {
	for !s.isAtEnd() && isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	s.addToken(tkn.Identifier(text), nil)
	return nil
}

func (s *Scanner) peekNext() byte {
	if s.current >= len(s.source)-1 {
		return 0
	}
	return s.source[s.current+1]
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c byte) bool {
	return isAlpha(c) || isDigit(c)
}
