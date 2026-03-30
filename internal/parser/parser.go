package parser

import (
	"fmt"
	"gosql-br/internal/ast"
	"gosql-br/internal/lexer"
)

const (
	_ int = iota
	LOWEST
	OR      // OU
	AND     // E
	COMPARE // ==, !=, >, <, EM/DENTRO
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
	}

	if p.peekTokenIs(lexer.WHERE) {
		p.nextToken() // consome o token atual
		p.nextToken() // consome o WHERE
		stmt.Condition = p.parseExpression(LOWEST)
	}

	return stmt
}

// Hierarquia de expressões para lidar com E, OU e Parênteses
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var left ast.Expression

	// CASO 1: Início com Parênteses (Agrupamento)
	if p.curTokenIs(lexer.LPAREN) {
		p.nextToken() // Consome '('
		left = p.parseExpression(LOWEST)

		// Após resolver o interior, o PRÓXIMO deve ser ')'
		if !p.expectPeek(lexer.RPAREN) {
			return nil
		}
		// Aqui o curToken agora é ')'
	} else {
		// CASO 2: Comparação normal (coluna == valor)
		left = p.parseComparison()
	}

	// LOOP DE CONEXÃO: Tenta unir o 'left' com E/OU enquanto a prioridade permitir
	// Importante: o loop olha o PEEK (próximo) para decidir se continua
	for !p.peekTokenIs(lexer.EOF) && !p.peekTokenIs(lexer.RPAREN) && precedence < p.getPrecedence(p.peekToken.Type) {
		if !p.peekTokenIs(lexer.AND) && !p.peekTokenIs(lexer.OR) {
			break
		}

		p.nextToken() // Move para o E ou OU
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

func (p *Parser) parseLogicalExpression(left ast.Expression) ast.Expression {
	expr := &ast.LogicalExpression{
		Left:     left,
		Operator: p.curToken.Literal,
	}

	precedence := p.getPrecedence(p.curToken.Type)
	p.nextToken() // Move para o início da próxima expressão (depois do E/OU)
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

// Auxiliares
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
	msg := fmt.Sprintf("esperava %d, obteve %d", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) isOperator(t lexer.TokenType) bool {
	return t == lexer.GT || t == lexer.GE || t == lexer.LT ||
		t == lexer.LE || t == lexer.EQ || t == lexer.NE || t == lexer.IN
}

func (p *Parser) Errors() []string { return p.errors }
