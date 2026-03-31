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

	// Funções de agregação
	COUNT        // CONTE
	DISTINCT     // DIFERENTES
	ORDER_BY     // ORDENE
	BY           // POR
	IGNORE       // IGNORAR
	EMPTY        // VAZIO
	MAX_NUM      // MAX
	MIN_NUM      // MIN
	MAX_DATE     // MAX_DATA
	MIN_DATE     // MIN_DATA
	PERCENTAGE   // PORCENTAGEM
	DAYS_BETWEEN // DIAS_ENTRE

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
	"PEGUE":       SELECT,
	"QUANDO":      WHERE,
	"USE":         USE,
	"DE":          FROM,
	"E":           AND,
	"OU":          OR,
	"EM":          IN,
	"DENTRO":      IN,
	"CONTE":       COUNT,
	"DIFERENTES":  DISTINCT,
	"ORDENE":      ORDER_BY,
	"POR":         BY,
	"IGNORAR":     IGNORE,
	"VAZIO":       EMPTY,
	"MAX":         MAX_NUM,
	"MIN":         MIN_NUM,
	"MAX_DATA":    MAX_DATE,
	"MIN_DATA":    MIN_DATE,
	"PORCENTAGEM": PERCENTAGE,
	"DIAS_ENTRE":  DAYS_BETWEEN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
