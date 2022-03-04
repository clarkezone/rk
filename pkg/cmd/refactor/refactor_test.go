package refactor

import "testing"

func Test_validateImage(t *testing.T) {
	var val = DoRefactor()
	if val != "refactor" {
		t.Errorf("incorrect return value")
	}
}
