package main

import (
	"time"
)

const (
	// TimestepSimulation is the clock speed of the Chip-8 emulator.
	TimestepSimulation = 2 * time.Millisecond

	// TimestepBatch is the time between executing batches of instructions.
	TimestepBatch = 100 * time.Millisecond
)

// Interpreter contains the current state of the Chip-8 interpreter as well as its connected hardware.
type Interpreter struct {
	// The Chip-8 language is capable of accessing up to 4KB (4,096 bytes) of RAM.
	memory [4096]uint8

	// Chip-8 has 16 general purpose 8-bit registers.
	registers [16]uint8

	// There is also a 16-bit register called I.
	i uint16

	// Chip-8 also has two special purpose 8-bit registers, for the delay and sound timers.
	dt uint8
	st uint8

	// dtch and stch are used for setting the delay timer register and sound timer register respectively.
	dtch chan uint8
	stch chan uint8

	// The program counter (PC) should be 16-bit, and is used to store the currently executing address.
	pc uint16

	// The stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.
	sp uint8

	// The stack is an array of 16 16-bit values.
	stack [16]uint8

	// The original implementation of the Chip-8 language used a 64x32-pixel monochrome display.
	display [32][64]uint8

	displaych chan [32][64]uint8

	// The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad.
	// Each receive will return a bitmask of the currently pressed keys.
	keypadch chan uint16

	stopch chan struct{}
}

func (i *Interpreter) init() {
	i.stopch = make(chan struct{})
	i.displaych = make(chan [32][64]uint8)
}

func (i *Interpreter) Run() {
	i.init()

	currentTime := time.Now()
	var accum time.Duration

	ticker := time.NewTicker(TimestepBatch)
	defer ticker.Stop()
	for {
		select {
		case newTime := <-ticker.C:
			frameTime := newTime.Sub(currentTime)
			currentTime = newTime
			accum += frameTime
			for accum >= TimestepSimulation {
				i.step()
				accum -= TimestepSimulation
			}
		case <-i.stopch:
			return
		}
	}
}

func (i *Interpreter) step() {
	// mock instructions which just move a cursor across the screen
	i.registers[1]++
	if i.registers[1] == 0 {
		x0 := i.registers[0]
		x1 := (x0 + 1) % 64
		i.display[0][x0] = 0
		i.display[0][x1] = 1
		i.registers[0] = x1

		// non blocking send to the display
		select {
		case i.displaych <- i.display:
		default:
		}
	}

}
