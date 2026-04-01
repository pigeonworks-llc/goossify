package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeCredentialScanning(t *testing.T) {
	t.Run("detects hardcoded secret patterns in Go files", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "main.go", `package main
const apiKey = "AKIAIOSFODNN7EXAMPLE"
`)
		a := New(dir)
		result := a.analyzeCredentialScanning()

		hasWarning := false
		for _, item := range result.Items {
			if item.Status == "missing" {
				hasWarning = true
			}
		}
		if !hasWarning {
			t.Error("expected warning for hardcoded AWS key pattern")
		}
	})

	t.Run("passes clean project", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "main.go", `package main
func main() {}
`)
		a := New(dir)
		result := a.analyzeCredentialScanning()

		for _, item := range result.Items {
			if item.Name == "Hardcoded secrets" && item.Status == "missing" {
				t.Error("false positive: clean project should not have secret warnings")
			}
		}
	})
}

func TestAnalyzeInternalReferences(t *testing.T) {
	t.Run("detects private IP addresses", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "config.go", `package config
const dbHost = "192.168.1.100"
`)
		a := New(dir)
		result := a.analyzeInternalReferences()

		hasWarning := false
		for _, item := range result.Items {
			if item.Status == "missing" {
				hasWarning = true
			}
		}
		if !hasWarning {
			t.Error("expected warning for private IP address")
		}
	})

	t.Run("detects internal domains", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "client.go", `package client
const endpoint = "http://api.internal.company.com/v1"
`)
		a := New(dir)
		result := a.analyzeInternalReferences()

		hasWarning := false
		for _, item := range result.Items {
			if item.Status == "missing" {
				hasWarning = true
			}
		}
		if !hasWarning {
			t.Error("expected warning for internal domain")
		}
	})

	t.Run("detects localhost references", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "server.go", `package server
const addr = "forgejo.localhost:3000"
`)
		a := New(dir)
		result := a.analyzeInternalReferences()

		hasWarning := false
		for _, item := range result.Items {
			if item.Status == "missing" {
				hasWarning = true
			}
		}
		if !hasWarning {
			t.Error("expected warning for localhost reference")
		}
	})

	t.Run("passes clean project", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "main.go", `package main
const endpoint = "https://api.example.com/v1"
`)
		a := New(dir)
		result := a.analyzeInternalReferences()

		for _, item := range result.Items {
			if item.Name == "Internal references" && item.Status == "missing" {
				t.Error("false positive: clean project should pass")
			}
		}
	})

	t.Run("ignores test files", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "server_test.go", `package server
const testAddr = "127.0.0.1:8080"
`)
		a := New(dir)
		result := a.analyzeInternalReferences()

		for _, item := range result.Items {
			if item.Name == "Internal references" && item.Status == "missing" {
				t.Error("test files should be excluded from internal reference checks")
			}
		}
	})
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
