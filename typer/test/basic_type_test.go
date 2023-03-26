package parser_typer_test

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestBasicTypeTrue(t *testing.T) {
	validProgram(t, testcode.BasicTypeTrue)
}

func TestBasicTypeFalse(t *testing.T) {
	validProgram(t, testcode.BasicTypeFalse)
}
