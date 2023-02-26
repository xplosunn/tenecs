package types

import (
	"fmt"
)

/*
There are different categories of types we care about:
1. Functions -> can construct types
2. TypeArgument -> can only happen when there's an unresolved generic in scope
3. "concrete" types
*/

type VariableType interface {
	sealedVariableType()
	VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void)
}

type StructFieldVariableType interface {
	sealedStructFieldVariableType()
	StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void)
}

func StructFieldVariableTypeFromVariableType(varType VariableType) (StructFieldVariableType, bool) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return caseTypeArgument, true
	} else if caseStruct != nil {
		return caseStruct, true
	} else if caseInterface != nil {
		return nil, false
	} else if caseFunction != nil {
		return nil, false
	} else if caseBasicType != nil {
		return caseBasicType, true
	} else if caseVoid != nil {
		return caseBasicType, true
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func VariableTypeFromStructFieldVariableType(structVarType StructFieldVariableType) VariableType {
	caseTypeArgument, caseStruct, caseBasicType, caseVoid := structVarType.StructFieldVariableTypeCases()
	if caseTypeArgument != nil {
		return caseTypeArgument
	} else if caseStruct != nil {
		return caseStruct
	} else if caseBasicType != nil {
		return caseBasicType
	} else if caseVoid != nil {
		return caseVoid
	} else {
		panic(fmt.Errorf("cases on %v", structVarType))
	}
}

type TypeArgument struct {
	Name string
}

func (t *TypeArgument) sealedVariableType() {}
func (t *TypeArgument) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return t, nil, nil, nil, nil, nil
}
func (t *TypeArgument) sealedStructFieldVariableType() {}
func (t *TypeArgument) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return t, nil, nil, nil
}

type Struct struct {
	Package string
	Name    string
	Fields  map[string]StructFieldVariableType
}

func (s *Struct) sealedVariableType() {}
func (s *Struct) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, s, nil, nil, nil, nil
}
func (s *Struct) sealedStructFieldVariableType() {}
func (s *Struct) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, s, nil, nil
}
func (s *Struct) sealedConstructableVariableTypeVariableType() {}
func (s *Struct) ConstructableVariableTypeCases() (*Struct, *Interface) {
	return s, nil
}

type Interface struct {
	Package   string
	Name      string
	Variables map[string]VariableType
}

func (i *Interface) sealedVariableType() {}
func (i *Interface) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, i, nil, nil, nil
}
func (i *Interface) sealedConstructableVariableTypeVariableType() {}
func (i *Interface) ConstructableVariableTypeCases() (*Struct, *Interface) {
	return nil, i
}

type Function struct {
	Generics   []string
	Arguments  []FunctionArgument
	ReturnType VariableType
}

func (f *Function) sealedVariableType() {}
func (f *Function) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, f, nil, nil
}

type FunctionArgument struct {
	Name         string
	VariableType VariableType
}

type BasicType struct {
	Type string
}

func (b *BasicType) sealedVariableType() {}
func (b *BasicType) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, b, nil
}
func (b *BasicType) sealedStructFieldVariableType() {}
func (b *BasicType) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, nil, b, nil
}

type Void struct {
}

func (v *Void) sealedVariableType() {}
func (v *Void) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, nil, v
}
func (v *Void) sealedStructFieldVariableType() {}
func (v *Void) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, nil, nil, v
}
