// Package assets provides embedded frontend assets.
package assets

import (
	"fmt"
	"io/fs"
	"path"
	"strings"
	"testing"
)

func TestWebEmbedContainsIndex(t *testing.T) {
	data, err := fs.ReadFile(Web, "web/index.html")
	if err != nil {
		t.Fatalf("read index: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("expected index content")
	}
}

func TestWebEmbedContainsInterLicense(t *testing.T) {
	data, err := fs.ReadFile(Web, "web/licenses/inter-OFL-1.1.txt")
	if err != nil {
		t.Fatalf("read embedded Inter license: %v", err)
	}
	notice := string(data)
	for _, want := range []string{
		"Copyright 2016 The Inter Project Authors (https://github.com/rsms/inter)",
		"SIL OPEN FONT LICENSE Version 1.1 - 26 February 2007",
		"Permission is hereby granted, free of charge",
	} {
		if !strings.Contains(notice, want) {
			t.Errorf("embedded Inter license is missing %q", want)
		}
	}
}

func TestWebEmbedContainsLocalFonts(t *testing.T) {
	files := embeddedFiles(t)
	stylesheets := embeddedStylesheets(files)
	subsets := []string{"cyrillic", "cyrillic-ext", "greek", "greek-ext", "latin", "latin-ext", "vietnamese"}
	weights := []string{"400", "500", "600", "700"}
	for _, weight := range weights {
		for _, subset := range subsets {
			prefix := "inter-" + subset + "-" + weight + "-normal"
			requireInterFontFace(t, files, stylesheets, prefix, weight)
		}
	}
	requireReferencedWOFF2(t, files, stylesheets, "materialdesignicons-webfont")
}

func TestWebEmbedHasNoRemoteFontReferences(t *testing.T) {
	for name, data := range embeddedFiles(t) {
		content := strings.ToLower(string(data))
		for _, host := range []string{"fonts.googleapis.com", "fonts.gstatic.com"} {
			if strings.Contains(content, host) {
				t.Errorf("embedded frontend file %s references remote font host %s", name, host)
			}
		}
	}
}

func embeddedFiles(t *testing.T) map[string][]byte {
	t.Helper()
	files := make(map[string][]byte)
	err := fs.WalkDir(Web, "web", func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk embedded frontend file %s: %w", name, walkErr)
		}
		if entry.IsDir() {
			return nil
		}
		data, err := fs.ReadFile(Web, name)
		if err != nil {
			return fmt.Errorf("read embedded frontend file %s: %w", name, err)
		}
		files[name] = data
		return nil
	})
	if err != nil {
		t.Fatalf("walk embedded frontend: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("embedded frontend is empty; rebuild the frontend before running asset tests")
	}
	return files
}

func embeddedStylesheets(files map[string][]byte) string {
	var stylesheets strings.Builder
	for name, data := range files {
		if path.Ext(name) == ".css" {
			stylesheets.Write(data)
		}
	}
	return stylesheets.String()
}

func requireInterFontFace(t *testing.T, files map[string][]byte, stylesheets, prefix, weight string) {
	t.Helper()
	assetName := requireReferencedWOFF2(t, files, stylesheets, prefix)
	block := fontFaceBlockForAsset(stylesheets, assetName)
	if block == "" {
		t.Fatalf("embedded font %s has no associated @font-face block", assetName)
	}
	normalized := strings.NewReplacer("'", "", "\"", "").Replace(strings.Join(strings.Fields(block), ""))
	declarations := []string{
		"font-family:Inter;",
		"font-weight:" + weight + ";",
		"font-style:normal;",
		"font-display:swap;",
	}
	for _, declaration := range declarations {
		if !strings.Contains(normalized, declaration) {
			t.Errorf("@font-face for %s is missing %s: %s", assetName, declaration, block)
		}
	}
}

func fontFaceBlockForAsset(stylesheets, assetName string) string {
	remaining := stylesheets
	for {
		faceStart := strings.Index(remaining, "@font-face")
		if faceStart < 0 {
			return ""
		}
		remaining = remaining[faceStart+len("@font-face"):]
		faceEnd := strings.Index(remaining, "}")
		if faceEnd < 0 {
			return ""
		}
		block := remaining[:faceEnd]
		if strings.Contains(block, assetName) {
			return block
		}
		remaining = remaining[faceEnd+1:]
	}
}

func requireReferencedWOFF2(t *testing.T, files map[string][]byte, stylesheets, prefix string) string {
	t.Helper()
	for name, data := range files {
		base := path.Base(name)
		if !strings.HasPrefix(base, prefix) || !strings.HasSuffix(base, ".woff2") {
			continue
		}
		if len(data) < 4 || string(data[:4]) != "wOF2" {
			t.Fatalf("embedded font %s is not a valid WOFF2 payload", name)
		}
		if !strings.Contains(stylesheets, base) {
			t.Fatalf("embedded font %s is not referenced by the production CSS", name)
		}
		return base
	}
	t.Fatalf("embedded frontend is missing %s*.woff2; rebuild the frontend before running asset tests", prefix)
	return ""
}
