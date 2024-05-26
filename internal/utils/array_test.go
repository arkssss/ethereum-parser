package utils

import "testing"

func TestInStringArray(t *testing.T) {
	type args struct {
		arr    []string
		target string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"test not in array", args{
			arr:    []string{"1", "2"},
			target: "3",
		}, false},
		{"test in array", args{
			arr:    []string{"1", "2"},
			target: "1",
		}, true},
		{"test empty", args{
			arr:    nil,
			target: "1",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InStringArray(tt.args.arr, tt.args.target); got != tt.want {
				t.Errorf("InStringArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
