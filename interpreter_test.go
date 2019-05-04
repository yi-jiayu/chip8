package main

import (
	"reflect"
	"testing"
)

func Test_instrAt(t *testing.T) {
	type args struct {
		data []uint8
		i    int
	}
	tests := []struct {
		name string
		args args
		want instruction
	}{
		{
			name: "0nnn SYS addr",
			args: args{
				data: []uint8{0x0a, 0xbc},
			},
			want: instruction{
				opcode: OpSYS_0nnn,
				addr:   0xabc,
			},
		},
		{
			name: "00E0 CLS",
			args: args{
				data: []uint8{0, 0xe0},
			},
			want: instruction{
				opcode: OpCLS_00E0,
			},
		},
		{
			name: "00EE RET",
			args: args{
				data: []uint8{0, 0xee},
			},
			want: instruction{
				opcode: OpRET_00EE,
			},
		},
		{
			name: "1nnn JP addr",
			args: args{
				data: []uint8{0x1a, 0xbc},
			},
			want: instruction{
				opcode: OpJP_1nnn,
				addr:   0xabc,
			},
		},
		{
			name: "2nnn CALL addr",
			args: args{
				data: []uint8{0x2a, 0xbc},
			},
			want: instruction{
				opcode: OpCALL_2nnn,
				addr:   0xabc,
			},
		},
		{
			name: "3xkk SE Vx, byte",
			args: args{
				data: []uint8{0x3a, 0xbc},
			},
			want: instruction{
				opcode: OpSE_3xkk,
				x:      0xa,
				byte:   0xbc,
			},
		},
		{
			name: "4xkk SNE Vx, byte",
			args: args{
				data: []uint8{0x4a, 0xbc},
			},
			want: instruction{
				opcode: OpSNE_4xkk,
				x:      0xa,
				byte:   0xbc,
			},
		},
		{
			name: "5xy0 SE Vx, Vy",
			args: args{
				data: []uint8{0x5a, 0xbc},
			},
			want: instruction{
				opcode: OpSE_5xy0,
				x:      0xa,
				y:      0xb,
			},
		},
		{
			name: "6xkk LD Vx, byte",
			args: args{
				data: []uint8{0x5a, 0xbc},
			},
			want: instruction{
				opcode: OpSE_5xy0,
				x:      0xa,
				y:      0xb,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := instrAt(tt.args.data, tt.args.i); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("instrAt() = %v, want %v", got, tt.want)
			}
		})
	}
}
