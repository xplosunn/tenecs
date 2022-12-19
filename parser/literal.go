package parser

import "github.com/alecthomas/participle/v2"

var literalUnion = participle.Union[Literal](LiteralFloat{}, LiteralInt{}, LiteralString{}, LiteralBool{})

type Literal interface {
	sealedLiteral()
}

type LiteralFloat struct {
	Value float64 `@Float`
}

func (literal LiteralFloat) sealedLiteral() {}

type LiteralInt struct {
	Value int `@Int`
}

func (literal LiteralInt) sealedLiteral() {}

type LiteralString struct {
	Value string `@String`
}

func (literal LiteralString) sealedLiteral() {}

type LiteralBool struct {
	Value bool `@("true" | "false")`
}

func (literal LiteralBool) sealedLiteral() {}

func LiteralFold[Result any](
	literal Literal,
	caseFloat func(arg float64) Result,
	caseInt func(arg int) Result,
	caseString func(arg string) Result,
	caseBool func(arg bool) Result,
) Result {
	litFloat, ok := literal.(LiteralFloat)
	if ok {
		return caseFloat(litFloat.Value)
	}
	litInt, ok := literal.(LiteralInt)
	if ok {
		return caseInt(litInt.Value)
	}
	litString, ok := literal.(LiteralString)
	if ok {
		return caseString(litString.Value)
	}
	litBool, ok := literal.(LiteralBool)
	if ok {
		return caseBool(litBool.Value)
	}

	var nilCase Result
	return nilCase
}
