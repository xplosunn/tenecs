package golang

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestRunCodeUnlessCached(t *testing.T) {
	result := RunCodeUnlessCached(t, `package main
	
import "fmt"
func main() {
    fmt.Println("hello world")
}
`)
	assert.Equal(t, "hello world\n", result)
}
