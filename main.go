package main

import (
	"fmt"
	"log"

	"github.com/gdamore/tcell"
)

func render(screen tcell.Screen, display [32][64]uint8) {
	for y, row := range display {
		for x, col := range row {
			var c rune
			if col == 1 {
				c = tcell.RuneBlock
			} else {
				c = ' '
			}
			screen.SetContent(x, y, c, nil, 0)
		}
	}
}

func resizeTerminal(w, h int) {
	fmt.Printf("\033[8;%d;%dt", h, w)
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}

	err = screen.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer screen.Fini()

	// try to resize terminal
	oldw, oldh := screen.Size()
	resizeTerminal(64, 32)
	defer resizeTerminal(oldw, oldh)

	display := make(chan [32][64]uint8)

	keymap := NewKeymap(DvorakLayout)
	keych, keypad := NewKeypad(keymap)

	ip := New(keypad, display)
	go ip.Run()

	events := make(chan tcell.Event)
	go func() {
		for {
			ev := screen.PollEvent()
			events <- ev
		}
	}()

	screen.Show()

loop:
	for {
		select {
		case ev := <-events:
			if key, ok := ev.(*tcell.EventKey); ok {
				if key.Key() == tcell.KeyCtrlC {
					ip.Stop()
					break loop
				}
				if key.Key() == tcell.KeyRune {
					keych <- key.Rune()
				}
			}
		case display := <-display:
			render(screen, display)
			screen.Show()
		}
	}
}
