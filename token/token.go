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
	IF       = "if"
	ELSE     = "else"
	RETURN   = "return"
	TRUE     = "true"
	FALSE    = "false"

	IDENT  = "IDENT" // Identifier
	INT    = "INT"   // Literal
	STRING = "STRING"

	// Operators
	ASSIGN = "="
	PLUS   = "+"
	MINUS  = "-"
	MUL    = "*"
	DIV    = "/"
	BANG   = "!"
	EQ     = "=="
	NEQ    = "!="

	LT = "<"
	GT = ">"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"
	LBRACKET = "["
	RBRACKET = "]"
)

var keywords map[string]TokenType = map[string]TokenType{
	"let":    LET,
	"fn":     FUNCTION,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

func LookupTokenType(literal string) TokenType {
	if tokType, ok := keywords[literal]; ok {
		return tokType
	}
	return IDENT
}
