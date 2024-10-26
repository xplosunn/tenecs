package node

import (
	"github.com/alecthomas/assert/v2"
	"testing"
)

func TestRunCodeBlockingAndReturningOutputWhenFinished(t *testing.T) {
	result, err := RunCodeBlockingAndReturningOutputWhenFinished(`console.log("hello world");`)
	assert.NoError(t, err)
	assert.Equal(t, "hello world\n", result)
}
