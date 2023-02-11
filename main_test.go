package main

import (
	"testing"
)

func TestValidArguments(t *testing.T) {
	args := []string{"a", "b", "c", "d"}
	if !validateArguments(args) {
		t.Error("args did not validate when they should have")
	}
}

func TestInvalidArguments(t *testing.T) {
	args := []string{"a", "b", "c"}
	if validateArguments(args) {
		t.Error("args validated when they should not have")
	}
}
