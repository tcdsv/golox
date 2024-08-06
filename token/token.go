package token

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
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) Token {
	return Token{
		literal,
		tokenType,
		lexeme,
		line,
	}
}

func Identifier(text string) TokenType {
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