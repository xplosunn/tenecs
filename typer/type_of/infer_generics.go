package type_of

import (
	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

func AttemptGenericInference(node desugar.Node, function *types.Function, argumentsPassed []desugar.NamedArgument, genericsPassed []desugar.TypeAnnotation, expectedReturnType *types.VariableType, file string, scope binding.Scope) ([]types.VariableType, *type_error.TypecheckError) {
	resolvedGenerics := []types.VariableType{}
	for genericIndex, functionGenericName := range function.Generics {
		if len(genericsPassed) > 0 {
			shouldInfer := false
			passed := genericsPassed[genericIndex]
			for _, element := range passed.OrTypes {
				var err *type_error.TypecheckError
				desugar.TypeAnnotationElementExhaustiveSwitch(
					element,
					func(underscoreTypeAnnotation desugar.SingleNameType) {
						if len(passed.OrTypes) > 1 {
							err = type_error.PtrOnNodef(file, underscoreTypeAnnotation.Node, "Cannot infer part of an or type")
							return
						}
						shouldInfer = true
					},
					func(typeAnnotation desugar.SingleNameType) {},
					func(typeAnnotation desugar.FunctionType) {},
				)
				if err != nil {
					return nil, err
				}
			}
			if !shouldInfer {
				varType, err := scopecheck.ValidateTypeAnnotationInScope(passed, file, scope)
				if err != nil {
					return nil, type_error.FromScopeCheckError(file, err)
				}
				resolvedGenerics = append(resolvedGenerics, varType)
				continue
			}
		}

		if len(function.Arguments) != len(argumentsPassed) {
			return nil, type_error.PtrOnNodef(file, node, "expected %d arguments but got %d", len(function.Arguments), len(argumentsPassed))
		}
		var found types.VariableType
		for i, arg := range argumentsPassed {
			var typeOfArgFunction types.VariableType
			_, _, _, caseParameterFunction, _ := function.Arguments[i].VariableType.VariableTypeCases()
			if caseParameterFunction != nil {
				if len(arg.Argument.AccessOrInvocationChain) == 0 {
					lambdaOrList, ok := arg.Argument.Expression.(desugar.LambdaOrList)
					if ok && lambdaOrList.Lambda != nil {
						lambda := *lambdaOrList.Lambda
						if lambdaOrList.Generics == nil {
							argumentTypes, ok, err := tryToDetermineFunctionArgumentTypes(resolvedGenerics, lambda, function, caseParameterFunction, file, scope)
							if err != nil {
								return nil, err
							}
							if !ok {
								continue
							}
							localScope := scope
							for i, argType := range argumentTypes {
								var err *binding.ResolutionError
								localScope, err = binding.CopyAddingLocalVariable(localScope, lambda.Signature.Parameters[i].Name, argType)
								if err != nil {
									return nil, type_error.FromResolutionError(file, lambda.Signature.Parameters[i].Name.Node, err)
								}
							}
							var returnType types.VariableType
							if lambda.Signature.ReturnType != nil {
								rType, err := scopecheck.ValidateTypeAnnotationInScope(*lambda.Signature.ReturnType, file, scope)
								if err != nil {
									return nil, type_error.FromScopeCheckError(file, err)
								}
								returnType = rType
							} else {
								rType, err := TypeOfBlock(lambda.Block, file, localScope)
								if err != nil {
									return nil, err
								}
								returnType = rType
							}
							arguments := []types.FunctionArgument{}
							for i, variableType := range argumentTypes {
								arguments = append(arguments, types.FunctionArgument{
									Name:         lambda.Signature.Parameters[i].Name.String,
									VariableType: variableType,
								})
							}
							typeOfArgFunction = &types.Function{
								Generics:   nil,
								Arguments:  arguments,
								ReturnType: returnType,
							}
						}
					}
				}
			}
			typeOfArg := typeOfArgFunction
			if typeOfArg == nil {
				typeOfThisArg, err := TypeOfExpressionBox(arg.Argument, file, scope)
				if err != nil {
					continue
				}
				typeOfArg = typeOfThisArg
			}
			maybeInferred, ok := tryToInferGeneric(functionGenericName, function.Arguments[i].VariableType, typeOfArg)
			if !ok {
				return nil, type_error.PtrOnNodef(file, node, "Could not infer generics, please annotate them")
			}
			if maybeInferred != nil {
				if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
					found = maybeInferred
				} else {
					return nil, type_error.PtrOnNodef(file, node, "Could not infer generics, please annotate them")
				}
			}
		}
		if found == nil && expectedReturnType != nil {
			caseTypeArgument, _, _, _, _ := function.ReturnType.VariableTypeCases()
			if caseTypeArgument != nil && caseTypeArgument.Name == functionGenericName {
				found = *expectedReturnType
			}
		}
		if found == nil {
			return nil, type_error.PtrOnNodef(file, node, "Could not infer generics, please annotate them")
		}
		resolvedGenerics = append(resolvedGenerics, found)
	}
	if len(resolvedGenerics) == len(function.Generics) {
		return resolvedGenerics, nil
	} else {
		return nil, type_error.PtrOnNodef(file, node, "Could not infer generics, please annotate them")
	}
}

func tryToDetermineFunctionArgumentTypes(
	resolvedGenerics []types.VariableType,
	lambda desugar.Lambda,
	function *types.Function,
	caseParameterFunction *types.Function,
	file string,
	scope binding.Scope,
) ([]types.VariableType, bool, *type_error.TypecheckError) {
	arguments := []types.VariableType{}
	for i, parameter := range lambda.Signature.Parameters {
		if parameter.Type == nil {
			typeOfParam, ok := tryToDetermineFunctionArgumentType(resolvedGenerics, function.Generics, caseParameterFunction.Arguments[i].VariableType)
			if !ok {
				return nil, false, nil
			}
			arguments = append(arguments, typeOfParam)
		} else {
			typeOfParam, err := scopecheck.ValidateTypeAnnotationInScope(*parameter.Type, file, scope)
			if err != nil {
				return nil, false, type_error.FromScopeCheckError(file, err)
			}
			arguments = append(arguments, typeOfParam)
		}
	}
	return arguments, true, nil
}

func tryToDetermineFunctionArgumentType(
	resolvedGenerics []types.VariableType,
	functionGenerics []string,
	argumentVariableType types.VariableType,
) (types.VariableType, bool) {
	caseTypeArg, caseList, caseKnownType, _, _ := argumentVariableType.VariableTypeCases()
	if caseTypeArg != nil {
		for i, generic := range functionGenerics {
			if generic == caseTypeArg.Name {
				if len(resolvedGenerics) > i {
					return resolvedGenerics[i], true
				}
			}
		}
		return nil, false
	} else if caseList != nil {
		return caseList, true
	} else if caseKnownType != nil {
		return caseKnownType, true
	} else {
		return nil, false
	}
}

func tryToInferGeneric(genericName string, functionVarType types.VariableType, argVarType types.VariableType) (types.VariableType, bool) {
	funcCaseTypeArgument, funcCaseList, funcCaseKnownType, funcCaseFunction, funcCaseOr := functionVarType.VariableTypeCases()
	if funcCaseTypeArgument != nil {
		if funcCaseTypeArgument.Name == genericName {
			return argVarType, true
		}
		return nil, true
	} else if funcCaseList != nil {
		argList, ok := argVarType.(*types.List)
		if ok {
			inferred, ok := tryToInferGeneric(genericName, funcCaseList.Generic, argList.Generic)
			if inferred != nil || !ok {
				return inferred, ok
			}
		}
		return nil, true
	} else if funcCaseKnownType != nil {
		argKnownType, ok := argVarType.(*types.KnownType)
		if ok && len(funcCaseKnownType.Generics) == len(argKnownType.Generics) {
			for i, _ := range funcCaseKnownType.Generics {
				inferred, ok := tryToInferGeneric(genericName, funcCaseKnownType.Generics[i], argKnownType.Generics[i])
				if inferred != nil || !ok {
					return inferred, ok
				}
			}
		}
		return nil, true
	} else if funcCaseFunction != nil {
		for _, generic := range funcCaseFunction.Generics {
			if generic == genericName {
				return nil, true
			}
		}
		argFunction, ok := argVarType.(*types.Function)
		if !ok {
			return nil, false
		}
		if len(funcCaseFunction.Arguments) != len(argFunction.Arguments) {
			return nil, false
		}
		var found types.VariableType
		for i, _ := range funcCaseFunction.Arguments {
			maybeInferred, ok := tryToInferGeneric(genericName, funcCaseFunction.Arguments[i].VariableType, argFunction.Arguments[i].VariableType)
			if !ok {
				return nil, false
			}
			if maybeInferred != nil {
				if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
					found = maybeInferred
				} else {
					return nil, false
				}
			}
		}
		maybeInferred, ok := tryToInferGeneric(genericName, funcCaseFunction.ReturnType, argFunction.ReturnType)
		if !ok {
			return nil, false
		}
		if maybeInferred != nil {
			if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
				found = maybeInferred
			} else {
				return nil, false
			}
		}
		return found, true
	} else if funcCaseOr != nil {
		_, _, _, _, caseArgOr := argVarType.VariableTypeCases()
		if caseArgOr != nil {
			remainingTypesToMatch := []types.VariableType{}
			for _, argVarType := range caseArgOr.Elements {
				matched := false
				for _, element := range funcCaseOr.Elements {
					if types.VariableTypeEq(argVarType, element) {
						matched = true
						break
					}
				}
				if !matched {
					remainingTypesToMatch = append(remainingTypesToMatch, argVarType)
				}
			}
			argVarType = &types.OrVariableType{Elements: remainingTypesToMatch}
		}

		var found types.VariableType
		for _, element := range funcCaseOr.Elements {
			maybeInferred, ok := tryToInferGeneric(genericName, element, argVarType)
			if !ok {
				return nil, false
			}
			if maybeInferred != nil {
				if found == nil || types.VariableTypeContainedIn(found, maybeInferred) {
					found = maybeInferred
				} else {
					return nil, false
				}
			}
		}
		return found, true
	} else {
		return nil, true
	}
}
