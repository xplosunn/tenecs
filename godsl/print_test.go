package godsl_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/godsl"
	"testing"
)

func TestPrintMain(t *testing.T) {
	main := godsl.NativeFunctionDeclaration("main").Parameters().Body(
		godsl.NativeFunctionInvocation().Import("fmt").Name("Println").Parameters(`"hello world"`),
	)

	assert.Equal(t, `import (
	"fmt"
)

func main() {
	fmt.Println("hello world")
}`, godsl.Print(main))
}
