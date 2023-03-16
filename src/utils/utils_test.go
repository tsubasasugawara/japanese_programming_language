package utils

import (
	"testing"
)

func TestToLower(t *testing.T) {
	tests := []struct {
		input  string
		expect string
	}{
		{"０１９７２３５", "0197235"},
		{"７８６１２３１２８９６４", "786123128964"},
	}

	for i, v := range tests {
		if res := ToLower(v.input); res != v.expect {
			t.Fatalf("test%d : got=%s expect=%s\n", i, res, v.expect)
		}
	}
}
