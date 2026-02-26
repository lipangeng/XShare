package main

import "testing"

func TestVersionStringNonEmpty(t *testing.T) {
	if version == "" {
		t.Fatal("version must not be empty")
	}
}
