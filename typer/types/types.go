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
	VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType)
}

type StructFieldVariableType interface {
	sealedStructFieldVariableType()
	StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType)
}

func StructFieldVariableTypeFromVariableType(varType VariableType) (StructFieldVariableType, bool) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid, caseArray, caseOr := varType.VariableTypeCases()
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
	} else if caseArray != nil {
		return caseArray, true
	} else if caseOr != nil {
		elements := []StructFieldVariableType{}
		for _, element := range caseOr.Elements {
			newElement, ok := StructFieldVariableTypeFromVariableType(element)
			if !ok {
				return nil, false
			}
			elements = append(elements, newElement)
		}
		return &OrStructFieldVariableType{
			Elements: elements,
		}, true
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func VariableTypeFromStructFieldVariableType(structVarType StructFieldVariableType) VariableType {
	caseTypeArgument, caseStruct, caseBasicType, caseVoid, caseArray, caseOr := structVarType.StructFieldVariableTypeCases()
	if caseTypeArgument != nil {
		return caseTypeArgument
	} else if caseStruct != nil {
		return caseStruct
	} else if caseBasicType != nil {
		return caseBasicType
	} else if caseVoid != nil {
		return caseVoid
	} else if caseArray != nil {
		return caseArray
	} else if caseOr != nil {
		elements := []VariableType{}
		for _, element := range caseOr.Elements {
			elements = append(elements, VariableTypeFromStructFieldVariableType(element))
		}
		return &OrVariableType{
			Elements: elements,
		}
	} else {
		panic(fmt.Errorf("cases on %v", structVarType))
	}
}

type OrVariableType struct {
	Elements []VariableType
}

func (o *OrVariableType) sealedVariableType() {}
func (o *OrVariableType) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, nil, nil, nil, nil, nil, nil, o
}

type OrStructFieldVariableType struct {
	Elements []StructFieldVariableType
}

func (o *OrStructFieldVariableType) sealedStructFieldVariableType() {}
func (o *OrStructFieldVariableType) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType) {
	return nil, nil, nil, nil, nil, o
}

type TypeArgument struct {
	Name string
}

func (t *TypeArgument) sealedVariableType() {}
func (t *TypeArgument) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return t, nil, nil, nil, nil, nil, nil, nil
}
func (t *TypeArgument) sealedStructFieldVariableType() {}
func (t *TypeArgument) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType) {
	return t, nil, nil, nil, nil, nil
}

type Struct struct {
	Package string
	Name    string
	Fields  map[string]StructFieldVariableType
}

func (s *Struct) sealedVariableType() {}
func (s *Struct) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, s, nil, nil, nil, nil, nil, nil
}
func (s *Struct) sealedStructFieldVariableType() {}
func (s *Struct) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType) {
	return nil, s, nil, nil, nil, nil
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
func (i *Interface) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, nil, i, nil, nil, nil, nil, nil
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
func (f *Function) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, nil, nil, f, nil, nil, nil, nil
}

type FunctionArgument struct {
	Name         string
	VariableType VariableType
}

type BasicType struct {
	Type string
}

func (b *BasicType) sealedVariableType() {}
func (b *BasicType) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, nil, nil, nil, b, nil, nil, nil
}
func (b *BasicType) sealedStructFieldVariableType() {}
func (b *BasicType) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType) {
	return nil, nil, b, nil, nil, nil
}

type Void struct {
}

func (v *Void) sealedVariableType() {}
func (v *Void) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, nil, nil, nil, nil, v, nil, nil
}
func (v *Void) sealedStructFieldVariableType() {}
func (v *Void) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType) {
	return nil, nil, nil, v, nil, nil
}

type Array struct {
	OfType StructFieldVariableType
}

func (a *Array) sealedVariableType() {}
func (a *Array) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void, *Array, *OrVariableType) {
	return nil, nil, nil, nil, nil, nil, a, nil
}
func (a *Array) sealedStructFieldVariableType() {}
func (a *Array) StructFieldVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void, *Array, *OrStructFieldVariableType) {
	return nil, nil, nil, nil, a, nil
}
