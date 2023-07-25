package dapper

var (
	openBrace        = []byte("{")
	openBraceNewLine = []byte("{\n")
	closeBrace       = []byte("}")
	openCloseBrace   = []byte("{}")

	openParen      = []byte("(")
	closeParen     = []byte(")")
	openCloseParen = []byte("()")

	keyValueSeparator = []byte(": ")

	space    = []byte(" ")
	asterisk = []byte("*")
	dot      = []byte(".")
	newLine  = []byte("\n")
)
