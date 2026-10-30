package oci

import "testing"

func TestParseReference(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"ghcr.io/Muqtxdir/debian:RC1", "ghcr.io/muqtxdir/debian:RC1"},
		{"localhost:5000/Muqtxdir/debian", "localhost:5000/muqtxdir/debian:latest"},
	}
	for _, test := range tests {
		reference, err := ParseReference(test.value)
		if err != nil {
			t.Fatalf("ParseReference(%q): %v", test.value, err)
		}
		if got := reference.Name(); got != test.want {
			t.Errorf("ParseReference(%q) = %q, want %q", test.value, got, test.want)
		}
	}
}

func TestParseRepository(t *testing.T) {
	repository, err := ParseRepository("ghcr.io/Muqtxdir/Python")
	if err != nil {
		t.Fatal(err)
	}
	if got := repository.Name(); got != "ghcr.io/muqtxdir/python" {
		t.Errorf("ParseRepository() = %q, want ghcr.io/muqtxdir/python", got)
	}
}
