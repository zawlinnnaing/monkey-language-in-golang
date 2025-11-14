package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// User identified token
	IDENT = "IDENT"
	INT   = "INT"

	ASSIGN = "="
	// TODO: add support for all operators (+,-,*,/)
	PLUS = "+"

	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	BANG     = "!"
	MINUS    = "-"
	SLASH    = "/"
	ASTERISK = "*"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	IF     = "IF"
	ELSE   = "ELSE"
	RETURN = "RETURN"
	TRUE   = "TRUE"
	FALSE  = "FALSE"

	STRING = "STRING"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"true":   TRUE,
	"false":  FALSE,
}

var AVAILABLE_TOKEN_TYPES []TokenType = []TokenType{
	ILLEGAL,
	EOF,
	IDENT,
	INT,
	ASSIGN,
	PLUS,
	COMMA,
	SEMICOLON,
	LPAREN,
	RPAREN,
	LBRACE,
	RBRACE,
	FUNCTION,
	LET,
}

func LookupIdentifier(token string) TokenType {
	if tok, ok := keywords[token]; ok {
		return tok
	}
	return IDENT
}

func New(tokenType TokenType, literal string) *Token {
	return &Token{
		Type:    tokenType,
		Literal: literal,
	}
}
