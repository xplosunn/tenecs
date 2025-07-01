package test_standard_library

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/codegen"
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	today := time.Now()
	day := today.Day()
	month := int(today.Month())
	year := today.Year()
	program := fmt.Sprintf(`package test

import tenecs.go.Runtime
import tenecs.test.UnitTest
import tenecs.test.GoIntegrationTest
import tenecs.test.GoIntegrationTestKit

_ := GoIntegrationTest("stdlib", "Time", (testkit: GoIntegrationTestKit, runtime: Runtime) => {
  today := runtime.time.today()
  testkit.assert.equal(%d, today.day)
  testkit.assert.equal(%d, today.month)
  testkit.assert.equal(%d, today.year)
})
`, day, month, year)
	parsed, err := parser.ParseString(program)
	assert.NoError(t, err)
	desugared := desugar.Desugar(*parsed)

	typed, err := typer.TypecheckSingleFile(desugared)
	assert.NoError(t, err)
	runTestInGolang(t, typed, codegen.FindTests(typed))
}
