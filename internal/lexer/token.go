package lexer

type TokenType uint8

const (
	ILLEGAL TokenType = iota
	EOF

	IDENT
	INT
	FLOAT
	STRING

	ASSIGN
	GT
	GE
	LT
	LE
	EQ
	NE

	SELECT
	WHERE
	USE
	FROM
	AND
	OR
	IN

	COMMA
	DOT
	ASTERISK
	LPAREN
	RPAREN
)

type Token struct {
	Type    TokenType
	Literal string
}

var keywords = map[string]TokenType{
	"PEGUE":  SELECT,
	"QUANDO": WHERE,
	"USE":    USE,
	"DE":     FROM,
	"E":      AND,
	"OU":     OR,
	"EM":     IN,
	"DENTRO": IN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
