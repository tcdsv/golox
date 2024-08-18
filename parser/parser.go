package parser

import (
	loxerror "golox/error"
	"golox/expr"
	"golox/stmt"
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

func (p *Parser) Parse() ([]stmt.Stmt, []error) {
	statements := []stmt.Stmt{}
	errors := []error{}
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			errors = append(errors, err)
		} else {
			statements = append(statements, statement)
		}
	}
	return statements, errors
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == tkn.EOF
}

func (p *Parser) check(tokenType tkn.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
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

func (p *Parser) declaration() (stmt.Stmt, error) {
	var stmt stmt.Stmt
	var err error
	if p.match(tkn.VAR) {
		stmt, err =  p.varDeclaration()
	} else {
		stmt, err = p.statement()
	}
	if err != nil {
		p.synchronize()
		return nil, err
	}
	return stmt, nil
}

func (p *Parser) varDeclaration() (stmt.Stmt, error) {

	err := p.consume(tkn.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}

	name := p.previous()
	var initializer expr.Expr
	if p.match(tkn.EQUAL) {
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		} 
	}

	err = p.consume(tkn.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return stmt.VarStmt{
		Name: name,
		Initializer: initializer,
	}, nil

}

func (p *Parser) statement() (stmt.Stmt, error) {
	if p.match(tkn.PRINT) {
		return p.printStatement()
	}
	if p.match(tkn.LEFT_BRACE) {
		return p.blockStatement()
	}
	return p.expressionStatement()
}

func (p *Parser) blockStatement() (stmt.Stmt, error) {

	statements := []stmt.Stmt{}

	for !p.check(tkn.RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, statement)
	}

	err := p.consume(tkn.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil, err
	}
	
	return stmt.BlockStmt{
		Statements: statements,
	}, nil

}

func (p *Parser) printStatement() (stmt.Stmt, error) {

	e, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(tkn.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	
	return stmt.PrintStmt{
		E: e,
	}, nil

}

func (p *Parser) expressionStatement() (stmt.Stmt, error) {
	
	e, err := p.expression()
	if err != nil {
		return nil, err
	}
	err = p.consume(tkn.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}
	return stmt.ExprStmt{
		E: e,
	}, nil

}

func (p *Parser) expression() (expr.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (expr.Expr, error) {
	
	e, err := p.equality()
	if err != nil {
		return nil, err
	}
	
	if p.match(tkn.EQUAL) {

		equals := p.previous()
		rightAssignment, err := p.assignment()
		if err != nil {
			return nil, err
		}

		varExpr, ok := e.(expr.VariableExpr)
		if ok {
			assignExpr := expr.AssignExpr{
				Name: varExpr.Name,
				Right: rightAssignment,
			}
			return assignExpr, nil		
		}

		return nil, loxerror.NewErrorFromToken(equals, "Invalid assignment target.")
	} 
	return e, nil

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

	if p.match(tkn.IDENTIFIER) {
		return expr.VariableExpr{
			Name: p.previous(),
		}, nil
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

func (p *Parser) synchronize() {
	
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == tkn.SEMICOLON {
			return
		}
		switch (p.peek().Type) {
			case tkn.CLASS:
        	case tkn.FUN:
        	case tkn.VAR:
        	case tkn.FOR:
        	case tkn.IF:
        	case tkn.WHILE:
        	case tkn.PRINT:
        	case tkn.RETURN:
          		return;
		}
		p.advance()
	}

}