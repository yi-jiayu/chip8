package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

// memoryOffsetProgram is the memory address that Chip-8 programs begin at.
const memoryOffsetProgram = 0x200

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

	dtstop chan<- struct{}
	ststop chan<- struct{}

	// The program counter (PC) should be 16-bit, and is used to store the currently executing address.
	pc uint16

	// The stack pointer (SP) can be 8-bit, it is used to point to the topmost level of the stack.
	sp uint8

	// The stack is an array of 16 16-bit values.
	stack [16]uint16

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

// Load loads a Chip-8 program into memory.
func (ip *Interpreter) Load(prog []byte) {
	copy(ip.memory[memoryOffsetProgram:], prog)
}

// Run starts the Chip-8 interpreter.
func (ip *Interpreter) Run() {
	// initialise stop channel
	ip.stopch = make(chan struct{})

	// load sprites
	ip.loadSprites()

	// start sound and delay timers
	ip.stget, ip.stset, ip.ststop = newTimer()
	ip.dtget, ip.dtset, ip.dtstop = newTimer()

	// set PC to program start address
	ip.pc = memoryOffsetProgram

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
			// stop delay and sound timers
			ip.ststop <- struct{}{}
			ip.dtstop <- struct{}{}
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

func (ip *Interpreter) currentInstr() instruction {
	return instruction{
		hi: ip.memory[ip.pc],
		lo: ip.memory[ip.pc+1],
	}
}

func (ip *Interpreter) step() {
	instr := ip.currentInstr()
	op := instr.opcode()
	log.Printf("opcode: %d, instr: 0x%02X%02X, pc: 0x%02X", op, instr.hi, instr.lo, ip.pc)
	switch op {
	case OpCLS_00E0:
		CLS_00E0(ip, instr)
	case OpRET_00EE:
		RET_00EE(ip, instr)
	case OpSYS_0nnn:
		SYS_0nnn(ip, instr)
	case OpJP_1nnn:
		JP_1nnn(ip, instr)
	case OpCALL_2nnn:
		CALL_2nnn(ip, instr)
	case OpSE_3xkk:
		SE_3xkk(ip, instr)
	case OpSNE_4xkk:
		SNE_4xkk(ip, instr)
	case OpSE_5xy0:
		SE_5xy0(ip, instr)
	case OpLD_6xkk:
		LD_6xkk(ip, instr)
	case OpADD_7xkk:
		ADD_7xkk(ip, instr)
	case OpLD_8xy0:
		LD_8xy0(ip, instr)
	case OpOR_8xy1:
		OR_8xy1(ip, instr)
	case OpAND_8xy2:
		AND_8xy2(ip, instr)
	case OpXOR_8xy3:
		XOR_8xy3(ip, instr)
	case OpADD_8xy4:
		ADD_8xy4(ip, instr)
	case OpSUB_8xy5:
		SUB_8xy5(ip, instr)
	case OpSHR_8xy6:
		SHR_8xy6(ip, instr)
	case OpSUBN_8xy7:
		SUBN_8xy7(ip, instr)
	case OpSHL_8xyE:
		SHL_8xyE(ip, instr)
	case OpSNE_9xy0:
		SNE_9xy0(ip, instr)
	case OpLD_Annn:
		LD_Annn(ip, instr)
	case OpJP_Bnnn:
		JP_Bnnn(ip, instr)
	case OpRND_Cxkk:
		RND_Cxkk(ip, instr)
	case OpDRW_Dxyn:
		DRW_Dxyn(ip, instr)
	case OpSKP_Ex9E:
		SKP_Ex9E(ip, instr)
	case OpSKNP_ExA1:
		SKNP_ExA1(ip, instr)
	case OpLD_Fx07:
		LD_Fx07(ip, instr)
	case OpLD_Fx0A:
		LD_Fx0A(ip, instr)
	case OpLD_Fx15:
		LD_Fx15(ip, instr)
	case OpLD_Fx18:
		LD_Fx18(ip, instr)
	case OpADD_Fx1E:
		ADD_Fx1E(ip, instr)
	case OpLD_Fx29:
		LD_Fx29(ip, instr)
	case OpLD_Fx33:
		LD_Fx33(ip, instr)
	case OpLD_Fx55:
		LD_Fx55(ip, instr)
	case OpLD_Fx65:
		LD_Fx65(ip, instr)
	default:
		panic(fmt.Sprintf("illegal opcode: 0x%X", op))
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
