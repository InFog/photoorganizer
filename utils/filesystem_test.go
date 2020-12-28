package utils

import (
	"testing"
)

func TestIsFileFalse(t *testing.T) {
	if IsFile("notafile.go") {
		t.Error("notafile.go should not exist")
	}
}

func TestIsFileTrue(t *testing.T) {
	if IsFile("utils/filesystem.go") {
		t.Error("utils/filesystem.go should exist")
	}
}
