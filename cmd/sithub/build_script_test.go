// Package main provides the SitHub server CLI.
package main

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildScriptContents(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok)

	root := filepath.Clean(filepath.Join(filepath.Dir(filename), "..", ".."))
	path := filepath.Join(root, "build.sh")

	// #nosec G304 -- reading a repo file in a test with a deterministic path.
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	contents := string(data)
	assert.Contains(t, contents, "npm ci")
	assert.Contains(t, contents, "npm run build")
	assert.Contains(t, contents, "tools/embed/copy.sh")
	assert.Contains(t, contents, "go build -o")
}
