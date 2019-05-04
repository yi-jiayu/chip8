package main

// keypad is the interface for a Chip-8 keypad.
type keypad interface {
	// IsPressed returns true if key is currently pressed.
	IsPressed(key uint8) bool
}

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
	display [64][32]rune

	// The computers which originally used the Chip-8 Language had a 16-key hexadecimal keypad.
	keypad keypad
}

type instruction struct {
	addr   uint16
	nibble uint8
	x      uint8
	y      uint8
	byte   uint8
}

func instrAt(data []uint8, i int) (instr instruction) {
	hi, lo := data[i], data[i+1]
	op := hi & 0xf0
	switch op {
	case 0:
		instr.addr = uint16(hi&0x0f)<<8 + uint16(lo)
	}
	return
}
