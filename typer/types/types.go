package types

import (
	"fmt"
)

type VariableType interface {
	sealedVariableType()
	VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void)
}

type StructVariableType interface {
	sealedStructVariableType()
	StructVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void)
}

type ConstructableVariableType interface {
	sealedConstructableVariableTypeVariableType()
	ConstructableVariableTypeCases() (*Struct, *Interface)
}

func StructVariableTypeFromVariableType(varType VariableType) (StructVariableType, bool) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return caseTypeArgument, true
	} else if caseStruct != nil {
		return *caseStruct, true
	} else if caseInterface != nil {
		return nil, false
	} else if caseFunction != nil {
		return nil, false
	} else if caseBasicType != nil {
		return *caseBasicType, true
	} else if caseVoid != nil {
		return *caseBasicType, true
	} else {
		panic(fmt.Errorf("code on %v", varType))
	}
}

func VariableTypeFromStructVariableType(structVarType StructVariableType) VariableType {
	caseTypeArgument, caseStruct, caseBasicType, caseVoid := structVarType.StructVariableTypeCases()
	if caseTypeArgument != nil {
		return *caseTypeArgument
	} else if caseStruct != nil {
		return *caseStruct
	} else if caseBasicType != nil {
		return *caseBasicType
	} else if caseVoid != nil {
		return *caseVoid
	} else {
		panic(fmt.Errorf("code on %v", structVarType))
	}
}

func VariableTypeFromConstructableVariableType(constructableVariableType ConstructableVariableType) VariableType {
	caseStruct, caseInterface := constructableVariableType.ConstructableVariableTypeCases()
	if caseStruct != nil {
		return *caseStruct
	} else if caseInterface != nil {
		return *caseInterface
	} else {
		panic(fmt.Errorf("code on %v", constructableVariableType))
	}
}

type TypeArgument struct {
	Name string
}

func (t TypeArgument) sealedVariableType() {}
func (t TypeArgument) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return &t, nil, nil, nil, nil, nil
}
func (t TypeArgument) sealedStructVariableType() {}
func (t TypeArgument) StructVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return &t, nil, nil, nil
}

type ResolvedTypeArgument struct {
	Name               string
	StructVariableType StructVariableType
}

type Struct struct {
	Package               string
	ResolvedTypeArguments []ResolvedTypeArgument
	Name                  string
	Fields                map[string]StructVariableType
}

func (s Struct) sealedVariableType() {}
func (s Struct) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, &s, nil, nil, nil, nil
}
func (s Struct) sealedStructVariableType() {}
func (s Struct) StructVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, &s, nil, nil
}
func (s Struct) sealedConstructableVariableTypeVariableType() {}
func (s Struct) ConstructableVariableTypeCases() (*Struct, *Interface) {
	return &s, nil
}

type Interface struct {
	Package string
	Name    string
}

func (i Interface) sealedVariableType() {}
func (i Interface) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, &i, nil, nil, nil
}
func (i Interface) sealedConstructableVariableTypeVariableType() {}
func (i Interface) ConstructableVariableTypeCases() (*Struct, *Interface) {
	return nil, &i
}

type Function struct {
	Generics   []string
	Arguments  []FunctionArgument
	ReturnType VariableType
}

func (f Function) sealedVariableType() {}
func (f Function) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, &f, nil, nil
}

type FunctionArgument struct {
	Name         string
	VariableType VariableType
}

type BasicType struct {
	Type string
}

func (b BasicType) sealedVariableType() {}
func (b BasicType) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, &b, nil
}
func (b BasicType) sealedStructVariableType() {}
func (b BasicType) StructVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, nil, &b, nil
}

type Void struct {
}

func (v Void) sealedVariableType() {}
func (v Void) VariableTypeCases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, nil, &v
}
func (v Void) sealedStructVariableType() {}
func (v Void) StructVariableTypeCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, nil, nil, &v
}
