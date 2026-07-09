package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMainHelp(t *testing.T) {
	t.Helper()

	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{"sithub", "--help"}
	main()
}

func TestVersionCommandPrintsVersion(t *testing.T) {
	originalVersion := version
	defer func() { version = originalVersion }()
	version = "9.9.9"

	cmd := newVersionCmd()
	var out bytes.Buffer
	cmd.SetOut(&out)

	require.NoError(t, cmd.Execute())
	assert.Equal(t, "9.9.9\n", out.String())
}
