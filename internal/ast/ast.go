package ast

import "gosql-br/internal/lexer"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

type Expression interface {
	Node
	expressionNode()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

type UseStatement struct {
	Token lexer.Token
	File  string
}

func (us *UseStatement) statementNode()       {}
func (us *UseStatement) TokenLiteral() string { return us.Token.Literal }

func (ss *SelectStatement) statementNode()       {}
func (ss *SelectStatement) TokenLiteral() string { return ss.Token.Literal }

type SelectStatement struct {
	Token     lexer.Token
	Columns   []string
	Condition Expression
}

type Condition struct {
	Left     string
	Operator string
	Right    string
}

type ComparisonExpression struct {
	Left     string // Nome da coluna
	Operator string // >, <, ==, !=, DENTRO
	Right    any    // Pode ser uma string, int ou uma lista []string para o DENTRO
}

func (ce *ComparisonExpression) TokenLiteral() string { return ce.Operator }
func (ce *ComparisonExpression) expressionNode()      {}

type LogicalExpression struct {
	Left     Expression // Pode ser uma Comparison ou outra Logical
	Operator string     // E, OU
	Right    Expression
}

func (le *LogicalExpression) TokenLiteral() string { return le.Operator }
func (le *LogicalExpression) expressionNode()      {}
