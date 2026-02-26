package api

import "testing"

func TestMethodForwardStartNonEmpty(t *testing.T) {
	if MethodForwardStart == "" {
		t.Fatal("MethodForwardStart must be non-empty")
	}
}
