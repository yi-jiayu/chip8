package main

//go:generate ./instructions.sh

// Opcodes for standard Chip-8 instructions
const (
	OpCLS_00E0 opcode = iota
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

type opcode uint8

type instruction struct {
	hi uint8
	lo uint8
}

func (instr instruction) addr() uint16 {
	return uint16(instr.hi&0x0f)<<8 | uint16(instr.lo)
}

func (instr instruction) nibble() uint8 {
	return instr.lo & 0xf
}

func (instr instruction) x() uint8 {
	return instr.hi & 0xf
}

func (instr instruction) y() uint8 {
	return instr.lo & 0xf0 >> 4
}

func (instr instruction) byte() uint8 {
	return instr.lo
}

func (instr instruction) opcode() opcode {
	op := instr.hi & 0xf0 >> 4
	switch op {
	case 0:
		switch instr.lo {
		case 0xE0:
			return OpCLS_00E0
		case 0xEE:
			return OpRET_00EE
		default:
			return OpSYS_0nnn
		}
	case 1:
		return OpJP_1nnn
	case 2:
		return OpCALL_2nnn
	case 3:
		return OpSE_3xkk
	case 4:
		return OpSNE_4xkk
	case 5:
		return OpSE_5xy0
	case 6:
		return OpLD_6xkk
	case 7:
		return OpADD_7xkk
	case 8:
		switch instr.lo & 0xF {
		case 0:
			return OpLD_8xy0
		case 1:
			return OpOR_8xy1
		case 2:
			return OpAND_8xy2
		case 3:
			return OpXOR_8xy3
		case 4:
			return OpADD_8xy4
		case 5:
			return OpSUB_8xy5
		case 6:
			return OpSHR_8xy6
		case 7:
			return OpSUBN_8xy7
		case 0xE:
			return OpSHL_8xyE
		}
	case 9:
		return OpSNE_9xy0
	case 0xA:
		return OpLD_Annn
	case 0xB:
		return OpJP_Bnnn
	case 0xC:
		return OpRND_Cxkk
	case 0xD:
		return OpDRW_Dxyn
	case 0xE:
		switch instr.lo {
		case 0x9E:
			return OpSKP_Ex9E
		case 0xA1:
			return OpSKNP_ExA1
		}
	case 0xF:
		switch instr.lo {
		case 0x07:
			return OpLD_Fx07
		case 0x0A:
			return OpLD_Fx0A
		case 0x15:
			return OpLD_Fx15
		case 0x18:
			return OpLD_Fx18
		case 0x1E:
			return OpADD_Fx1E
		case 0x29:
			return OpLD_Fx29
		case 0x33:
			return OpLD_Fx33
		case 0x55:
			return OpLD_Fx55
		case 0x65:
			return OpLD_Fx65
		}
	}
	return 0xFF
}
