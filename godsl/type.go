package godsl

import "fmt"

type Type interface {
	sealedType()
	typeToString() string
}

type goType struct {
	toString string
}

func (g goType) sealedType() {}

func (g goType) typeToString() string {
	return g.toString
}

func TypeAnyList() Type {
	return goType{"[]any"}
}

func TypeList(of Type) Type {
	return goType{"[]" + of.typeToString()}
}

func TypeStructOrInterface() Type {
	return goType{"map[string]any"}
}

func TypeFunc(numberOfArguments int) Type {
	if numberOfArguments < 0 {
		panic(fmt.Sprintf("tried to create TypeFunc with numberOfArguments = %d", numberOfArguments))
	}
	toString := "func("
	for i := 0; i < numberOfArguments; i++ {
		if i > 0 {
			toString += ","
		}
		toString += "any"
	}
	toString += ")any"
	return goType{toString}
}
