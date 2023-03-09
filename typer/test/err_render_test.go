package parser_typer_test

import (
	"fmt"
	"github.com/alecthomas/assert/v2"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer"
	"github.com/xplosunn/tenecs/typer/type_error"
	"testing"
)

func TestSuiteRenderErrorUppercasePackage(t *testing.T) {
	type Case struct {
		program  string
		expected string
	}

	cases := []Case{}

	cases = append(cases, Case{
		program: `package MyPackage`,
		expected: `
| 1 | package MyPackage
              ^ package name should start with a lowercase letter
`,
	})

	cases = append(cases, Case{
		program: `package MyPackage

`,
		expected: `
| 1 | package MyPackage
              ^ package name should start with a lowercase letter
| 2 | `,
	})

	cases = append(cases, Case{
		program: `
package MyPackage


`,
		expected: `| 1 | 
| 2 | package MyPackage
              ^ package name should start with a lowercase letter
| 3 | 
| 4 | `,
	})

	for i, testCase := range cases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			res, err := parser.ParseString(testCase.program)
			assert.NoError(t, err)

			_, err = typer.Typecheck(*res)
			assert.Error(t, err, "Didn't get an typererror")

			typecheckErr, ok := err.(*type_error.TypecheckError)
			assert.True(t, ok)

			rendered, err := type_error.Render(testCase.program, typecheckErr)
			assert.NoError(t, err)

			assert.Equal(t, testCase.expected, rendered)
		})
	}
}
