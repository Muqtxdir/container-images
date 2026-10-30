package main

import "testing"

func TestParseReference(t *testing.T) {
	tests := []struct {
		value string
		want  string
	}{
		{"ghcr.io/Muqtxdir/container-images/debian:stable", "ghcr.io/muqtxdir/container-images/debian:stable"},
		{"ghcr.io/muqtxdir/container-images/ubuntu:RC1", "ghcr.io/muqtxdir/container-images/ubuntu:RC1"},
		{"localhost:5000/Images/debian:latest", "localhost:5000/images/debian:latest"},
		{"ghcr.io/Muqtxdir/debian", "ghcr.io/muqtxdir/debian:latest"},
		{
			"ghcr.io/Muqtxdir/debian@sha256:8d161b9a394233296daab64e653f1ba877a0885b8e40943e9ed4a08bec17a71b",
			"ghcr.io/muqtxdir/debian@sha256:8d161b9a394233296daab64e653f1ba877a0885b8e40943e9ed4a08bec17a71b",
		},
	}
	for _, test := range tests {
		reference, err := parseReference(test.value)
		if err != nil {
			t.Errorf("parseReference(%q): %v", test.value, err)
			continue
		}
		if got := reference.Name(); got != test.want {
			t.Errorf("parseReference(%q) = %q, want %q", test.value, got, test.want)
		}
	}
}
