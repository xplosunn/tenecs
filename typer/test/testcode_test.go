package parser_typer_test

// ###############################################
// # This file is generated via code-generation. #
// # Check gen_test.go                           #
// ###############################################

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestArrayOfArray(t *testing.T) {
	validProgram(t, testcode.ArrayOfArray)
}

func TestArrayVariableWithEmptyArray(t *testing.T) {
	validProgram(t, testcode.ArrayVariableWithEmptyArray)
}

func TestArrayVariableWithTwoElementArray(t *testing.T) {
	validProgram(t, testcode.ArrayVariableWithTwoElementArray)
}

func TestBasicTypeFalse(t *testing.T) {
	validProgram(t, testcode.BasicTypeFalse)
}

func TestBasicTypeTrue(t *testing.T) {
	validProgram(t, testcode.BasicTypeTrue)
}

func TestGenericFunctionDeclared(t *testing.T) {
	validProgram(t, testcode.GenericFunctionDeclared)
}

func TestGenericFunctionDoubleInvoked(t *testing.T) {
	validProgram(t, testcode.GenericFunctionDoubleInvoked)
}

func TestGenericFunctionFixingArray(t *testing.T) {
	validProgram(t, testcode.GenericFunctionFixingArray)
}

func TestGenericFunctionInvoked1(t *testing.T) {
	validProgram(t, testcode.GenericFunctionInvoked1)
}

func TestGenericFunctionInvoked2(t *testing.T) {
	validProgram(t, testcode.GenericFunctionInvoked2)
}

func TestGenericFunctionInvoked3(t *testing.T) {
	validProgram(t, testcode.GenericFunctionInvoked3)
}

func TestGenericFunctionInvoked4(t *testing.T) {
	validProgram(t, testcode.GenericFunctionInvoked4)
}

func TestGenericFunctionSingleElementArray(t *testing.T) {
	validProgram(t, testcode.GenericFunctionSingleElementArray)
}

func TestGenericFunctionTakingArray(t *testing.T) {
	validProgram(t, testcode.GenericFunctionTakingArray)
}

func TestGenericImplementedInterfaceFunctionAllAnnotated(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionAllAnnotated)
}

func TestGenericImplementedInterfaceFunctionAnnotatedArg(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionAnnotatedArg)
}

func TestGenericImplementedInterfaceFunctionAnnotatedReturnType(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionAnnotatedReturnType)
}

func TestGenericImplementedInterfaceFunctionNotAnnotated(t *testing.T) {
	validProgram(t, testcode.GenericImplementedInterfaceFunctionNotAnnotated)
}

func TestGenericInterfaceFunction(t *testing.T) {
	validProgram(t, testcode.GenericInterfaceFunction)
}

func TestGenericStruct(t *testing.T) {
	validProgram(t, testcode.GenericStruct)
}

func TestGenericStructFunction(t *testing.T) {
	validProgram(t, testcode.GenericStructFunction)
}

func TestGenericStructInstance(t *testing.T) {
	validProgram(t, testcode.GenericStructInstance)
}

func TestInterfaceEmpty(t *testing.T) {
	validProgram(t, testcode.InterfaceEmpty)
}

func TestInterfaceReturningAnotherInterfaceInVariable(t *testing.T) {
	validProgram(t, testcode.InterfaceReturningAnotherInterfaceInVariable)
}

func TestInterfaceVariableFunctionOneArg(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableFunctionOneArg)
}

func TestInterfaceVariableFunctionTwoArgs(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableFunctionTwoArgs)
}

func TestInterfaceVariableFunctionZeroArgs(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableFunctionZeroArgs)
}

func TestInterfaceVariableString(t *testing.T) {
	validProgram(t, testcode.InterfaceVariableString)
}

func TestInterfaceWithSeparateModuleEmpty1(t *testing.T) {
	validProgram(t, testcode.InterfaceWithSeparateModuleEmpty1)
}

func TestInterfaceWithSeparateModuleEmpty2(t *testing.T) {
	validProgram(t, testcode.InterfaceWithSeparateModuleEmpty2)
}

func TestInterfaceWithSeparateModuleVariableString(t *testing.T) {
	validProgram(t, testcode.InterfaceWithSeparateModuleVariableString)
}

func TestMainProgramWithAnotherFunctionTakingConsole(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingConsole)
}

func TestMainProgramWithAnotherFunctionTakingConsoleAndMessage(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingConsoleAndMessage)
}

func TestMainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingConsoleAndMessageFromAnotherFunction)
}

func TestMainProgramWithAnotherFunctionTakingRuntime(t *testing.T) {
	validProgram(t, testcode.MainProgramWithAnotherFunctionTakingRuntime)
}

func TestMainProgramWithArgAnnotatedArg(t *testing.T) {
	validProgram(t, testcode.MainProgramWithArgAnnotatedArg)
}

func TestMainProgramWithArgAnnotatedArgAndReturn(t *testing.T) {
	validProgram(t, testcode.MainProgramWithArgAnnotatedArgAndReturn)
}

func TestMainProgramWithArgAnnotatedReturn(t *testing.T) {
	validProgram(t, testcode.MainProgramWithArgAnnotatedReturn)
}

func TestMainProgramWithIf(t *testing.T) {
	validProgram(t, testcode.MainProgramWithIf)
}

func TestMainProgramWithIfElse(t *testing.T) {
	validProgram(t, testcode.MainProgramWithIfElse)
}

func TestMainProgramWithInnerFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithInnerFunction)
}

func TestMainProgramWithSingleExpression(t *testing.T) {
	validProgram(t, testcode.MainProgramWithSingleExpression)
}

func TestMainProgramWithVariableWithFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunction)
}

func TestMainProgramWithVariableWithFunctionTakingFunction(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionTakingFunction)
}

func TestMainProgramWithVariableWithFunctionTakingFunctionFromStdLib1(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionTakingFunctionFromStdLib1)
}

func TestMainProgramWithVariableWithFunctionTakingFunctionFromStdLib2(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionTakingFunctionFromStdLib2)
}

func TestMainProgramWithVariableWithFunctionWithTypeInferred(t *testing.T) {
	validProgram(t, testcode.MainProgramWithVariableWithFunctionWithTypeInferred)
}

func TestModuleCreation1(t *testing.T) {
	validProgram(t, testcode.ModuleCreation1)
}

func TestModuleCreation2(t *testing.T) {
	validProgram(t, testcode.ModuleCreation2)
}

func TestModuleCreation3(t *testing.T) {
	validProgram(t, testcode.ModuleCreation3)
}

func TestModuleSelfCreation(t *testing.T) {
	validProgram(t, testcode.ModuleSelfCreation)
}

func TestModuleWithConstructorEmpty(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorEmpty)
}

func TestModuleWithConstructorWithArgUnused(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorWithArgUnused)
}

func TestModuleWithConstructorWithArgUsed(t *testing.T) {
	validProgram(t, testcode.ModuleWithConstructorWithArgUsed)
}

func TestNullFunction(t *testing.T) {
	validProgram(t, testcode.NullFunction)
}

func TestNullValue(t *testing.T) {
	validProgram(t, testcode.NullValue)
}

func TestOrArrayFunction(t *testing.T) {
	validProgram(t, testcode.OrArrayFunction)
}

func TestOrFunction(t *testing.T) {
	validProgram(t, testcode.OrFunction)
}

func TestOrVariableWithEmptyArray(t *testing.T) {
	validProgram(t, testcode.OrVariableWithEmptyArray)
}

func TestOrVariableWithTwoElementArray(t *testing.T) {
	validProgram(t, testcode.OrVariableWithTwoElementArray)
}

func TestStructAsVariable(t *testing.T) {
	validProgram(t, testcode.StructAsVariable)
}

func TestStructFunctionAccess(t *testing.T) {
	validProgram(t, testcode.StructFunctionAccess)
}

func TestStructVariableAccess(t *testing.T) {
	validProgram(t, testcode.StructVariableAccess)
}

func TestStructWithConstructorAnotherStruct1(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorAnotherStruct1)
}

func TestStructWithConstructorAnotherStruct2(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorAnotherStruct2)
}

func TestStructWithConstructorEmpty(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorEmpty)
}

func TestStructWithConstructorWithBooleans(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorWithBooleans)
}

func TestStructWithConstructorWithString(t *testing.T) {
	validProgram(t, testcode.StructWithConstructorWithString)
}

func TestTestsUnit(t *testing.T) {
	validProgram(t, testcode.TestsUnit)
}

func TestWhenExplicitExhaustive(t *testing.T) {
	validProgram(t, testcode.WhenExplicitExhaustive)
}

func TestWhenOtherMultipleTypes(t *testing.T) {
	validProgram(t, testcode.WhenOtherMultipleTypes)
}

func TestWhenOtherSingleType(t *testing.T) {
	validProgram(t, testcode.WhenOtherSingleType)
}

