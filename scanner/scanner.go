package scanner

import (
	loxerror "golox/error"
)

const (
	LEFT_PAREN TokenType = iota
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	// One or two character tokens.
	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	// Literals.
	IDENTIFIER
	STRING
	NUMBER

	// Keywords.
	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	TRUE
	VAR
	WHILE

	EOF
)

type TokenType int

type Token struct {
	Literal interface{}
	Type    TokenType
	Lexeme  string
	line    int
}

type Scanner struct {
	source  string
	start   int
	current int
	line    int
	tokens  []Token
}

func NewScanner(source string) *Scanner {
	return &Scanner{source, 0, 0, 1, []Token{}}
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		literal,
		tokenType,
		lexeme,
		line,
	}
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{
		literal,
		tokenType,
		text,
		s.line,
	})
}

func (s *Scanner) advance() byte {
	char := s.source[s.current]
	s.current++
	return char
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) Scan() ([]Token, error) {

	for !s.isAtEnd() {
		err := s.scanToken()
		if err != nil {
			return nil, err
		}
	}
	s.tokens = append(s.tokens, NewToken(EOF, "", nil, s.line))
	return s.tokens, nil
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

func (s *Scanner) scanString() error {

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
	s.addToken(STRING, text)
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
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '/':
		if s.match('/') { //comment
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
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

	value := s.source[s.start:s.current]
	s.addToken(NUMBER, value)
	return nil
}

func (s *Scanner) scanIdentifier() error {
	for !s.isAtEnd() && isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	s.addToken(s.identiferType(text), nil)
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

func (s *Scanner) identiferType(text string) TokenType {
	switch text {
	case "and":
		return AND
	case "class":
		return CLASS
	case "else":
		return ELSE
	case "false":
		return FALSE
	case "fun":
		return FUN
	case "for":
		return FOR
	case "if":
		return IF
	case "nil":
		return NIL
	case "or":
		return OR
	case "print":
		return PRINT
	case "return":
		return RETURN
	case "super":
		return SUPER
	case "this":
		return THIS
	case "true":
		return TRUE
	case "var":
		return VAR
	case "while":
		return WHILE
	default:
		return IDENTIFIER
	}
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
