package main

import (
	"testing"
)

func Test_instruction_addr(t *testing.T) {
	instr := instruction{0x0A, 0xBC}
	var want uint16 = 0xABC
	if got := instr.addr(); got != want {
		t.Errorf("want = 0x%X, got = 0x%X", want, got)
	}
}

func Test_instruction_nibble(t *testing.T) {
	instr := instruction{0x0A, 0xBC}
	var want uint8 = 0xC
	if got := instr.nibble(); got != want {
		t.Errorf("want = 0x%X, got = 0x%X", want, got)
	}
}

func Test_instruction_x(t *testing.T) {
	instr := instruction{0x0A, 0xBC}
	var want uint8 = 0xA
	if got := instr.x(); got != want {
		t.Errorf("want = 0x%X, got = 0x%X", want, got)
	}
}

func Test_instruction_y(t *testing.T) {
	instr := instruction{0x0A, 0xBC}
	var want uint8 = 0xB
	if got := instr.y(); got != want {
		t.Errorf("want = 0x%X, got = 0x%X", want, got)
	}
}

func Test_instruction_byte(t *testing.T) {
	instr := instruction{0x0A, 0xBC}
	var want uint8 = 0xBC
	if got := instr.byte(); got != want {
		t.Errorf("want = 0x%X, got = 0x%X", want, got)
	}
}

func Test_instruction_opcode(t *testing.T) {
	tests := []struct {
		name        string
		instruction instruction
		want        opcode
	}{
		{
			name:        "0nnn SYS addr",
			instruction: instruction{0x0a, 0xbc},
			want:        OpSYS_0nnn,
		},
		{
			name:        "00E0 CLS",
			instruction: instruction{0, 0xe0},
			want:        OpCLS_00E0,
		},
		{
			name:        "00EE RET",
			instruction: instruction{0, 0xee},
			want:        OpRET_00EE,
		},
		{
			name:        "1nnn JP addr",
			instruction: instruction{0x1a, 0xbc},
			want:        OpJP_1nnn,
		},
		{
			name:        "2nnn CALL addr",
			instruction: instruction{0x2a, 0xbc},
			want:        OpCALL_2nnn,
		},
		{
			name:        "3xkk SE Vx, byte",
			instruction: instruction{0x3a, 0xbc},
			want:        OpSE_3xkk,
		},
		{
			name:        "4xkk SNE Vx, byte",
			instruction: instruction{0x4a, 0xbc},
			want:        OpSNE_4xkk,
		},
		{
			name:        "5xy0 SE Vx, Vy",
			instruction: instruction{0x5a, 0xbc},
			want:        OpSE_5xy0,
		},
		{
			name:        "6xkk LD Vx, byte",
			instruction: instruction{0x6a, 0xbc},
			want:        OpLD_6xkk,
		},
		{
			name:        " 7xkk ADD Vx, byte",
			instruction: instruction{0x7a, 0xbc},
			want:        OpADD_7xkk,
		},
		{
			name:        "8xy0 LD Vx, Vy",
			instruction: instruction{0x8a, 0xb0},
			want:        OpLD_8xy0,
		},
		{
			name:        "8xy1 OR Vx, Vy",
			instruction: instruction{0x8a, 0xb1},
			want:        OpOR_8xy1,
		},
		{
			name:        "8xy2 AND Vx, Vy",
			instruction: instruction{0x8a, 0xb2},
			want:        OpAND_8xy2,
		},
		{
			name:        "8xy3 XOR Vx, Vy",
			instruction: instruction{0x8a, 0xb3},
			want:        OpXOR_8xy3,
		},
		{
			name:        "8xy4 ADD Vx, Vy",
			instruction: instruction{0x8a, 0xb4},
			want:        OpADD_8xy4,
		},
		{
			name:        "8xy5 SUB Vx, Vy",
			instruction: instruction{0x8a, 0xb5},
			want:        OpSUB_8xy5,
		},
		{
			name:        "8xy6 SHR Vx {, Vy}",
			instruction: instruction{0x8a, 0xb6},
			want:        OpSHR_8xy6,
		},
		{
			name:        "8xy7 SUBN Vx, Vy",
			instruction: instruction{0x8a, 0xb7},
			want:        OpSUBN_8xy7,
		},
		{
			name:        "8xyE SHL Vx {, Vy}",
			instruction: instruction{0x8a, 0xbE},
			want:        OpSHL_8xyE,
		},
		{
			name:        "9xy0 SNE Vx, Vy",
			instruction: instruction{0x9a, 0xb0},
			want:        OpSNE_9xy0,
		},
		{
			name:        "Annn LD I, addr",
			instruction: instruction{0xAa, 0xbc},
			want:        OpLD_Annn,
		},
		{
			name:        "Bnnn JP V0, addr",
			instruction: instruction{0xBa, 0xbc},
			want:        OpJP_Bnnn,
		},
		{
			name:        "Cxkk RND Vx, byte",
			instruction: instruction{0xCa, 0xbc},
			want:        OpRND_Cxkk,
		},
		{
			name:        "Dxyn DRW Vx, Vy, nibble",
			instruction: instruction{0xDa, 0xbc},
			want:        OpDRW_Dxyn,
		},
		{
			name:        "Ex9E SKP Vx",
			instruction: instruction{0xEa, 0x9E},
			want:        OpSKP_Ex9E,
		},
		{
			name:        "ExA1 SKNP Vx",
			instruction: instruction{0xEa, 0xA1},
			want:        OpSKNP_ExA1,
		},
		{
			name:        "Fx07 LD Vx, DT",
			instruction: instruction{0xFa, 0x07},
			want:        OpLD_Fx07,
		},
		{
			name:        "Fx0A LD Vx, K",
			instruction: instruction{0xFa, 0x0A},
			want:        OpLD_Fx0A,
		},
		{
			name:        "Fx15 LD DT, Vx",
			instruction: instruction{0xFa, 0x15},
			want:        OpLD_Fx15,
		},
		{
			name:        "Fx18 LD ST, Vx",
			instruction: instruction{0xFa, 0x18},
			want:        OpLD_Fx18,
		},
		{
			name:        "Fx1E ADD I, Vx",
			instruction: instruction{0xFa, 0x1E},
			want:        OpADD_Fx1E,
		},
		{
			name:        "Fx29 LD F, Vx",
			instruction: instruction{0xFa, 0x29},
			want:        OpLD_Fx29,
		},
		{
			name:        "Fx33 LD B, Vx",
			instruction: instruction{0xFa, 0x33},
			want:        OpLD_Fx33,
		},
		{
			name:        "Fx55 LD [I], Vx",
			instruction: instruction{0xFa, 0x55},
			want:        OpLD_Fx55,
		},
		{
			name:        "Fx65 LD Vx, [I]",
			instruction: instruction{0xFa, 0x65},
			want:        OpLD_Fx65,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.instruction.opcode(); got != tt.want {
				t.Errorf("instruction.opcode() = %v, want %v", got, tt.want)
			}
		})
	}
}
