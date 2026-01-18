package main

import (
	"os"
	"testing"
)

func TestMainHelp(t *testing.T) {
	t.Helper()

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"sithub", "--help"}
	main()
}
