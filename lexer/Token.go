package lexer

import (
	"fmt"
)

// TokenType describes all possible Tokens.
type TokenType int

const (
	TOKEN_ERROR TokenType = iota
	TOKEN_EOF

	TOKEN_NEWSLIDE

	TOKEN_TEXT
	TOKEN_TITLE // #
	TOKEN_MEDIA // @

	TOKEN_COLOR_BG // _
	TOKEN_COLOR_FG // ^

	TOKEN_COMMENT // !
)

const EOF rune = 0
const CONTROL = '/'
const TITLE = '#'
const MEDIA = '@'
const COLOR_BG = '_'
const COLOR_FG = '^'
const COMMENT = '!'

// Token represents a TokenType and a Value.
type Token struct {
	Typ TokenType
	Val string
}

func (t Token) String() string {
	switch t.Typ {
	case TOKEN_EOF:
		return "EOF"
	case TOKEN_ERROR:
		return t.Val
	}
	if len(t.Val) > 10 {
		return fmt.Sprintf("%.10q...", t.Val)
	}
	return fmt.Sprintf("%q", t.Val)
}
