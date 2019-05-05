package main

import (
	"math/rand"
	"time"
)

const (
	// TimestepSimulation is the clock speed of the Chip-8 emulator.
	TimestepSimulation = 2 * time.Millisecond

	// TimestepBatch is the time between executing batches of instructions.
	TimestepBatch = 100 * time.Millisecond
)

// Sprites for the Chip-8 hexadecimal font
var (
	Sprites = []uint8{
		// 0
		0xF0, 0x90, 0x90, 0x90, 0xF0,
		0, 0, 0, // padding

		// 1
		0x20, 0x60, 0x20, 0x20, 0x70,
		0, 0, 0, // padding

		// 2
		0xF0, 0x10, 0xF0, 0x80, 0xF0,
		0, 0, 0, // padding

		// 3
		0xF0, 0x10, 0xF0, 0x10, 0xF0,
		0, 0, 0, // padding

		// 4
		0x90, 0x90, 0xF0, 0x10, 0x10,
		0, 0, 0, // padding

		// 5
		0xF0, 0x80, 0xF0, 0x10, 0xF0,
		0, 0, 0, // padding

		// 6
		0xF0, 0x80, 0xF0, 0x90, 0xF0,
		0, 0, 0, // padding

		// 7
		0xF0, 0x10, 0x20, 0x40, 0x40,
		0, 0, 0, // padding

		// 8
		0xF0, 0x90, 0xF0, 0x90, 0xF0,
		0, 0, 0, // padding

		// 9
		0xF0, 0x90, 0xF0, 0x10, 0xF0,
		0, 0, 0, // padding

		// A
		0xF0, 0x90, 0xF0, 0x90, 0x90,
		0, 0, 0, // padding

		// B
		0xE0, 0x90, 0xE0, 0x90, 0xE0,
		0, 0, 0, // padding

		// C
		0xF0, 0x80, 0x80, 0x80, 0xF0,
		0, 0, 0, // padding

		// D
		0xE0, 0x90, 0x90, 0x90, 0xE0,
		0, 0, 0, // padding

		// E
		0xF0, 0x80, 0xF0, 0x80, 0xF0,
		0, 0, 0, // padding

		// F
		0xF0, 0x80, 0xF0, 0x80, 0x80,
		0, 0, 0, // padding
	}
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
	// dtget and stget are used for getting the current value of the delay and sound timers.
	// dtset and stset are used for setting the current value of the delay and sound timers.
	dtget <-chan uint8
	stget <-chan uint8
	dtset chan<- uint8
	stset chan<- uint8

	// The program counter (PC) should be 16-bit, and is used to store the currently executing address.
	pc uint16

	// The stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.
	sp uint8

	// The stack is an array of 16 16-bit values.
	stack    [16]uint16
	stackptr int

	// The original implementation of the Chip-8 language used a 64x32-pixel monochrome display.
	display [32][8]uint8

	// The current state of the display is sent to displaych whenever it is drawn.
	displaych chan<- [32][8]uint8

	// The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad.
	// Each receive will return a bitmask of the currently pressed keys.
	keypadch <-chan uint16

	stopch chan struct{}
}

func (ip *Interpreter) init() {
	ip.stopch = make(chan struct{})
	ip.displaych = make(chan [32][8]uint8)
}

// New returns a new Chip-8 interpreter.
func New(keypad <-chan uint16, display chan<- [32][8]uint8) *Interpreter {
	return &Interpreter{
		keypadch:  keypad,
		displaych: display,
	}
}

func (ip *Interpreter) Run() {
	ip.stopch = make(chan struct{})
	ip.loadSprites()

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
				ip.step()
				accum -= TimestepSimulation
			}
		case <-ip.stopch:
			return
		}
	}
}

func (ip *Interpreter) Stop() {
	ip.stopch <- struct{}{}
}

func (ip *Interpreter) render() {
	// non blocking send to the display
	select {
	case ip.displaych <- ip.display:
	default:
	}
}
func (ip *Interpreter) step() {
	// mock instructions which just move a cursor across the screen
	if ip.pc == 0 {
		ip.registers[1] = 7
		ip.registers[2] = 1
		ip.pc = 1
	}
	ip.registers[0]++
	if ip.registers[0] == 0 {
		section := ip.registers[1]
		val := ip.registers[2]
		if val != 1 {
			val >>= 1
		} else {
			ip.display[0][section] = 0
			section = (section + 1) % 8
			ip.registers[1] = section
			val = 0x80
		}
		ip.registers[2] = val
		ip.display[0][section] = val
		ip.render()
	}
}

func (ip *Interpreter) rand() uint8 {
	b := make([]byte, 1)
	rand.Read(b)
	return b[0]
}

func (ip *Interpreter) loadSprites() {
	copy(ip.memory[:], Sprites)
}
