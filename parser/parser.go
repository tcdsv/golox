package parser

import (
	"golox/expr"
	"golox/scanner"
)

type Parser struct {
	tokens   []scanner.Token
	position int
}

func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *Parser) Parse() (expr.Expr, error) {
	return p.expression()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	return p.peek().Type == tokenType
}

func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.position++
	}
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.position-1]
}

func (p *Parser) peek() scanner.Token {
	return p.tokens[p.position]
}

func (p *Parser) match(tokens ...scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, t := range tokens {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expression() (expr.Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (expr.Expr, error) {

	e, err := p.comparison()
	if err != nil {
		return nil, err
	}
	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}		

		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    right,
		}
	}

	return e, nil

}

func (p *Parser) comparison() (expr.Expr, error) {
	e, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		right, err := p.term()
		if err != nil {
			return nil, err
		}
		
		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    right,
		}
	}

	return e, nil
}

func (p *Parser) term() (expr.Expr, error) {

	e, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.PLUS, scanner.MINUS) {
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    right,
		}
	}

	return e, nil

}

func (p *Parser) factor() (expr.Expr, error) {

	e, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.SLASH, scanner.STAR) {
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    right,
		}
	}

	return e, nil

}

func (p *Parser) unary() (expr.Expr, error) {

	if p.match(scanner.BANG, scanner.MINUS) {
		uexp, err := p.unary()
		if err != nil {
			return nil, err
		}
		e := expr.UnaryExpr{
			Operator: p.previous(),
			Expr:     uexp,
		}
		return e, nil
	}
	return p.primary()

}

func (p *Parser) primary() (expr.Expr, error) {
	
	if p.match(scanner.NUMBER, scanner.STRING, scanner.TRUE, scanner.FALSE, scanner.NIL) {
		e := expr.LiteralExpr{
			Value: p.previous().Literal,
		}
		return e, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		err = p.consume(scanner.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return expr.GroupingExpr{
			Expr: e,
		}, nil

	}

	return nil, scanner.NewErrorFromToken(p.peek(), "Expect expression.")
}

func (p *Parser) consume(tokenType scanner.TokenType, message string) error {
	if (p.check(tokenType)) {
		p.advance()
		return nil
	}
	
	return scanner.NewErrorFromToken(p.peek(), message)
}