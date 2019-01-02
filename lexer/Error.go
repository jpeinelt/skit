package lexer

const (
	lexerErrorUnexpectedEOF    string = "Lexer: Unexpected end of file."
	lexerErrorExpectedControl  string = "Lexer: Expected Control Character or /."
	lexerErrorExpectedTitle    string = "Lexer: Expected title text."
	lexerErrorExpectedMediaURL string = "Lexer: Expected image or gif path."
	lexerErrorExpectedColor    string = "Lexer: Expected color code between 0 and 255."
)
