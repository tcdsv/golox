package parser

import (
	loxerror "golox/error"
	"golox/expr"
	"golox/stmt"
	tkn "golox/token"
	loxvalue "golox/value"
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

	err := p.consume(tkn.IDENTIFIER, loxerror.PARSE_ERROR_VARIABLE_EXPR_MISSING_NAME)
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

	err = p.consume(tkn.SEMICOLON, loxerror.PARSE_ERROR_VARIABLE_EXPR_MISSING_SEMICOLON)
	if err != nil {
		return nil, err
	}

	return stmt.VarStmt{
		Name: name,
		Initializer: initializer,
	}, nil

}

func (p *Parser) statement() (stmt.Stmt, error) {
	
	if p.match(tkn.FOR) {
		return p.forStatement()
	}
	if p.match(tkn.IF) {
		return p.ifStatement()
	}
	if p.match(tkn.PRINT) {
		return p.printStatement()
	}
	if p.match(tkn.LEFT_BRACE) {
		return p.blockStatement()
	}
	if p.match(tkn.WHILE) {
		return p.while()
	}
	return p.expressionStatement()

}

func (p *Parser) forStatement() (stmt.Stmt, error) {

	var err error
	err = p.consume(tkn.LEFT_PAREN, "Expect '(' after 'for'.");
	if err != nil {
		return nil, err
	}

	var initializer stmt.Stmt
	if p.match(tkn.SEMICOLON) {
		initializer = nil
	} else if p.match(tkn.VAR) {
		initializer, err = p.varDeclaration()
	} else {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}

	var condition expr.Expr
	if !p.check(tkn.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
    err = p.consume(tkn.SEMICOLON, "Expect ';' after loop condition.");
	if err != nil {
		return nil, err
	}

	var increment expr.Expr
	if !p.check(tkn.RIGHT_PAREN) {
		increment, err = p.expression();
		if err != nil {
			return nil, err
		}
	}
	err = p.consume(tkn.RIGHT_PAREN, "Expect ')' after for clauses.");
	if err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}

	if increment != nil {
		body = stmt.BlockStmt{
			Statements: []stmt.Stmt{
				body, stmt.ExprStmt{E: increment},
			},
		}
	}

	if condition == nil {
		condition = expr.LiteralExpr{
			Value: loxvalue.Boolean{Value: true},
		}
	}

	body = stmt.WhileStmt{
		Condition: condition,
		Body: body,
	}

	if initializer != nil {
		body = stmt.BlockStmt{
			Statements: []stmt.Stmt{
				initializer, body,
			},
		}
	}

	return body, nil

}


func (p *Parser) while() (stmt.Stmt, error) {

	p.consume(tkn.LEFT_PAREN, "Expect '(' after 'while'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(tkn.RIGHT_PAREN, "Expect '(' after while condition.")
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return stmt.WhileStmt{
		Condition: condition,
		Body: body,
	}, nil

}

func (p *Parser) ifStatement() (stmt.Stmt, error) {

	p.consume(tkn.LEFT_PAREN, "Expect '(' after 'if'.")
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	p.consume(tkn.RIGHT_PAREN, "Expect ')' after if condition.")

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	var elseBranch stmt.Stmt
	if p.match(tkn.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}

	return stmt.IfStmt{
		Condition: condition,
		ThenBrnach: thenBranch,
		ElseBranch: elseBranch,
	},nil

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

func (p *Parser) or() (expr.Expr, error) {
	
	and, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(tkn.OR) {

		operator := p.previous()

		right, err := p.and()
		if err != nil {
			return nil, err
		}

		and = expr.LogicalExpr{
			Operator: operator,
			Left: and,
			Right: right,
		}

	}

	return and, nil

}

func (p *Parser) and() (expr.Expr, error) {
	
	e, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(tkn.AND) {
		
		operator := p.previous()

		right, err := p.equality()
		if err != nil {
			return nil, err
		}

		e = expr.LogicalExpr{
			Operator: operator,
			Left: e,
			Right: right, 
		}

	}
	
	return e, nil

}

func (p *Parser) assignment() (expr.Expr, error) {
	
	e, err := p.or()
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

		err = p.consume(tkn.RIGHT_PAREN, loxerror.PARSE_ERROR_MISSING_RIGHT_PAREN)
		if err != nil {
			return nil, err
		}

		return expr.GroupingExpr{
			Expr: e,
		}, nil

	}

	return nil, loxerror.NewErrorFromToken(p.peek(), loxerror.PARSE_ERROR_MISSING_EXPRESSION)
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