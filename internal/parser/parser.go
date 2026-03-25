package parser

import (
	"fmt"
	"gosql-br/internal/ast"
	"gosql-br/internal/lexer"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  lexer.Token
	peekToken lexer.Token
	errors    []string
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

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
		return stmt
	}

	if !p.expectPeek(lexer.IDENT) {
		return nil
	}
	stmt.Columns = append(stmt.Columns, p.curToken.Literal)

	for p.peekTokenIs(lexer.COMMA) {
		p.nextToken()

		if !p.expectPeek(lexer.IDENT) {
			return nil
		}

		stmt.Columns = append(stmt.Columns, p.curToken.Literal)
	}
	if p.peekTokenIs(lexer.WHERE) {
		p.nextToken()
		p.nextToken()
		stmt.Condition = p.parseCondition()
	}

	return stmt
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

func (p *Parser) curTokenIs(t lexer.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t lexer.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t lexer.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t lexer.TokenType) {
	msg := fmt.Sprintf("esperava que o próximo token fosse %d, mas obteve %d", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseCondition() *ast.Condition {
	cond := &ast.Condition{}

	if !p.curTokenIs(lexer.IDENT) {
		return nil
	}
	cond.Left = p.curToken.Literal
	p.nextToken()

	if !p.isOperator(p.curToken.Type) {
		return nil
	}
	cond.Operator = p.curToken.Literal
	p.nextToken()

	if p.curToken.Type != lexer.INT && p.curToken.Type != lexer.STRING {
		return nil
	}
	cond.Right = p.curToken.Literal

	return cond
}

func (p *Parser) isOperator(t lexer.TokenType) bool {
	return t == lexer.GT || t == lexer.GE || t == lexer.LT ||
		t == lexer.LE || t == lexer.EQ || t == lexer.NE
}

func (p *Parser) Errors() []string {
	return p.errors
}
