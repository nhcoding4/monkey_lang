package main

type TokenType string

type Token struct {
	tokenType    TokenType
	literal      string
	line, column int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	IDENT  = "IDENT"
	INT    = "INT"
	FLOAT  = "FLOAT"
	STRING = "STRING"

	ASSIGN  = "="
	PLUS    = "+"
	MINUS   = "-"
	ASTERIX = "*"
	SLASH   = "/"
	MODULO  = "%"

	BANG  = "!"
	EQ    = "=="
	NOTEQ = "!="
	LT    = "<"
	GT    = ">"
	LTEQ  = "<="
	GTEQ  = ">="

	COMMA     = ","
	SEMICOLON = ";"
	COLON     = ":"

	LPAREN   = "("
	RPAREN   = ")"
	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	FUNCTION = "FUNCTION"
	LET      = "LET"
	IF       = "IF"
	ELSE     = "ELSE"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
)

func lookupIdent(ident string) TokenType {
	switch ident {
	case "fn":
		return FUNCTION
	case "let":
		return LET
	case "if":
		return IF
	case "else":
		return ELSE
	case "true":
		return TRUE
	case "false":
		return FALSE
	case "return":
		return RETURN
	default:
		return IDENT
	}
}
