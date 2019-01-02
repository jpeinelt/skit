package main

import (
	"encoding/json"
	"log"
	"peek/parser"
)

func main() {
	testInput := `
		/! first Slide
		/_ 222
		/^ 39
		/# Vacation Time is Open Source Time
		Check out my repos at github.com/jpeinelt and gitlab.com/jpeinelt.

		/! second Slide
		/_ 114
		/^ 238
		The best Band ever: harpodell.de

		/! third Slide
		/_ 0
		/^ 15
		/# AWWWWW
		/@ ./cat.png
	`
	parsedPresentation := parser.Parse(testInput)

	out, err := json.Marshal(&parsedPresentation)
	if err != nil {
		panic(err)
	}
	log.Println(string(out))
}
