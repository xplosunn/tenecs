package parser

import (
	"github.com/alecthomas/participle/v2"
)

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
	Value bool `@"true"`
	False bool `| @"false"`
}

func (literal LiteralBool) sealedLiteral() {}

func LiteralExhaustiveSwitch(
	literal Literal,
	caseFloat func(literal float64),
	caseInt func(literal int),
	caseString func(literal string),
	caseBool func(literal bool),
) {
	litFloat, ok := literal.(LiteralFloat)
	if ok {
		caseFloat(litFloat.Value)
		return
	}
	litInt, ok := literal.(LiteralInt)
	if ok {
		caseInt(litInt.Value)
		return
	}
	litString, ok := literal.(LiteralString)
	if ok {
		caseString(litString.Value)
		return
	}
	litBool, ok := literal.(LiteralBool)
	if ok {
		caseBool(litBool.Value)
		return
	}
}
