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

func (p *Parser) isAtEnd() bool {
	return p.current().Type == scanner.EOF
}

func (p *Parser) advance() {
	p.position++
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.position-1]
}

func (p *Parser) current() scanner.Token {
	return p.tokens[p.position]
}

// Pre-conditions:
// - Current position is defined.
// - List of tokens is provided.

// Post-conditions:
// - Returns false if the current position is at the end.
// - Returns false if the token at the current position does not match any token in the list.
// - Returns true if the token at the current position matches a token in the list.
// - In case of a match, the current position is incremented by one.

func (p *Parser) match(tokens ...scanner.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	for _, t := range tokens {
		if p.current().Type == t {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) expression() expr.Expr {
	return p.equality()
}

func (p *Parser) equality() expr.Expr {
	e := p.comparison()

	for p.match(scanner.BANG_EQUAL, scanner.EQUAL_EQUAL) {
		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    p.comparison(),
		}
	}

	return e
}

func (p *Parser) comparison() expr.Expr {
	e := p.term()

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    p.term(),
		}
	}

	return e
}

func (p *Parser) term() expr.Expr {
	e := p.factor()

	for p.match(scanner.PLUS, scanner.MINUS) {
		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    p.factor(),
		}
	}

	return e
}

func (p *Parser) factor() expr.Expr {
	e := p.unary()

	for p.match(scanner.SLASH, scanner.STAR) {
		e = expr.BinaryExpr{
			Operator: p.previous(),
			Left:     e,
			Right:    p.unary(),
		}
	}

	return e
}

func (p *Parser) unary() expr.Expr {
	if p.match(scanner.BANG, scanner.MINUS) {
		e := expr.UnaryExpr{
			Operator: p.previous(),
			Expr:     p.unary(),
		}
		return e
	}
	return p.primary()
}

func (p *Parser) primary() expr.Expr {
	if p.match(scanner.NUMBER, scanner.STRING, scanner.TRUE, scanner.FALSE, scanner.NIL) {
		e := expr.LiteralExpr{
			Value: p.previous().Literal,
		}
		return e
	}
	if p.match(scanner.LEFT_PAREN) {
		e := p.expression()
		if p.match(scanner.RIGHT_PAREN) {
			e = expr.GroupingExpr{
				Expr: e,
			}
			return e
		}
		// error
	}
	return nil
}