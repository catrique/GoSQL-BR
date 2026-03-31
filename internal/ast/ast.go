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

type SelectStatement struct {
	Token       lexer.Token
	Columns     []string
	Functions   []FunctionCall
	Condition   Expression
	OrderBy     string
	IgnoreEmpty string
}

func (ss *SelectStatement) statementNode()       {}
func (ss *SelectStatement) TokenLiteral() string { return ss.Token.Literal }

type ComparisonExpression struct {
	Left     string
	Operator string
	Right    any
}

func (ce *ComparisonExpression) TokenLiteral() string { return ce.Operator }
func (ce *ComparisonExpression) expressionNode()      {}

type LogicalExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (le *LogicalExpression) TokenLiteral() string { return le.Operator }
func (le *LogicalExpression) expressionNode()      {}

type DaysBetweenExpression struct {
	Column1  string
	Column2  string
	Operator string
	Value    string
}

func (db *DaysBetweenExpression) TokenLiteral() string { return "DIAS_ENTRE" }
func (db *DaysBetweenExpression) expressionNode()      {}

type FunctionCall struct {
	Name   string
	Column string
	Inner  *FunctionCall
}
