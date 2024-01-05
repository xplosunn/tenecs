package golang

import (
	"errors"
	"fmt"
	"github.com/alecthomas/assert/v2"
	"hash/fnv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func RunCodeUnlessCached(t *testing.T, code string) string {
	wd, err := os.Getwd()
	assert.NoError(t, err)

	projectDir := wd

	for !strings.HasSuffix(projectDir, "tenecs") {
		projectDir = filepath.Dir(projectDir)
	}

	cacheDir := filepath.Join(projectDir, ".cache")

	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		err := os.Mkdir(cacheDir, os.ModePerm)
		assert.NoError(t, err)
	}

	hash := func() uint64 {
		h := fnv.New64a()
		_, err := h.Write([]byte(code))
		assert.NoError(t, err)
		return h.Sum64()
	}()
	cacheFile := filepath.Join(
		cacheDir,
		strings.ReplaceAll(
			filepath.Join(wd, fmt.Sprintf("%s-%d.txt", t.Name(), hash)),
			"/",
			"__",
		),
	)
	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		file, err := os.Create(cacheFile)
		assert.NoError(t, err)

		result, err := RunCodeBlockingAndReturningOutputWhenFinished(code)
		assert.NoError(t, err)

		_, err = file.WriteString(result)
		assert.NoError(t, err)
		return result
	} else {
		fileContent, err := os.ReadFile(cacheFile)
		assert.NoError(t, err)
		return string(fileContent)
	}

}

func RunCodeBlockingAndReturningOutputWhenFinished(code string) (string, error) {
	dir, err := os.MkdirTemp("", "")
	if err != nil {
		return "", err
	}
	generatedFilePath := filepath.Join(dir, "main.go")
	if err != nil {
		return "", err
	}
	_, err = os.Create(generatedFilePath)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(generatedFilePath, []byte(code), 0644)
	if err != nil {
		return "", err
	}
	buildCmd := exec.Command("go", "build", generatedFilePath)
	buildCmd.Dir = dir
	err = buildCmd.Run()
	if err != nil {
		return "", errors.New("error running " + generatedFilePath + ": " + err.Error())
	}

	runCmd := exec.Command("./main")
	runCmd.Dir = dir
	output, err := runCmd.CombinedOutput()
	if err != nil {
		return "", errors.New("error running " + generatedFilePath + ": " + err.Error())
	}
	return string(output), nil
}
