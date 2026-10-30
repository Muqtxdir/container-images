package main

import "testing"

func TestIsExcluded(t *testing.T) {
	tests := []struct {
		relativePath string
		want         bool
	}{
		{"bin/python3.14", false},
		{"bin/pip3.14", true},
		{"lib/libpython3.14.so.1.0", true},
		{"lib/python3.14/os.py", false},
		{"lib/python3.14t/site-packages/pip/__init__.py", true},
		{"lib/python3.14/tkinter/__init__.py", true},
		{"lib/tcl9.0/init.tcl", true},
		{"share/terminfo/x/xterm", true},
	}
	for _, test := range tests {
		if got := isExcluded(test.relativePath); got != test.want {
			t.Errorf("isExcluded(%q) = %v, want %v", test.relativePath, got, test.want)
		}
	}
}
