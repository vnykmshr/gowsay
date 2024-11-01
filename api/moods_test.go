package api

import (
	"testing"
)

func Test_getRandomMood(t *testing.T) {
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
			if _, got := moods[getRandomMood()]; got != tt.present {
				t.Errorf("getRandomMood() = %v, present %v", got, tt.present)
			}
		})
	}
}
