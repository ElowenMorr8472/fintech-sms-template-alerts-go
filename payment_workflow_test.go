package main

import "testing"

func TestNeedsConfirmation(t *testing.T) {
	tests := []struct {
		name  string
		cents int
		want  bool
	}{{"small", 9999, false}, {"threshold", 100000, true}, {"large", 250000, true}}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := NeedsConfirmation(Payment{AmountCents: tc.cents}); got != tc.want {
				t.Fatalf("got %v, want %v", got, tc.want)
			}
		})
	}
}
