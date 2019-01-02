package lexer

const (
	lexerErrorUnexpectedEof    string = "Unexpected end of file."
	lexerErrorExpectedControl  string = "Expected Control Character or /."
	lexerErrorExpectedTitle    string = "Expected title text."
	lexerErrorExpectedMediaUrl string = "Expected image or gif path."
	lexerErrorExpectedColor    string = "Expected color code between 0 and 256."
)
