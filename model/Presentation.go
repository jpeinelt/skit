package model

// Presentation represents a sequence of Slides.
type Presentation struct {
	Slides []Slide
}

// Slide represents a presentation slide with a
// title, text, media, background and foreground colors.
// All fields are optional.
type Slide struct {
	Title   string
	Text    string
	Media   string
	ColorBg int
	ColorFg int
}
