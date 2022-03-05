package refactor

import "testing"

func Test_validateImage(t *testing.T) {
	var val = DoRefactor("foo")
	if val != "refactor bar" {
		t.Errorf("incorrect return value")
	}
}

func Test_validateNoArg(t *testing.T) {
	var val = DoRefactor("")
	if val != "refactor nothing" {
		t.Errorf("incorrect return value")
	}
}
