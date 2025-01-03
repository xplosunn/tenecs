package node

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func RunCodeBlockingAndReturningOutputWhenFinished(t *testing.T, code string) (string, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	generatedFilePath := filepath.Join(dir, "main.js")
	if t != nil {
		t.Log("File ran: " + generatedFilePath)
	}
	_, err = os.Create(generatedFilePath)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(generatedFilePath, []byte(code), 0644)
	if err != nil {
		return "", err
	}

	runCmd := exec.Command("node", generatedFilePath)
	runCmd.Dir = dir
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return "", errors.New("error running " + generatedFilePath + ": " + err.Error())
	}
	return string(output), nil
}
