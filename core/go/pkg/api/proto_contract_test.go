package api

import "testing"

func TestControlMethodForwardStartExists(t *testing.T) {
	if MethodForwardStart == "" {
		t.Fatal("MethodForwardStart must be defined from generated protobuf")
	}
}
