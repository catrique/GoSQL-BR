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

// --- Statements ---

type UseStatement struct {
	Token lexer.Token
	File  string
}

func (us *UseStatement) statementNode()       {}
func (us *UseStatement) TokenLiteral() string { return us.Token.Literal }

type SelectStatement struct {
	Token       lexer.Token
	Columns     []string       // nomes de colunas normais
	Functions   []FunctionCall // funções como CONTE, MAX, PORCENTAGEM...
	Condition   Expression
	OrderBy     string // nome da coluna para ORDENE POR
	IgnoreEmpty string // nome da coluna cujos vazios devem ser ignorados
}

func (ss *SelectStatement) statementNode()       {}
func (ss *SelectStatement) TokenLiteral() string { return ss.Token.Literal }

// --- Expressions ---

type ComparisonExpression struct {
	Left     string
	Operator string // >, <, ==, !=, DENTRO, VAZIO, NAO_VAZIO
	Right    any    // string, int ou []string para DENTRO; nil para VAZIO
}

func (ce *ComparisonExpression) TokenLiteral() string { return ce.Operator }
func (ce *ComparisonExpression) expressionNode()      {}

type LogicalExpression struct {
	Left     Expression
	Operator string // E, OU
	Right    Expression
}

func (le *LogicalExpression) TokenLiteral() string { return le.Operator }
func (le *LogicalExpression) expressionNode()      {}

// DaysBetweenExpression representa: DIAS_ENTRE(coluna1, coluna2) >= N
// Coluna1 deve ser a data mais velha, Coluna2 a mais recente
type DaysBetweenExpression struct {
	Column1  string // data mais velha
	Column2  string // data mais recente
	Operator string // >, >=, <, <=, ==, !=
	Value    string // número de dias para comparar
}

func (db *DaysBetweenExpression) TokenLiteral() string { return "DIAS_ENTRE" }
func (db *DaysBetweenExpression) expressionNode()      {}

// --- Funções de Agregação ---

// FunctionCall representa: CONTE, DIFERENTES, MAX, MIN, MAX_DATA, MIN_DATA, PORCENTAGEM
// Exemplos:
//   CONTE             → Name="CONTE", Column=""
//   DIFERENTES(municipio) → Name="DIFERENTES", Column="municipio"
//   CONTE(DIFERENTES(municipio)) → Name="CONTE", Column="municipio", Inner=&FunctionCall{Name:"DIFERENTES"}
//   MAX(idade)        → Name="MAX", Column="idade"
//   PORCENTAGEM       → Name="PORCENTAGEM", Column=""
type FunctionCall struct {
	Name   string        // CONTE, DIFERENTES, MAX, MIN, MAX_DATA, MIN_DATA, PORCENTAGEM
	Column string        // coluna alvo (vazio para CONTE simples e PORCENTAGEM)
	Inner  *FunctionCall // para CONTE(DIFERENTES(...))
}
