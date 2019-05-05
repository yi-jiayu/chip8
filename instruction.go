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
	lo, hi uint8
}

func (instr instruction) addr() uint16 {
	return uint16(instr.hi&0x0f)<<8 + uint16(instr.lo)
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
	default:
		panic("illegal opcode")
	}
}
