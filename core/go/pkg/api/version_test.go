package api

import "testing"

func TestProtocolVersionDefined(t *testing.T) {
	if ProtocolVersion == "" {
		t.Fatal("ProtocolVersion must be defined")
	}
}
