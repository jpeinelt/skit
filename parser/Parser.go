package parser

import (
	"log"
	"peek/lexer"
	"peek/model"
	"strconv"
)

func isEOF(token lexer.Token) bool {
	return token.Typ == lexer.TOKEN_EOF
}

func Parse(input string) model.Presentation {
	output := model.Presentation{
		Slides: make([]model.Slide, 0),
	}

	var token lexer.Token
	slide := model.Slide{}

	log.Println("Start lexer and parser")

	l := lexer.BeginLexing(input)

	for {
		token = l.NextToken()
		if isEOF(token) {
			output.Slides = append(output.Slides, slide)
			break
		}
		switch token.Typ {
		case lexer.TOKEN_ERROR:
			panic(token.Val)
		case lexer.TOKEN_NEWSLIDE:
			output.Slides = append(output.Slides, slide)
			slide = model.Slide{}
		case lexer.TOKEN_TEXT:
			if len(slide.Text) == 0 {
				slide.Text = token.Val
			} else {
				slide.Text += "\n" + token.Val
			}
		case lexer.TOKEN_TITLE:
			slide.Title = token.Val
		case lexer.TOKEN_MEDIA:
			slide.Media = token.Val
		case lexer.TOKEN_COLOR_BG:
			c, err := strconv.Atoi(token.Val)
			if err != nil || c < 0 || c >= 256 {
				// panic("Parser expected background color code between 0 and 255 but was" + string(c))
				panic(c)
			}
			slide.ColorBg = c
		case lexer.TOKEN_COLOR_FG:
			c, err := strconv.Atoi(token.Val)
			if err != nil || c < 0 || c >= 256 {
				panic("Parser expected foreground color code between 0 and 255.")
			}
			slide.ColorFg = c
		default:
			// do nothing
		}
	}
	log.Println("Parser shutdown")
	return output
}
