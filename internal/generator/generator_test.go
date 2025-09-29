package generator

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSanitizePackageName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  string
	}{
		{input: "MyApp", want: "myapp"},
		{input: "my-app", want: "my_app"},
		{input: "123app", want: "_123app"},
		{input: "こんにちは", want: "こんにちは"},
		{input: "", want: "app"},
	}

	for _, tc := range cases {
		if got := sanitizePackageName(tc.input); got != tc.want {
			t.Errorf("sanitizePackageName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestToExportedName(t *testing.T) {
	t.Parallel()

	cases := []struct {
		input string
		want  string
	}{
		{input: "my-app", want: "MyApp"},
		{input: "  multi word name  ", want: "MultiWordName"},
		{input: "", want: "App"},
	}

	for _, tc := range cases {
		if got := toExportedName(tc.input); got != tc.want {
			t.Errorf("toExportedName(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestGeneratorGenerateLibraryUsesBasePath(t *testing.T) {
	tempDir := t.TempDir()
	projectPath := filepath.Join(tempDir, "my-cli")

	cfg := &ProjectConfig{
		Name:           "my-cli",
		Type:           "library",
		Author:         "Example Author",
		Email:          "author@example.com",
		GitHubUsername: "example",
	}

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to read working directory: %v", err)
	}

	gen := New(projectPath, cfg)
	if err := gen.Generate(); err != nil {
		t.Fatalf("Generate() returned error: %v", err)
	}

	expectedFile := filepath.Join(projectPath, "pkg", cfg.PackageName, cfg.PackageName+".go")
	if _, err := os.Stat(expectedFile); err != nil {
		t.Fatalf("expected file not created: %v", err)
	}

	currentWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to read working directory after generation: %v", err)
	}

	if currentWD != originalWD {
		t.Fatalf("working directory changed: got %s, want %s", currentWD, originalWD)
	}
}
