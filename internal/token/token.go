package token

// Global tokens
const (
	INVALID    = "INVALID"
	END        = "END"
	IDENTIFIER = "IDENTIFIER"
	INTEGER    = "INTEGER"
	WHILE      = "WHILE"
	IF         = "IF"
	ELSE       = "ELSE"
	PRINT      = "PRINT"
)

// Token - Creates a token struct
type Token struct {
	Type  string
	Value string
}

// Map for matching keywords
var keywords = map[string]string{
	"while": WHILE,
	"if":    IF,
	"else":  ELSE,
	"print": PRINT,
}

// IsKeyword looks in the keywords map to see
// if the identifier is a keyword
func IsKeyword(identifier string) string {
	if token, ok := keywords[identifier]; ok {
		return token
	}
	return IDENTIFIER
}
