package main

import "testing"

func TestPrepareMessage(t *testing.T) {
	for _, c := range []struct {
		branch, message, want string
	}{
		{"ARC-123", "My commit", "ARC-123: My commit"},
		{"master", "My commit", "My commit"},
		{"ARC-123", "ARC-123: My commit", "ARC-123: My commit"},
		{"testing", "My commit", "My commit"},
		{"ARC-456", "ARC-123: My commit", "ARC-456: My commit"},
	} {
		got := PrepareMessage(c.branch, c.message)
		if got != c.want {
			t.Errorf("PrepareMessage(%q, %q) == %q, want %q", c.branch, c.message, got, c.want)
		}
	}
}
