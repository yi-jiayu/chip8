package main

type Opcode uint8

// Opcodes for standard Chip-8 instructions
const (
	OpCLS_00E0 Opcode = iota
	OpRET_00EE
	OpSYS_0nnn
	OpJP_1nnn
	OpCALL_2nnn
	OpSE_3xkk
	OpSNE_4xkk
	OpSE_5xy0
	OpLD_6xkk
	OpADD_7xkk
	OpLD_8xy0
	OpOR_8xy1
	OpAND_8xy2
	OpXOR_8xy3
	OpADD_8xy4
	OpSUB_8xy5
	OpSHR_8xy6
	OpSUBN_8xy7
	OpSHL_8xyE
	OpSNE_9xy0
	OpLD_Annn
	OpJP_Bnnn
	OpRND_Cxkk
	OpDRW_Dxyn
	OpSKP_Ex9E
	OpSKNP_ExA1
	OpLD_Fx07
	OpLD_Fx0A
	OpLD_Fx15
	OpLD_Fx18
	OpADD_Fx1E
	OpLD_Fx29
	OpLD_Fx33
	OpLD_Fx55
	OpLD_Fx65
)

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
	opcode Opcode
	addr   uint16
	nibble uint8
	x      uint8
	y      uint8
	byte   uint8
}

func instrAt(data []uint8, i int) (instr instruction) {
	hi, lo := data[i], data[i+1]
	op := hi & 0xf0 >> 4
	switch op {
	case 0:
		switch lo {
		case 0xE0:
			instr.opcode = OpCLS_00E0
		case 0xEE:
			instr.opcode = OpRET_00EE
		default:
			instr.opcode = OpSYS_0nnn
			instr.addr = uint16(hi&0x0f)<<8 + uint16(lo)
		}
	case 1:
		instr.opcode = OpJP_1nnn
		instr.addr = uint16(hi&0x0f)<<8 + uint16(lo)
	case 2:
		instr.opcode = OpCALL_2nnn
		instr.addr = uint16(hi&0x0f)<<8 + uint16(lo)
	case 3:
		instr.opcode = OpSE_3xkk
		instr.x = hi & 0x0f
		instr.byte = lo
	case 4:
		instr.opcode = OpSNE_4xkk
		instr.x = hi & 0x0f
		instr.byte = lo
	case 5:
		instr.opcode = OpSE_5xy0
		instr.x = hi & 0x0f
		instr.y = lo & 0xf0 >> 4
	}
	return
}
