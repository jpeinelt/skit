package lexer

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type lexFn func(*Lexer) lexfn

type Lexer struct {
	name   string
	input  string
	tokens chan Token
	state  lexfn

	start int
	pos   int
	width int
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) currentInput() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) emit(tokenType TokenType) {
	l.tokens <- Token{tokenType, val: l.input[l.start:l.pos]}
	l.start = l.pos
}

func (l *Lexer) errorfn(format string, args ...interface{}) lexFn {
	l.tokens <- Token{TOKEN_ERROR, val: fmt.Sprintf(format, args...)}
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

func (l *Lexer) NextToken() Token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			l.state = l.state(l)
		}
	}
	panic("Lexer invalid state (nextToken)")
}

func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *Lexer) Run() {
	for state := lexBegin; state != nil; {
		state = state(l)
	}
	l.shutdown()
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
	if l.isEOF {
		return nil
	}
	if strings.HasPrefix(lexer.currentInput(), CONTROL) {
		return lexControl
	} else {
		return lexText
	}
}

func lexControl(l *Lexer) lexFn {
	l.ignore()
	switch n := l.next(); n {
	case CONTROL:
		lexText
	case TITLE:
		lexTitle
	case MEDIA:
		lexMedia
	case COLOR_BG:
		lexColorBg
	case COLOR_FG:
		lexColorFg
	case COMMENT:
		lexComment
	default:
		return l.errorfn(lexerErrorExpectedControl)
	}
}

func (l *Lexer) ignoreUntilTextInLine() {
	for {
		l.ignore()
		n := l.next()
		if !unicode.IsSpace(n) || n == "\n" || n == EOF {
			break
		}
	}
}

func scanLine(l *Lexer, token TokenType, nextFn lexFn) lexFn {
	for {
		ch := l.next()
		if ch == "\n" || ch == EOF {
			l.backup()
			l.emit(token)
			nextFn
		}
	}
}

func lexText(l *Lexer) lexFn {
	return scanLine(l, TOKEN_TEXT, lexNewSlide)
}

func lexTitle(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := l.currentInput()
	if unicode.IsSpace(c) {
		return l.errorFn(lexerErrorExpectedTitle)
	}
	return scanLine(l, TOKEN_TITLE, lexNewSlide)
}

func lexNewSlide(l *Lexer) lexFn {
	n := l.next()
	p := l.peek()
	if n == "\n" && p == "\n" {
		l.emit(TOKEN_NEWSLIDE)
	}
	l.ignore()
	lexBegin
}

func lexMedia(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := l.currentInput()
	if unicode.IsSpace(c) {
		return l.errorFn(lexerErrorExpectedMediaUrl)
	}
	return scanLine(l, TOKEN_MEDIA, lexNewSlide)
}

func lexColorBg(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	c := lcurrentInput()
	if unicode.IsSpace(c) {
		return l.errorFn(lexerErrorExpectedColor)
	}
	lexNewSlide
}

func lexColorFg(l *Lexer) lexFn {

	lexNewSlide
}

func lexComment(l *Lexer) lexFn {
	l.ignoreUntilTextInLine()
	return scanline(l, TOKEN_COMMENT, lexNewSlide)
}
