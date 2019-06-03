package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Keywords
	FUNCTION = "FN"
	LET      = "LET"

	IDENT = "IDENT" // Identifier
	INT   = "INT"   // Literal

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
)

var keywords map[string]TokenType = map[string]TokenType{
	"let": LET,
	"fn":  FUNCTION,
}

func LookupTokenType(literal string) TokenType {
	if tokType, ok := keywords[literal]; ok {
		return tokType
	}
	return IDENT
}
