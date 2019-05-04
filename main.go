package main

import (
	"log"

	"github.com/gdamore/tcell"
)

func main() {
	log.Println("Creating new screen")
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Initialising screen")
	err = screen.Init()
	if err != nil {
		log.Fatal(err)
	}

	screen.SetContent(0, 0, 'H', nil, 0)

	screen.Show()

	log.Println("Polling for events")
	events := make(chan tcell.Event)
	go func() {
		for {
			ev := screen.PollEvent()
			events <- ev
		}
	}()
loop:
	for ev := range events {
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Rune() == 'q' {
				break loop
			}
			log.Println(ev.Name(), ev.Key())
		default:
			log.Println(ev)
		}
	}
}
