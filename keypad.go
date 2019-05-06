package main

import (
	"time"
)

const (
	KeyResetInterval = 50 * time.Millisecond
)

var (
	QwertyLayout = "x123qweasdzc4rfv"
	DvorakLayout = "q123',.aoe;j4puk"
)

type Keymap map[rune]uint16

// NewKeymap returns a mapping from standard keyboard keys to Chip-8 keypad keys.
//
// The Chip-8 keypad contains 16 keys in the following layout, labelled with their hexadecimal index:
//
// 	1 2 3 C
// 	4 5 6 D
// 	7 8 9 E
// 	A 0 B F
//
// To map these keys to
//
// 	1 2 3 4
// 	Q W E R
// 	A S D F
// 	Z X C V
//
// on a Qwerty keyboard, we pass a string with each letter in the corresponding index:
//
// 	x123qweasdzc4rfv
//
// x is the first letter because it corresponds to the 0 key on the Chip-8 keypad.
func NewKeymap(layout string) Keymap {
	keymap := make(Keymap)
	for i, r := range layout {
		keymap[r] = 1 << uint(i)
	}
	return keymap
}

// NewKeypad returns a pair of channels, the first sending incoming keyboard events
// and the second for reading the current keyboard state.
func NewKeypad(keymap Keymap) (chan<- rune, <-chan uint16) {
	var state uint16
	keych := make(chan rune)
	statech := make(chan uint16)
	ticker := time.NewTicker(KeyResetInterval)

	go func() {
		for {
			select {
			case <-ticker.C:
				// reset pressed keys every KeyResetInterval
				state = 0
			case key := <-keych:
				if mask, ok := keymap[key]; ok {
					state |= mask
				}
			case statech <- state:
			}
		}
	}()

	return keych, statech
}
