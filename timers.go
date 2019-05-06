package main

import (
	"time"
)

// The Chip-8 timers run at 60 Hz
const timerInterval = time.Second / 60

func newTimer() (<-chan uint8, chan<- uint8, chan<- struct{}) {
	getch := make(chan uint8)
	setch := make(chan uint8)
	stopch := make(chan struct{})
	go func() {
		var val uint8
		ticker := time.NewTicker(timerInterval)
		defer ticker.Stop()
		for {
			select {
			case getch <- val:
			case v := <-setch:
				val = v
			case <-ticker.C:
				if val > 0 {
					val--
				}
			case <-stopch:
				return
			}
		}
	}()
	return getch, setch, stopch
}
