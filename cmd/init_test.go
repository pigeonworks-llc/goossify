package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateProjectDirectory(t *testing.T) {
	t.Helper()

	tempDir := t.TempDir()

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		if chdirErr := os.Chdir(originalWD); chdirErr != nil {
			t.Fatalf("failed to restore working directory: %v", chdirErr)
		}
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("failed to switch to temp dir: %v", err)
	}

	projectPath, err := createProjectDirectory("sample-app")
	if err != nil {
		t.Fatalf("createProjectDirectory returned error: %v", err)
	}

	expected := filepath.Join(tempDir, "sample-app")

	resolvedProject, err := filepath.EvalSymlinks(projectPath)
	if err != nil {
		t.Fatalf("failed to resolve project path symlinks: %v", err)
	}
	resolvedExpected, err := filepath.EvalSymlinks(expected)
	if err != nil {
		t.Fatalf("failed to resolve expected path symlinks: %v", err)
	}

	if resolvedProject != resolvedExpected {
		t.Fatalf("unexpected project path: got %s, want %s", resolvedProject, resolvedExpected)
	}

	if _, err := os.Stat(projectPath); err != nil {
		t.Fatalf("project directory not created: %v", err)
	}

	currentWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to read current working directory: %v", err)
	}

	resolvedCurrent, err := filepath.EvalSymlinks(currentWD)
	if err != nil {
		t.Fatalf("failed to resolve current working directory: %v", err)
	}
	resolvedTempDir, err := filepath.EvalSymlinks(tempDir)
	if err != nil {
		t.Fatalf("failed to resolve temp dir: %v", err)
	}

	if resolvedCurrent != resolvedTempDir {
		t.Fatalf("working directory changed unexpectedly: got %s, want %s", resolvedCurrent, resolvedTempDir)
	}
}
