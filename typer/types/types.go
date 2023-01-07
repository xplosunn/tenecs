package types

import (
	"fmt"
)

type VariableType interface {
	sealedVariableType()
	Cases() (*Struct, *Interface, *Function, *BasicType, *Void)
}

type StructVariableType interface {
	sealedStructVariableType()
	StructCases() (*Struct, *BasicType, *Void)
}

type ConstructableVariableType interface {
	sealedConstructableVariableTypeVariableType()
	ConstructableCases() (*Struct, *Interface)
}

func StructVariableTypeFromVariableType(varType VariableType) (StructVariableType, bool) {
	caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
	if caseStruct != nil {
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
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func VariableTypeFromStructVariableType(structVarType StructVariableType) VariableType {
	caseStruct, caseBasicType, caseVoid := structVarType.StructCases()
	if caseStruct != nil {
		return *caseStruct
	} else if caseBasicType != nil {
		return *caseBasicType
	} else if caseVoid != nil {
		return *caseVoid
	} else {
		panic(fmt.Errorf("cases on %v", structVarType))
	}
}

func VariableTypeFromConstructableVariableType(constructableVariableType ConstructableVariableType) VariableType {
	caseStruct, caseInterface := constructableVariableType.ConstructableCases()
	if caseStruct != nil {
		return *caseStruct
	} else if caseInterface != nil {
		return *caseInterface
	} else {
		panic(fmt.Errorf("cases on %v", constructableVariableType))
	}
}

type Struct struct {
	Package string
	Name    string
}

func (s Struct) sealedVariableType() {}
func (s Struct) Cases() (*Struct, *Interface, *Function, *BasicType, *Void) {
	return &s, nil, nil, nil, nil
}
func (s Struct) sealedStructVariableType() {}
func (s Struct) StructCases() (*Struct, *BasicType, *Void) {
	return &s, nil, nil
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
func (i Interface) Cases() (*Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, &i, nil, nil, nil
}
func (i Interface) sealedConstructableVariableTypeVariableType() {}
func (i Interface) ConstructableCases() (*Struct, *Interface) {
	return nil, &i
}

type Function struct {
	Arguments  []FunctionArgument
	ReturnType VariableType
}

func (f Function) sealedVariableType() {}
func (f Function) Cases() (*Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, &f, nil, nil
}

type FunctionArgument struct {
	Name         string
	VariableType VariableType
}

type BasicType struct {
	Type string
}

func (b BasicType) sealedVariableType() {}
func (b BasicType) Cases() (*Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, &b, nil
}
func (b BasicType) sealedStructVariableType() {}
func (b BasicType) StructCases() (*Struct, *BasicType, *Void) {
	return nil, &b, nil
}

type Void struct {
}

func (v Void) sealedVariableType() {}
func (v Void) Cases() (*Struct, *Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, nil, &v
}
func (v Void) sealedStructVariableType() {}
func (v Void) StructCases() (*Struct, *BasicType, *Void) {
	return nil, nil, &v
}
