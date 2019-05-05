package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gdamore/tcell"
)

func render(screen tcell.Screen, display [32][8]uint8) {
	for y, row := range display {
		for x0, col := range row {
			for x1 := 0; x1 < 8; x1++ {
				var c rune
				var mask uint8 = 1 << (7 - uint(x1))
				if col&mask > 0 {
					c = tcell.RuneBlock
				} else {
					c = ' '
				}
				x := x0*8 + x1
				screen.SetContent(x, y, c, nil, 0)
			}
		}
	}
}

func resizeTerminal(w, h int) {
	fmt.Printf("\033[8;%d;%dt", h, w)
}

func init() {
	f, err := os.Create("chip8.log")
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)
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

	display := make(chan [32][8]uint8)

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
