package parser

import (
	"errors"
	"fmt"
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

func (p *Parser) Parse() expr.Expr {
	expr := p.expression()
	return expr
}

func handleError(err *error) {

	fmt.Println((*err).Error())

	if r := recover(); r != nil {
		*err = errors.New("abs")
	}

	fmt.Println((*err).Error())

}

func (p *Parser) isAtEnd() bool {
	// todo:
	// remove this conditional.
	// The parser assumes that the list of tokens ends with an EOF token
	if p.position >= len(p.tokens) {
		return true
	}
	return p.peek().Type == scanner.EOF
}

func (p *Parser) check(tokenType scanner.TokenType) bool {
	/*if p.isAtEnd() {
		return false
	}*/
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

func (p *Parser) expression() expr.Expr {
	// testPanic()
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
