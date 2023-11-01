package parser

import (
	"github.com/alecthomas/participle/v2"
)

var literalUnion = participle.Union[Literal](LiteralFloat{}, LiteralInt{}, LiteralString{}, LiteralBool{}, LiteralNull{})

type Literal interface {
	sealedLiteral()
}

type LiteralFloat struct {
	Value float64 `@Float`
}

func (literal LiteralFloat) sealedLiteral() {}

type LiteralInt struct {
	Negative bool `@"-"?`
	Value    int  `@Int`
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

type LiteralNull struct {
	Value bool `@"null"`
}

func (literal LiteralNull) sealedLiteral() {}

func LiteralExhaustiveSwitch(
	literal Literal,
	caseFloat func(literal float64),
	caseInt func(literal int),
	caseString func(literal string),
	caseBool func(literal bool),
	caseNull func(),
) {
	litFloat, ok := literal.(LiteralFloat)
	if ok {
		caseFloat(litFloat.Value)
		return
	}
	litInt, ok := literal.(LiteralInt)
	if ok {
		value := litInt.Value
		if litInt.Negative {
			value = 0 - value
		}
		caseInt(value)
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
	_, ok = literal.(LiteralNull)
	if ok {
		caseNull()
		return
	}
}
