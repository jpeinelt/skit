package lexer

import (
	"fmt"
)

type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_NEWLINE

	TOKEN_TEXT
	TOKEN_TITLE // #
	TOKEN_MEDIA // @

	TOKEN_ALIGN_LEFT   // <
	TOKEN_ALIGN_RIGHT  // >
	TOKEN_ALIGN_CENTER // =

	TOKEN_COLOR_BG // _
	TOKEN_COLOR_FG // ^

	TOKEN_MONOSPACE   // '
	TOKEN_BULLETPOINT // -
	TOKEN_COMMENT     // !
)

const EOF rune = 0
const CONTROL = "/"
const TITLE = "#"
const MEDIA = "@"
const ALIGN_LEFT = "<"
const ALIGN_RIGHT = ">"
const ALIGN_CENTER = "="
const COLOR_BG = "_"
const COLOR_FG = "^"
const MONOSPACE = "'"
const BULLETPOINT = "-"
const COMMENT = "!"

type Token struct {
	typ TokenType
	val string
}

func (t Token) String() string {
	switch t.typ {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_ERROR:
		return t.val
	}
	if len(t.val) > 10 {
		return fmt.Sprintf("%.10q...", t.val)
	}
	return fmt.Sprintf("%q", t.val)
}
