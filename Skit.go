package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"peek/model"
	"peek/parser"

	termbox "github.com/nsf/termbox-go"

	"github.com/jpeinelt/image2ascii/convert"

	_ "image/jpeg"
	_ "image/png"

	"github.com/jpeinelt/gocmd"
	"github.com/jpeinelt/gocui"
)

var (
	currentSlide   = 0
	presentation   model.Presentation
	converter      = convert.NewImageConverter()
	convertOptions = convert.DefaultOptions
)

func main() {
	flags := struct {
		Help    bool   `short:"h" long:"help" description:"Display usage" global:"true"`
		Version bool   `short:"v" long:"version" description:"Display version"`
		Load    string `short:"l" long:"load" description:"loads slides file from given path"`
	}{}

	gocmd.HandleFlag("Load", func(cmd *gocmd.Cmd, args []string) error {
		fileName := flags.Load
		if len(fileName) > 0 {
			skitFile(fileName)
		} else {
			demo()
		}
		return nil
	})

	// Init the app
	gocmd.New(gocmd.Options{
		Name:        "Skit",
		Version:     "1.0.0",
		Description: "A basic presentation app for the command line.",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}

func demo() {
	testInput := `
		/! first Slide
		/_ 222
		/^ 39
		/# Vacation Time is Open Source Time
		Check out my repos at github.com/jpeinelt and gitlab.com/jpeinelt.

		/! second Slide
		/# Best Band!
		/_ 114
		/^ 238
		The best Band ever: harpodell.de

		/! third Slide
		/_ 0
		/^ 15
		/# AWWWWW
		/@ ./cat.png
	`
	presentation = parser.Parse(testInput)
	present()
}

func skitFile(fileName string) {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Panic("Cannot read file")
	}
	input := string(buf)
	presentation = parser.Parse(input)
	present()
}

func present() {
	// ascii converter
	convertOptions.FixedWidth = 80
	convertOptions.FixedHeight = 40

	// create UI
	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)
	if err := initKeybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := newView(g, presentation.Slides[currentSlide]); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func initKeybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, nextSlide); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, previousSlide); err != nil {
		return err
	}
	return nil
}

func newView(g *gocui.Gui, slide model.Slide) error {
	maxX, maxY := g.Size()
	if len(slide.Media) > 0 {
		g.BgColor = gocui.Attribute(termbox.ColorBlack)
		g.FgColor = gocui.Attribute(termbox.ColorWhite)
	} else {
		g.BgColor = gocui.Attribute(slide.ColorBg)
		g.FgColor = gocui.Attribute(slide.ColorFg)
	}
	name := string(currentSlide)
	v, err := g.SetView(name, -1, 0, maxX, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		title := fmt.Sprintf(" %v ", slide.Title)
		v.Title = title
		if len(slide.Media) > 0 {
			g.BgColor = gocui.Attribute(termbox.ColorBlack)
			g.FgColor = gocui.Attribute(termbox.ColorBlack)
			fmt.Fprint(v, converter.ImageFile2ASCIIString(slide.Media, &convertOptions))
		} else {
			text := fmt.Sprintf("\n\n\n%v", slide.Text)
			fmt.Fprintln(v, text)
		}

	}
	return nil
}

func layout(g *gocui.Gui) error {
	return nil
}

func nextSlide(g *gocui.Gui, v *gocui.View) error {
	currentSlide++
	if currentSlide >= len(presentation.Slides) {
		currentSlide = len(presentation.Slides) - 1
	}
	return newView(g, presentation.Slides[currentSlide])
}

func previousSlide(g *gocui.Gui, v *gocui.View) error {
	if currentSlide == 0 {
		return nil
	}
	if err := g.DeleteView(string(currentSlide)); err != nil {
		return err
	}
	currentSlide--
	if currentSlide < 0 {
		currentSlide = 0
	}
	slide := presentation.Slides[currentSlide]
	g.BgColor = gocui.Attribute(slide.ColorBg)
	g.FgColor = gocui.Attribute(slide.ColorFg)
	return nil
}
