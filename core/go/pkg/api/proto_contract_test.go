package api

import "testing"

func TestMethodForwardStartMatchesControlMethodToken(t *testing.T) {
	const expected = "forward.start"
	if MethodForwardStart != expected {
		t.Fatalf("MethodForwardStart = %q, want %q", MethodForwardStart, expected)
	}
}
