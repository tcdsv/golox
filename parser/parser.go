package parser

import (
	loxerror "golox/error"
	"golox/expr"
	tkn "golox/token"
)

type Parser struct {
	tokens   []tkn.Token
	position int
}

func NewParser(tokens []tkn.Token) *Parser {
	return &Parser{
		tokens:   tokens,
		position: 0,
	}
}

func (p *Parser) Parse() (expr.Expr, error) {
	return p.expression()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tkn.EOF
}

func (p *Parser) check(tokenType tkn.TokenType) bool {
	return p.peek().Type == tokenType
}

func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.position++
	}
}

func (p *Parser) previous() tkn.Token {
	return p.tokens[p.position-1]
}

func (p *Parser) peek() tkn.Token {
	return p.tokens[p.position]
}

func (p *Parser) match(tokens ...tkn.TokenType) bool {
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
	for p.match(tkn.BANG_EQUAL, tkn.EQUAL_EQUAL) {
		operator := p.previous()

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}		

		e = expr.BinaryExpr{
			Operator: operator,
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

	for p.match(tkn.GREATER, tkn.GREATER_EQUAL, tkn.LESS, tkn.LESS_EQUAL) {
		operator := p.previous()

		right, err := p.term()
		if err != nil {
			return nil, err
		}
		
		e = expr.BinaryExpr{
			Operator: operator,
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

	for p.match(tkn.PLUS, tkn.MINUS) {
		operator := p.previous()

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		e = expr.BinaryExpr{
			Operator: operator,
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

	for p.match(tkn.SLASH, tkn.STAR) {
		operator := p.previous()
		
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		e = expr.BinaryExpr{
			Operator: operator,
			Left:     e,
			Right:    right,
		}
	}

	return e, nil

}

func (p *Parser) unary() (expr.Expr, error) {

	if p.match(tkn.BANG, tkn.MINUS) {
		operator := p.previous()
		uexp, err := p.unary()
		if err != nil {
			return nil, err
		}
		e := expr.UnaryExpr{
			Operator: operator,
			Right:    uexp,
		}
		return e, nil
	}
	return p.primary()

}

func (p *Parser) primary() (expr.Expr, error) {
	
	if p.match(tkn.NUMBER, tkn.STRING, tkn.TRUE, tkn.FALSE, tkn.NIL) {
		e := expr.LiteralExpr{
			Value: p.previous().Literal,
		}
		return e, nil
	}

	if p.match(tkn.LEFT_PAREN) {
		e, err := p.expression()
		if err != nil {
			return nil, err
		}

		err = p.consume(tkn.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}

		return expr.GroupingExpr{
			Expr: e,
		}, nil

	}

	return nil, loxerror.NewErrorFromToken(p.peek(), "Expect expression.")
}

func (p *Parser) consume(tokenType tkn.TokenType, message string) error {
	if (p.check(tokenType)) {
		p.advance()
		return nil
	}
	
	return loxerror.NewErrorFromToken(p.peek(), message)
}