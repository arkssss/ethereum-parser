package utils

import (
	"testing"
)

func TestHexToInt(t *testing.T) {
	type args struct {
		hex string
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"test hex 1", args{hex: "0x1"}, 1},
		{"test hex 1", args{hex: "0xA"}, 10},
		{"test hex 1", args{hex: "0x11"}, 17},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HexToInt(tt.args.hex); got != tt.want {
				t.Errorf("HexToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntToHex1(t *testing.T) {
	type args struct {
		i int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test int 1", args{i: 1}, "0x1"},
		{"test int 1", args{i: 10}, "0xa"},
		{"test int 1", args{i: 17}, "0x11"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IntToHex(tt.args.i); got != tt.want {
				t.Errorf("IntToHex() = %v, want %v", got, tt.want)
			}
		})
	}
}
