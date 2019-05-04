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
			name: "0nnn",
			args: args{
				data: []uint8{0x0a, 0xbc},
			},
			want: instruction{
				addr: 0xabc,
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
