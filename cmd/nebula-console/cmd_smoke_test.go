package main

import (
	"bytes"
	"testing"
)

func TestCLISmoke(t *testing.T) {
	var out bytes.Buffer
	if err := runCommand([]string{"health"}, &out); err != nil {
		t.Fatal(err)
	}
	if out.String() != "gesture-nebula ok\n" {
		t.Fatalf("output %q", out.String())
	}
}
