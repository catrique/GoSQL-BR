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
	Condition *Condition
}

type Condition struct {
	Left     string
	Operator string
	Right    string
}
