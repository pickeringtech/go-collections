package main

import (
	"encoding/json"
	"os"
	"testing"
)

func writeFileForTest(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}

func readFileForTest(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	return string(b)
}

func readFileForTestErr(path string) (string, error) {
	b, err := os.ReadFile(path)
	return string(b), err
}

func readBadgeForTest(t *testing.T, path string) Badge {
	t.Helper()
	var b Badge
	if err := json.Unmarshal([]byte(readFileForTest(t, path)), &b); err != nil {
		t.Fatalf("unmarshal badge %s: %v", path, err)
	}
	return b
}
