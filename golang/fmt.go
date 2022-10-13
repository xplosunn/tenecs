package golang

import (
	"github.com/alecthomas/assert/v2"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func Fmt(t *testing.T, fileContent string) string {
	dir, err := os.MkdirTemp("", "")
	assert.NoError(t, err)
	filePath := filepath.Join(dir, t.Name()+".go")

	_, err = os.Create(filePath)

	contentBytes := []byte(fileContent)
	err = os.WriteFile(filePath, contentBytes, 0644)
	assert.NoError(t, err)

	cmd := exec.Command("gofmt", "-w", filePath)
	cmd.Dir = dir
	err = cmd.Run()
	if err != nil {
		t.Log(filePath)
	}
	assert.NoError(t, err)

	formatted, err := os.ReadFile(filePath)
	assert.NoError(t, err)

	return string(formatted)
}
