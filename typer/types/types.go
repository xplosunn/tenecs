package types

import (
	"fmt"
)

type VariableType interface {
	sealedVariableType()
	Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void)
}

type StructVariableType interface {
	sealedStructVariableType()
	StructCases() (*TypeArgument, *Struct, *BasicType, *Void)
}

type ConstructableVariableType interface {
	sealedConstructableVariableTypeVariableType()
	ConstructableCases() (*Struct, *Interface)
}

func StructVariableTypeFromVariableType(varType VariableType) (StructVariableType, bool) {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
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
	caseTypeArgument, caseStruct, caseBasicType, caseVoid := structVarType.StructCases()
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
	caseStruct, caseInterface := constructableVariableType.ConstructableCases()
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
func (t TypeArgument) Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return &t, nil, nil, nil, nil, nil
}
func (t TypeArgument) sealedStructVariableType() {}
func (t TypeArgument) StructCases() (*TypeArgument, *Struct, *BasicType, *Void) {
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
}

func (s Struct) sealedVariableType() {}
func (s Struct) Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, &s, nil, nil, nil, nil
}
func (s Struct) sealedStructVariableType() {}
func (s Struct) StructCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, &s, nil, nil
}
func (s Struct) sealedConstructableVariableTypeVariableType() {}
func (s Struct) ConstructableCases() (*Struct, *Interface) {
	return &s, nil
}

type Interface struct {
	Package string
	Name    string
}

func (i Interface) sealedVariableType() {}
func (i Interface) Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, &i, nil, nil, nil
}
func (i Interface) sealedConstructableVariableTypeVariableType() {}
func (i Interface) ConstructableCases() (*Struct, *Interface) {
	return nil, &i
}

type Function struct {
	Generics   []string
	Arguments  []FunctionArgument
	ReturnType VariableType
}

func (f Function) sealedVariableType() {}
func (f Function) Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
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
func (b BasicType) Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, &b, nil
}
func (b BasicType) sealedStructVariableType() {}
func (b BasicType) StructCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, nil, &b, nil
}

type Void struct {
}

func (v Void) sealedVariableType() {}
func (v Void) Cases() (*TypeArgument, *Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, nil, &v
}
func (v Void) sealedStructVariableType() {}
func (v Void) StructCases() (*TypeArgument, *Struct, *BasicType, *Void) {
	return nil, nil, nil, &v
}
