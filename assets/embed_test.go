// Package assets provides embedded frontend assets.
package assets

import (
	"io"
	"testing"
)

func TestWebEmbedContainsIndex(t *testing.T) {
	file, err := Web.Open("web/index.html")
	if err != nil {
		t.Fatalf("open index: %v", err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			t.Fatalf("close index: %v", err)
		}
	}()

	data, err := io.ReadAll(file)
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected index content")
	}
}
