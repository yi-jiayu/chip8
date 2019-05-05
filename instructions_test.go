package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type registers struct {
	r [16]uint8
}

func (regs registers) set(idx int, val uint8) registers {
	regs.r[idx] = val
	return regs
}

func newInstructionAddr(addr uint16) instruction {
	return instruction{
		hi: uint8(addr >> 8),
		lo: uint8(addr & 0xFF),
	}
}

func newInstructionXYN(x, y, n uint8) instruction {
	return instruction{
		hi: x,
		lo: y<<4 + n,
	}
}

func newInstructionXKk(x, kk uint8) instruction {
	return instruction{
		hi: x,
		lo: kk,
	}
}

func TestSUB_8xy5(t *testing.T) {
	tests := []struct {
		name     string
		ip       Interpreter
		instr    instruction
		expected Interpreter
	}{
		{
			name: "Vx > Vy",
			ip: Interpreter{
				registers: [16]uint8{3, 1},
			},
			instr: newInstructionXYN(0, 1, 0),
			expected: Interpreter{
				registers: registers{[16]uint8{2, 1}}.set(VF, 1).r,
				pc:        1,
			},
		},
		{
			name: "Vx == Vy",
			ip: Interpreter{
				registers: [16]uint8{3, 3},
			},
			instr: newInstructionXYN(0, 1, 0),
			expected: Interpreter{
				registers: registers{[16]uint8{0, 3}}.set(VF, 1).r,
				pc:        1,
			},
		},
		{
			name: "Vx < Vy",
			ip: Interpreter{
				registers: [16]uint8{1, 3},
			},
			instr: newInstructionXYN(0, 1, 0),
			expected: Interpreter{
				registers: [16]uint8{0xFE, 3},
				pc:        1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SUB_8xy5(&tt.ip, tt.instr)
			if diff := cmp.Diff(tt.expected, tt.ip, cmp.AllowUnexported(Interpreter{})); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestSHL_8xyE(t *testing.T) {
	tests := []struct {
		name     string
		ip       Interpreter
		instr    instruction
		expected Interpreter
	}{
		{
			name: "MSB == 0",
			ip: Interpreter{
				registers: [16]uint8{0x70},
			},
			instr: newInstructionXYN(0, 0, 0),
			expected: Interpreter{
				registers: [16]uint8{0xE0},
				pc:        1,
			},
		},
		{
			name: "MSB == 1",
			ip: Interpreter{
				registers: [16]uint8{0xF0},
			},
			instr: newInstructionXYN(0, 0, 0),
			expected: Interpreter{
				registers: registers{[16]uint8{0xE0}}.set(VF, 1).r,
				pc:        1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SHL_8xyE(&tt.ip, tt.instr)
			if diff := cmp.Diff(tt.expected, tt.ip, cmp.AllowUnexported(Interpreter{})); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func BenchmarkADD_8xy4(b *testing.B) {
	ip := new(Interpreter)
	instr := instruction{0x01, 0x10}
	for i := 0; i < b.N; i++ {
		for x := uint8(0); x < 0xFF; x++ {
			for y := uint8(0); y < 0xFF; y++ {
				ip.registers[0] = x
				ip.registers[1] = y
				ADD_8xy4(ip, instr)
			}
		}
	}
}
