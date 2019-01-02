package lexer

import (
	"fmt"
	"log"
	"unicode"
	"unicode/utf8"
)

type lexFn func(*Lexer) lexFn

// Lexer with an input, tokens and a state.
// Internally it tracks also a start pointer, a position pointer and
// the width of the current lexed item.
type Lexer struct {
	input  string
	tokens chan Token
	state  lexFn

	start int
	pos   int
	width int
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) currentInput() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.start:l.pos])
	return r
}

func (l *Lexer) emit(tokenType TokenType) {
	l.tokens <- Token{Typ: tokenType, Val: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) errorFn(format string, args ...interface{}) lexFn {
	l.tokens <- Token{Typ: TOKEN_ERROR, Val: fmt.Sprintf(format, args...)}
	return nil
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) isEOF() bool {
	return l.pos >= utf8.RuneCountInString(l.input)
}

func (l *Lexer) isWhitespace() bool {
	ch, _ := utf8.DecodeRuneInString(l.input[l.pos:])
	return unicode.IsSpace(ch)
}

func (l *Lexer) next() rune {
	if l.pos >= utf8.RuneCountInString(l.input) {
		l.width = 0
		return EOF
	}
	result, width := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = width
	l.pos += l.width
	return result
}

// NextToken returns the next lexed Token in our channel. If there is none,
// it loops over doing nothing until we lexed a new Token.
func (l *Lexer) NextToken() Token {
	for {
		if l.state == nil {
			l.shutdown()
			return Token{Typ: TOKEN_EOF, Val: ""}
		}
		select {
		case token := <-l.tokens:
			return token
		default:
			l.state = l.state(l)
		}
	}
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// BeginLexing returns a new Lexer with a given input, a start state and
// a buffered Token channel.
func BeginLexing(input string) *Lexer {
	l := &Lexer{
		input:  input,
		state:  lexBegin,
		tokens: make(chan Token, 2),
	}
	return l
}

func (l *Lexer) shutdown() {
	close(l.tokens)
}

func (l *Lexer) skipWhitespace() {
	for {
		ch := l.next()
		if !unicode.IsSpace(ch) {
			break
		}
		if ch == EOF {
			l.emit(TOKEN_EOF)
			break
		}
		l.ignore()
	}
}

func lexBegin(l *Lexer) lexFn {
	l.skipWhitespace()
	if l.isEOF() {
		return nil
	}
	if l.currentInput() == CONTROL {
		return lexControl
	}
	return lexText
}

func lexControl(l *Lexer) lexFn {
	l.ignore()
	switch n := l.next(); n {
	case CONTROL:
		return lexText
	case TITLE:
		return lexTitle
	case MEDIA:
		return lexMedia
	case COLOR_BG:
		return lexColorBg
	case COLOR_FG:
		return lexColorFg
	case COMMENT:
		return lexComment
	default:
		return l.errorFn(lexerErrorExpectedControl)
	}
}

func (l *Lexer) ignoreUntilTextInLine() {
	for {
		l.ignore()
		n := l.next()
		if !unicode.IsSpace(n) || n == '\n' || n == EOF {
			break
		}
	}
}

func scanLine(l *Lexer, token TokenType) {
	for {
		ch := l.next()
		if ch == '\n' || ch == EOF {
			l.backup()
			l.emit(token)
			break
		}
	}
}

func lexText(l *Lexer) lexFn {
	scanLine(l, TOKEN_TEXT)
	return lexNewSlide
}

func lexTitle(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := l.currentInput()
	if unicode.IsSpace(c) {
		return l.errorFn(lexerErrorExpectedTitle)
	}
	scanLine(l, TOKEN_TITLE)
	return lexNewSlide
}

func lexNewSlide(l *Lexer) lexFn {
	n := l.next()
	p := l.peek()
	if n == '\n' && p == '\n' {
		l.emit(TOKEN_NEWSLIDE)
	}
	l.ignore()
	return lexBegin
}

func lexMedia(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := l.currentInput()
	if unicode.IsSpace(c) {
		return l.errorFn(lexerErrorExpectedMediaURL)
	}
	scanLine(l, TOKEN_MEDIA)
	return lexNewSlide
}

func lexColorBg(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := l.currentInput()
	if unicode.IsSpace(c) || !unicode.IsNumber(c) {
		return l.errorFn(lexerErrorExpectedColor)
	}
	for {
		n := l.next()
		if n == EOF || n == '\n' {
			l.backup()
			l.emit(TOKEN_COLOR_BG)
			break
		}
		if !unicode.IsNumber(c) {
			return l.errorFn(lexerErrorExpectedColor)
		}
	}
	return lexNewSlide
}

func lexColorFg(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := l.currentInput()
	if unicode.IsSpace(c) || !unicode.IsNumber(c) {
		return l.errorFn(lexerErrorExpectedColor)
	}
	for {
		n := l.next()
		if n == EOF || n == '\n' {
			l.backup()
			l.emit(TOKEN_COLOR_FG)
			break
		}
		if !unicode.IsNumber(c) {
			return l.errorFn(lexerErrorExpectedColor)
		}
	}
	return lexNewSlide
}

func lexComment(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	scanLine(l, TOKEN_COMMENT)
	return lexNewSlide
}
