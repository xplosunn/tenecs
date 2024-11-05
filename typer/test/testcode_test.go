package parser_typer_test

// ###############################################
// # This file is generated via code-generation. #
// # Check typer/test_generate/main.go           #
// ###############################################

import (
	"github.com/xplosunn/tenecs/testcode"
	"testing"
)

func TestArrowInvocationOneArg(t *testing.T) {
	validProgram(t, testcode.ArrowInvocationOneArg)
}

func TestArrowInvocationOneArgChain(t *testing.T) {
	validProgram(t, testcode.ArrowInvocationOneArgChain)
}

func TestArrowInvocationThreeArg(t *testing.T) {
	validProgram(t, testcode.ArrowInvocationThreeArg)
}

func TestArrowInvocationTwoArg(t *testing.T) {
	validProgram(t, testcode.ArrowInvocationTwoArg)
}

func TestBasicTypeFalse(t *testing.T) {
	validProgram(t, testcode.BasicTypeFalse)
}

func TestBasicTypeTrue(t *testing.T) {
	validProgram(t, testcode.BasicTypeTrue)
}

func TestCommentEverywhere(t *testing.T) {
	validProgram(t, testcode.CommentEverywhere)
}

func TestFunctionsCallAndThenCall(t *testing.T) {
	validProgram(t, testcode.FunctionsCallAndThenCall)
}

func TestFunctionsNamedArg(t *testing.T) {
	validProgram(t, testcode.FunctionsNamedArg)
}

func TestGenericFromJson(t *testing.T) {
	validProgram(t, testcode.GenericFromJson)
}

func TestGenericFunctionDeclared(t *testing.T) {
	validProgram(t, testcode.GenericFunctionDeclared)
}

func TestGenericFunctionDoubleInvoked(t *testing.T) {
	validProgram(t, testcode.GenericFunctionDoubleInvoked)
}

func TestGenericFunctionFixingList(t *testing.T) {
	validProgram(t, testcode.GenericFunctionFixingList)
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

func TestGenericFunctionSingleElementList(t *testing.T) {
	validProgram(t, testcode.GenericFunctionSingleElementList)
}

func TestGenericFunctionTakingList(t *testing.T) {
	validProgram(t, testcode.GenericFunctionTakingList)
}

func TestGenericIO(t *testing.T) {
	validProgram(t, testcode.GenericIO)
}

func TestGenericImplementedStructFunctionAllAnnotated(t *testing.T) {
	validProgram(t, testcode.GenericImplementedStructFunctionAllAnnotated)
}

func TestGenericImplementedStructFunctionAnnotatedArg(t *testing.T) {
	validProgram(t, testcode.GenericImplementedStructFunctionAnnotatedArg)
}

func TestGenericImplementedStructFunctionAnnotatedReturnType(t *testing.T) {
	validProgram(t, testcode.GenericImplementedStructFunctionAnnotatedReturnType)
}

func TestGenericImplementedStructFunctionNotAnnotated(t *testing.T) {
	validProgram(t, testcode.GenericImplementedStructFunctionNotAnnotated)
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

func TestGenericsInferFunctionResult(t *testing.T) {
	validProgram(t, testcode.GenericsInferFunctionResult)
}

func TestGenericsInferHigherOrderFunction(t *testing.T) {
	validProgram(t, testcode.GenericsInferHigherOrderFunction)
}

func TestGenericsInferHigherOrderFunctionOr(t *testing.T) {
	validProgram(t, testcode.GenericsInferHigherOrderFunctionOr)
}

func TestGenericsInferHigherOrderFunctionOr2(t *testing.T) {
	validProgram(t, testcode.GenericsInferHigherOrderFunctionOr2)
}

func TestGenericsInferIdentity(t *testing.T) {
	validProgram(t, testcode.GenericsInferIdentity)
}

func TestGenericsInferList(t *testing.T) {
	validProgram(t, testcode.GenericsInferList)
}

func TestGenericsInferOrSecondArgument(t *testing.T) {
	validProgram(t, testcode.GenericsInferOrSecondArgument)
}

func TestGenericsInferTypeParameter(t *testing.T) {
	validProgram(t, testcode.GenericsInferTypeParameter)
}

func TestGenericsInferTypeParameterPartialLeft(t *testing.T) {
	validProgram(t, testcode.GenericsInferTypeParameterPartialLeft)
}

func TestImportAliasMain(t *testing.T) {
	validProgram(t, testcode.ImportAliasMain)
}

func TestListOfList(t *testing.T) {
	validProgram(t, testcode.ListOfList)
}

func TestListVariableWithEmptyList(t *testing.T) {
	validProgram(t, testcode.ListVariableWithEmptyList)
}

func TestListVariableWithTwoElementList(t *testing.T) {
	validProgram(t, testcode.ListVariableWithTwoElementList)
}

func TestMainProgramAnnotatedType(t *testing.T) {
	validProgram(t, testcode.MainProgramAnnotatedType)
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

func TestMainProgramWithIfElseIf(t *testing.T) {
	validProgram(t, testcode.MainProgramWithIfElseIf)
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

func TestNullFunction(t *testing.T) {
	validProgram(t, testcode.NullFunction)
}

func TestNullValue(t *testing.T) {
	validProgram(t, testcode.NullValue)
}

func TestOrFunction(t *testing.T) {
	validProgram(t, testcode.OrFunction)
}

func TestOrListFunction(t *testing.T) {
	validProgram(t, testcode.OrListFunction)
}

func TestOrVariableWithEmptyList(t *testing.T) {
	validProgram(t, testcode.OrVariableWithEmptyList)
}

func TestOrVariableWithTwoElementList(t *testing.T) {
	validProgram(t, testcode.OrVariableWithTwoElementList)
}

func TestRecursionFactorial(t *testing.T) {
	validProgram(t, testcode.RecursionFactorial)
}

func TestShortCircuitExplicit(t *testing.T) {
	validProgram(t, testcode.ShortCircuitExplicit)
}

func TestShortCircuitInferLeft(t *testing.T) {
	validProgram(t, testcode.ShortCircuitInferLeft)
}

func TestShortCircuitInferRight(t *testing.T) {
	validProgram(t, testcode.ShortCircuitInferRight)
}

func TestShortCircuitTwice(t *testing.T) {
	validProgram(t, testcode.ShortCircuitTwice)
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

func TestTypealiasGeneric(t *testing.T) {
	validProgram(t, testcode.TypealiasGeneric)
}

func TestTypealiasGenericOr(t *testing.T) {
	validProgram(t, testcode.TypealiasGenericOr)
}

func TestTypealiasGenericOrUsed(t *testing.T) {
	validProgram(t, testcode.TypealiasGenericOrUsed)
}

func TestTypealiasGenericUsed(t *testing.T) {
	validProgram(t, testcode.TypealiasGenericUsed)
}

func TestTypealiasGenericUsedGeneric(t *testing.T) {
	validProgram(t, testcode.TypealiasGenericUsedGeneric)
}

func TestTypealiasNested(t *testing.T) {
	validProgram(t, testcode.TypealiasNested)
}

func TestTypealiasSimple(t *testing.T) {
	validProgram(t, testcode.TypealiasSimple)
}

func TestTypealiasSimpleOr(t *testing.T) {
	validProgram(t, testcode.TypealiasSimpleOr)
}

func TestTypealiasSimpleOrUsed(t *testing.T) {
	validProgram(t, testcode.TypealiasSimpleOrUsed)
}

func TestTypealiasSimpleUsed(t *testing.T) {
	validProgram(t, testcode.TypealiasSimpleUsed)
}

func TestWhenAnnotatedVariable(t *testing.T) {
	validProgram(t, testcode.WhenAnnotatedVariable)
}

func TestWhenExplicitExhaustive(t *testing.T) {
	validProgram(t, testcode.WhenExplicitExhaustive)
}

func TestWhenGeneric(t *testing.T) {
	validProgram(t, testcode.WhenGeneric)
}

func TestWhenOtherMultipleTypes(t *testing.T) {
	validProgram(t, testcode.WhenOtherMultipleTypes)
}

func TestWhenOtherSingleType(t *testing.T) {
	validProgram(t, testcode.WhenOtherSingleType)
}

func TestWhenStruct(t *testing.T) {
	validProgram(t, testcode.WhenStruct)
}

