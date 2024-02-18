package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	// Identifiers + Literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	// Operators
	ASSIGN   = "="
	PLUS     = "+"
	MINUS    = "-"
	BANG     = "!"
	ASTERISK = "*"
	SLASH    = "/"
	LT       = "<"
	GT       = ">"
	EQ       = "=="
	NOT_EQ   = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	// Keywords
	FUNCTION = "FUNCTION" // THILHIHNA
	LET      = "LET"      // HUCHIN
	TRUE     = "TRUE"     // TAK
	FALSE    = "FALSE"    // ZUAU
	IF       = "IF"       // AHIHLEH
	ELSE     = "ELSE"     // AHIHKEILEH
	RETURN   = "RETURN"   // LEHKIK
)

var keywords = map[string]TokenType{
	"thilhihna":  FUNCTION,
	"huchin":     LET,
	"tak":        TRUE,
	"zuau":       FALSE,
	"ahihleh":    IF,
	"ahihkeileh": ELSE,
	"lehkik":     RETURN,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
