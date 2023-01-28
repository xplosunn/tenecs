package testgen_test

import (
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/formatter"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/testgen"
	"github.com/xplosunn/tenecs/typer"
	"testing"
)

func TestOneLinerString(t *testing.T) {
	programString := `package pkg

interface HelloWorldProducer {
  public helloWorld: () -> String
}

implementing HelloWorldProducer module helloWorldProducer {
  public helloWorld := () => {
    "hello world!"
  }
}

`
	constructorName := "helloWorldProducer"
	targetFunctionName := "helloWorld"

	expectedOutput := `implementing UnitTests module generated() {
  public tests := (registry: UnitTestRegistry): UnitTestRegistry => {
    registry.test("hello world!", testCasehelloworld)
  }

  testCasehelloworld := (assert: UnitTestRegistry): UnitTestRegistry => {
    module := helloWorldProducer()
    result := module.helloWorld()
    expected := "hello world!"
    assert.equal<String>(result, expected)
  }

}`

	parsed, err := parser.ParseString(programString)
	assert.NoError(t, err)
	typed, err := typer.Typecheck(*parsed)
	assert.NoError(t, err)
	generated, err := testgen.Generate(*typed, constructorName, targetFunctionName)
	assert.NoError(t, err)
	formatted := formatter.DisplayModule(*generated)
	assert.Equal(t, expectedOutput, formatted)
}
