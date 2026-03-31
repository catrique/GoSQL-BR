package parser

import (
	"fmt"
	"gosql-br/internal/ast"
	"gosql-br/internal/lexer"
)

const (
	_ int = iota
	LOWEST
	OR
	AND
	COMPARE
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case lexer.USE:
		return p.parseUseStatement()
	case lexer.SELECT:
		return p.parseSelectStatement()
	default:
		return nil
	}
}

func (p *Parser) parseSelectStatement() *ast.SelectStatement {
	stmt := &ast.SelectStatement{Token: p.curToken}

	if p.peekTokenIs(lexer.ASTERISK) {
		p.nextToken()
		stmt.Columns = append(stmt.Columns, "*")
	} else {
		p.nextToken()
		if err := p.readColumnOrFunction(stmt); err != nil {
			p.errors = append(p.errors, err.Error())
			return nil
		}
		for p.peekTokenIs(lexer.COMMA) {
			p.nextToken()
			p.nextToken()
			if err := p.readColumnOrFunction(stmt); err != nil {
				p.errors = append(p.errors, err.Error())
				return nil
			}
		}
	}

	if p.peekTokenIs(lexer.WHERE) {
		p.nextToken()
		p.nextToken()
		stmt.Condition = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(lexer.IGNORE) {
		p.nextToken()
		if !p.expectPeek(lexer.EMPTY) {
			return nil
		}
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.IgnoreEmpty = p.curToken.Literal
	}

	if p.peekTokenIs(lexer.ORDER_BY) {
		p.nextToken()
		if !p.expectPeek(lexer.BY) {
			return nil
		}
		if !p.expectPeek(lexer.IDENT) {
			return nil
		}
		stmt.OrderBy = p.curToken.Literal
	}

	return stmt
}

func (p *Parser) readColumnOrFunction(stmt *ast.SelectStatement) error {
	switch p.curToken.Type {
	case lexer.COUNT:
		fn, err := p.parseFunctionCall("CONTE")
		if err != nil {
			return err
		}
		stmt.Functions = append(stmt.Functions, *fn)

	case lexer.DISTINCT:
		fn, err := p.parseFunctionCall("DIFERENTES")
		if err != nil {
			return err
		}
		stmt.Functions = append(stmt.Functions, *fn)

	case lexer.MAX_NUM:
		fn, err := p.parseFunctionCall("MAX")
		if err != nil {
			return err
		}
		stmt.Functions = append(stmt.Functions, *fn)

	case lexer.MIN_NUM:
		fn, err := p.parseFunctionCall("MIN")
		if err != nil {
			return err
		}
		stmt.Functions = append(stmt.Functions, *fn)

	case lexer.MAX_DATE:
		fn, err := p.parseFunctionCall("MAX_DATA")
		if err != nil {
			return err
		}
		stmt.Functions = append(stmt.Functions, *fn)

	case lexer.MIN_DATE:
		fn, err := p.parseFunctionCall("MIN_DATA")
		if err != nil {
			return err
		}
		stmt.Functions = append(stmt.Functions, *fn)

	case lexer.PERCENTAGE:
		stmt.Functions = append(stmt.Functions, ast.FunctionCall{Name: "PORCENTAGEM"})

	case lexer.IDENT:
		stmt.Columns = append(stmt.Columns, p.curToken.Literal)

	default:
		return fmt.Errorf("esperava coluna ou função, obteve '%s'", p.curToken.Literal)
	}
	return nil
}

func (p *Parser) parseFunctionCall(name string) (*ast.FunctionCall, error) {
	fn := &ast.FunctionCall{Name: name}

	if !p.peekTokenIs(lexer.LPAREN) {
		return fn, nil
	}
	p.nextToken()
	p.nextToken()

	if p.curToken.Type == lexer.DISTINCT {
		inner, err := p.parseFunctionCall("DIFERENTES")
		if err != nil {
			return nil, err
		}
		fn.Inner = inner
	} else if p.curTokenIs(lexer.IDENT) {
		fn.Column = p.curToken.Literal
	} else {
		return nil, fmt.Errorf("argumento inválido para %s", name)
	}

	if !p.expectPeek(lexer.RPAREN) {
		return nil, fmt.Errorf("esperava ')' após argumento de %s", name)
	}
	return fn, nil
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	var left ast.Expression

	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken()
		left = p.parseExpression(LOWEST)
		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
	} else if p.curToken.Type == lexer.DAYS_BETWEEN {
		left = p.parseDaysBetween()
	} else {
		left = p.parseComparison()
	}

	for !p.peekTokenIs(lexer.EOF) &&
		!p.peekTokenIs(lexer.RPAREN) &&
		!p.peekTokenIs(lexer.IGNORE) &&
		!p.peekTokenIs(lexer.ORDER_BY) &&
		precedence < p.getPrecedence(p.peekToken.Type) {

		if !p.peekTokenIs(lexer.AND) && !p.peekTokenIs(lexer.OR) {
			break
		}
		p.nextToken()
		left = p.parseLogicalExpression(left)
	}

	return left
}

func (p *Parser) parseComparison() ast.Expression {
	if !p.curTokenIs(lexer.IDENT) {
		return nil
	}

	expr := &ast.ComparisonExpression{Left: p.curToken.Literal}
	p.nextToken()

	if p.curTokenIs(lexer.EQ) || p.curTokenIs(lexer.NE) {
		op := p.curToken.Literal
		p.nextToken()
		if p.curToken.Type == lexer.EMPTY {
			if op == "==" {
				expr.Operator = "VAZIO"
			} else {
				expr.Operator = "NAO_VAZIO"
			}
			expr.Right = nil
			return expr
		}
		expr.Operator = op
		expr.Right = p.curToken.Literal
		return expr
	}

	if !p.isOperator(p.curToken.Type) {
		return nil
	}
	expr.Operator = p.curToken.Literal
	opType := p.curToken.Type
	p.nextToken()

	if opType == lexer.IN {
		expr.Right = p.parseList()
	} else {
		expr.Right = p.curToken.Literal
	}

	return expr
}

func (p *Parser) parseDaysBetween() ast.Expression {
	expr := &ast.DaysBetweenExpression{}

	if !p.expectPeek(lexer.LPAREN) {
		return nil
	}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	expr.Column1 = p.curToken.Literal

	if !p.expectPeek(lexer.COMMA) {
		return nil
	}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	expr.Column2 = p.curToken.Literal

	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}

	p.nextToken()
	if !p.isOperator(p.curToken.Type) {
		p.errors = append(p.errors, fmt.Sprintf("DIAS_ENTRE: esperava operador, obteve '%s'", p.curToken.Literal))
		return nil
	}
	expr.Operator = p.curToken.Literal

	p.nextToken()
	expr.Value = p.curToken.Literal

	return expr
}

func (p *Parser) parseLogicalExpression(left ast.Expression) ast.Expression {
	expr := &ast.LogicalExpression{
		Left:     left,
		Operator: p.curToken.Literal,
	}
	precedence := p.getPrecedence(p.curToken.Type)
	p.nextToken()
	expr.Right = p.parseExpression(precedence)
	return expr
}

func (p *Parser) parseList() []string {
	var list []string
	if !p.curTokenIs(lexer.LPAREN) {
		return nil
	}
	p.nextToken()
	list = append(list, p.curToken.Literal)
	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.curToken.Literal)
	}
	if !p.expectPeek(lexer.RPAREN) {
		return nil
	}
	return list
}

func (p *Parser) parseUseStatement() *ast.UseStatement {
	stmt := &ast.UseStatement{Token: p.curToken}
	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.File = p.curToken.Literal
	if p.peekTokenIs(lexer.DOT) {
		p.nextToken()
		p.nextToken()
		stmt.File += "." + p.curToken.Literal
	}
	return stmt
}

func (p *Parser) getPrecedence(t lexer.TokenType) int {
	switch t {
	case lexer.OR:
		return OR
	case lexer.AND:
		return AND
	default:
		return LOWEST
	}
}

func (p *Parser) curTokenIs(t lexer.TokenType) bool  { return p.curToken.Type == t }
func (p *Parser) peekTokenIs(t lexer.TokenType) bool { return p.peekToken.Type == t }

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("esperava token %d, obteve '%s' (%d)", t, p.peekToken.Literal, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) isOperator(t lexer.TokenType) bool {
	return t == lexer.GT || t == lexer.GE || t == lexer.LT ||
		t == lexer.LE || t == lexer.EQ || t == lexer.NE || t == lexer.IN
}

func (p *Parser) Errors() []string { return p.errors }
