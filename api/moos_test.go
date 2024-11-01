package api

import (
	"testing"
)

func Test_getRandomMoo(t *testing.T) {
	tests := []struct {
		name    string
		present bool
	}{
		{
			name:    "t1",
			present: true,
		},
		{
			name:    "t2",
			present: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(moos, getRandomMoo()); got != tt.present {
				t.Errorf("getRandomMoo() = %v, present %v", got, tt.present)
			}
		})
	}
}

func Test_contains(t *testing.T) {
	type args struct {
		s []string
		e string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "t1",
			args: args{
				s: []string{"abc", "xyz"},
				e: "abc",
			},
			want: true,
		},
		{
			name: "t1",
			args: args{
				s: []string{"abc", "xyz"},
				e: "000",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := contains(tt.args.s, tt.args.e); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}
