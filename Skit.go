package main

import (
	"encoding/json"
	"fmt"
	"github.com/jpeinelt/gocui"
	"log"
	"peek/model"
	"peek/parser"
)

var (
	currentSlide = 0
	presentation model.Presentation
)

func main() {
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

	out, err := json.Marshal(&presentation)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(out))

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
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, nextSlide); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyBackspace, gocui.ModNone, previousSlide); err != nil {
		return err
	}
	return nil
}

func newView(g *gocui.Gui, slide model.Slide) error {
	maxX, maxY := g.Size()
	g.SelBgColor = gocui.Attribute(slide.ColorBg)
	g.SelFgColor = gocui.Attribute(slide.ColorFg)
	title := slide.Title
	v, err := g.SetView(title, 0, 0, maxX, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, slide.Text)
	}
	return nil
}

func layout(g *gocui.Gui) error {
	return nil
}

func nextSlide(g *gocui.Gui, v *gocui.View) error {
	currentSlide++
	return newView(g, presentation.Slides[currentSlide])
}

func previousSlide(g *gocui.Gui, v *gocui.View) error {
	currentView := presentation.Slides[currentSlide].Title
	if err := g.DeleteView(currentView); err != nil {
		return err
	}
	currentSlide--
	return nil
}
