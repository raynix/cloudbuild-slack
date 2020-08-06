package main

import "testing"

func TestHello(t *testing.T) {
	expected := "Hello, world!"
	result := hello()
	if result != expected {
		t.Error("Not likely but hello world failed.")
	}
}
