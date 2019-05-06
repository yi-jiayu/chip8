package main

import (
	"fmt"
	"io/ioutil"
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
	// read rom data from stdin
	prog, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

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
	// load program
	ip.Load(prog)
	go ip.Run()

	screen.Show()

	// start display loop
	go func() {
		for display := range display {
			render(screen, display)
			screen.Show()
		}
	}()

	for {
		ev := screen.PollEvent()
		if key, ok := ev.(*tcell.EventKey); ok {
			if key.Key() == tcell.KeyCtrlC {
				ip.Stop()
				break
			}
			if key.Key() == tcell.KeyRune {
				keych <- key.Rune()
			}
		}
	}
}
