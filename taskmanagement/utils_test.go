package taskmanagement

import (
	"fmt"
	"os"
	"testing"
)

func TestCreatePathWithNoChunks(t *testing.T) {
	home, _ := os.UserHomeDir()

	result := createPath()
	expected := fmt.Sprintf("%v/%v", home, appFolderName)

	if result != expected {
		t.Fatalf("Expected: %v, Received: %v", expected, result)
	}
}

func TestCreatePathWithChunks(t *testing.T) {
	home, _ := os.UserHomeDir()
	result := createPath("sub", "inner.txt")
	expected := fmt.Sprintf("%v/%v/%v/%v", home, appFolderName, "sub", "inner.txt")

	if result != expected {
		t.Fatalf("Expected: %v, Received: %v", expected, result)
	}
}
