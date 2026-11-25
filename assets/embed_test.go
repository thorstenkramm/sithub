// Package assets provides embedded frontend assets.
package assets

import (
	"fmt"
	"io/fs"
	"net/url"
	"path"
	"regexp"
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
	faces := embeddedFontFaces(t, files)
	subsets := []string{"cyrillic", "cyrillic-ext", "greek", "greek-ext", "latin", "latin-ext", "vietnamese"}
	weights := []string{"400", "500", "600", "700"}
	for _, weight := range weights {
		for _, subset := range subsets {
			prefix := "inter-" + subset + "-" + weight + "-normal"
			requireInterFontFace(t, faces, prefix, weight)
		}
	}
	requireWOFF2Face(t, faces, "materialdesignicons-webfont")
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

func requireInterFontFace(t *testing.T, faces map[string]string, prefix, weight string) {
	t.Helper()
	block := requireWOFF2Face(t, faces, prefix)
	normalized := strings.NewReplacer("'", "", "\"", "").Replace(strings.Join(strings.Fields(block), ""))
	declarations := []string{
		"font-family:Inter;",
		"font-weight:" + weight + ";",
		"font-style:normal;",
		"font-display:swap;",
	}
	for _, declaration := range declarations {
		if !strings.Contains(normalized, declaration) {
			t.Errorf("@font-face for %s is missing %s: %s", prefix, declaration, block)
		}
	}
}

func embeddedFontFaces(t *testing.T, files map[string][]byte) map[string]string {
	t.Helper()
	blocks := regexp.MustCompile(`@font-face\s*\{[^}]*\}`)
	sources := regexp.MustCompile(`url\(\s*["']?([^\s"')]+)["']?\s*\)`)
	faces := make(map[string]string)
	for name, data := range files {
		if path.Ext(name) != ".css" {
			continue
		}
		base := url.URL{Path: strings.TrimPrefix(name, "web")}
		for _, block := range blocks.FindAllString(string(data), -1) {
			for _, match := range sources.FindAllStringSubmatch(block, -1) {
				if file := checkedFontURL(t, files, &base, match[1]); path.Ext(file) == ".woff2" {
					faces[path.Base(file)] = block
				}
			}
		}
	}
	return faces
}

func checkedFontURL(t *testing.T, files map[string][]byte, base *url.URL, src string) string {
	t.Helper()
	ref, err := base.Parse(src)
	if err != nil || ref.Host != "" || (ref.Scheme != "" && ref.Scheme != "data") {
		t.Fatalf("%s: expected local font URL, got %q", base.Path, src)
	}
	if ref.Scheme == "data" {
		return ""
	}
	font, ok := files["web"+ref.Path]
	if !ok {
		t.Fatalf("%s: font URL %q resolves to missing asset %s", base.Path, src, ref.Path)
	}
	if path.Ext(ref.Path) == ".woff2" && (len(font) < 4 || string(font[:4]) != "wOF2") {
		t.Fatalf("invalid WOFF2 signature: %s", ref.Path)
	}
	return ref.Path
}

func requireWOFF2Face(t *testing.T, faces map[string]string, prefix string) string {
	t.Helper()
	for name, block := range faces {
		if strings.HasPrefix(name, prefix) {
			return block
		}
	}
	t.Fatalf("embedded frontend is missing %s*.woff2; rebuild the frontend before running asset tests", prefix)
	return ""
}
