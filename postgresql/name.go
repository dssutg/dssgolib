package postgresql

import "fmt"

func EscapePostgresIdent(ident string) string {
	if IsKeyword(ident) {
		return fmt.Sprintf("%q", ident)
	}
	return ident
}
