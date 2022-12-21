package typer

import (
	"encoding/json"
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/parser"
	"reflect"
	"unicode"
)

func Typecheck(parsed parser.FileTopLevel) error {
	pkg, imports, modules := parser.FileTopLevelFields(parsed)
	err := validatePackage(pkg)
	if err != nil {
		return err
	}
	universe, err := resolveImports(imports, StdLib)
	if err != nil {
		return err
	}
	modulesMap, parserModulesMap, err := validateModulesImplements(modules, universe)
	if err != nil {
		return err
	}
	err = validateModulesVariableTypesAndExpressions(modulesMap, parserModulesMap, universe)
	if err != nil {
		return err
	}
	err = validateModulesVariableFunctionBlocks(modulesMap, parserModulesMap, universe)
	if err != nil {
		return err
	}

	return nil
}

type Universe struct {
	TypeByTypeName     immutable.Map[string, VariableType]
	TypeByVariableName immutable.Map[string, VariableType]
}

func NewUniverseFromDefaults() Universe {
	mapBuilder := immutable.NewMapBuilder[string, VariableType](nil)

	for key, value := range DefaultTypesAvailableWithoutImport {
		mapBuilder.Set(key, value)
	}
	return Universe{
		TypeByTypeName:     *mapBuilder.Map(),
		TypeByVariableName: *immutable.NewMap[string, VariableType](nil),
	}
}

func NewUniverseFromInterface(interf Interface) Universe {
	mapBuilder := immutable.NewMapBuilder[string, VariableType](nil)

	for key, value := range interf.Variables {
		mapBuilder.Set(key, value)
	}
	return Universe{
		TypeByTypeName:     *immutable.NewMap[string, VariableType](nil),
		TypeByVariableName: *mapBuilder.Map(),
	}
}

func copyUniverseAddingType(universe Universe, typeName string, varType VariableType) (Universe, *TypecheckError) {
	_, ok := universe.TypeByTypeName.Get(typeName)
	if ok {
		bytes, err := json.Marshal(universe.TypeByTypeName)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("type already exists %s in %s", typeName, string(bytes))
	}
	return Universe{
		TypeByTypeName:     *universe.TypeByTypeName.Set(typeName, varType),
		TypeByVariableName: universe.TypeByVariableName,
	}, nil
}

func copyUniverseAddingVariable(universe Universe, variableName string, varType VariableType) (Universe, *TypecheckError) {
	_, ok := universe.TypeByVariableName.Get(variableName)
	if ok {
		bytes, err := json.Marshal(universe.TypeByVariableName)
		if err != nil {
			panic(err)
		}
		return universe, PtrTypeCheckErrorf("variable already exists %s in %s", variableName, string(bytes))
	}
	return Universe{
		TypeByTypeName:     universe.TypeByTypeName,
		TypeByVariableName: *universe.TypeByVariableName.Set(variableName, varType),
	}, nil
}

func copyUniverseAddingVariables(universe Universe, variables map[string]VariableType) (Universe, *TypecheckError) {
	result := universe
	for name, varType := range variables {
		updatedResult, err := copyUniverseAddingVariable(result, name, varType)
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
}

func copyUniverseAddingFunctionArguments(universe Universe, functionArguments []FunctionArgument) (Universe, *TypecheckError) {
	result := universe
	for _, argument := range functionArguments {
		updatedResult, err := copyUniverseAddingVariable(result, argument.Name, argument.VariableType)
		if err != nil {
			return result, err
		}
		result = updatedResult
	}
	return result, nil
}

func validatePackage(node parser.Package) *TypecheckError {
	identifier := parser.PackageFields(node)
	for _, r := range identifier {
		if !unicode.IsLower(r) {
			return PtrTypeCheckErrorf("package name should start with a lowercase letter")
		} else {
			return nil
		}
	}
	return nil
}

func resolveImports(nodes []parser.Import, stdLib Package) (Universe, *TypecheckError) {
	universe := NewUniverseFromDefaults()
	for _, node := range nodes {
		dotSeparatedNames := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			return universe, PtrTypeCheckErrorf("all interfaces belong to a package")
		}
		currPackage := stdLib
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name]
				if !ok {
					return universe, PtrTypeCheckErrorf("no package " + name + " found")
				}
				currPackage = p
				continue
			}
			interf, ok := currPackage.Interfaces[name]
			if !ok {
				return universe, PtrTypeCheckErrorf("no interface " + name + " found")
			}
			updatedUniverse, err := copyUniverseAddingType(universe, name, interf)
			if err != nil {
				return universe, err
			}
			universe = updatedUniverse
		}
	}
	return universe, nil
}

func validateModulesImplements(nodes []parser.Module, universe Universe) (map[string]*Module, map[string]parser.Module, *TypecheckError) {
	modulesMap := map[string]*Module{}
	parserModulesMap := map[string]parser.Module{}
	for _, node := range nodes {
		name, implements, declarations := parser.ModuleFields(node)
		_ = declarations
		_, ok := modulesMap[name]
		if ok {
			return nil, nil, PtrTypeCheckErrorf("another module declared with name %s", name)
		}
		implementedInterfaces, err := validateImplementedInterfacesDoNotConflict(implements, universe)
		if err != nil {
			return nil, nil, err
		}
		modulesMap[name] = &Module{
			Name:       name,
			Implements: implementedInterfaces,
		}
		parserModulesMap[name] = node
	}
	return modulesMap, parserModulesMap, nil
}

func validateImplementedInterfacesDoNotConflict(implements []string, universe Universe) ([]Interface, *TypecheckError) {
	implementedInterfaces := []Interface{}
	for _, implement := range implements {
		varType, ok := universe.TypeByTypeName.Get(implement)
		if !ok {
			return implementedInterfaces, PtrTypeCheckErrorf("not found interface with name %s", implement)
		}
		caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
		if caseInterface != nil {
			implementedInterfaces = append(implementedInterfaces, *caseInterface)
		} else if caseFunction != nil {
			return implementedInterfaces, PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implement, printableName(varType))
		} else if caseBasicType != nil {
			return implementedInterfaces, PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implement, printableName(varType))
		} else if caseVoid != nil {
			return implementedInterfaces, PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implement, printableName(varType))
		} else {
			panic(fmt.Errorf("cases on %v", varType))
		}
	}
	allInterfaceVariableNames := map[string]string{}
	for _, implementedInterface := range implementedInterfaces {
		for varName, _ := range implementedInterface.Variables {
			conflictingInterfaceName, ok := allInterfaceVariableNames[varName]
			if ok {
				return nil, PtrTypeCheckErrorf("incompatible interfaces implemented because both shared a variable name '%s': %s, %s", varName, implementedInterface.Name, conflictingInterfaceName)
			}
			allInterfaceVariableNames[varName] = implementedInterface.Name
		}
	}
	return implementedInterfaces, nil
}

func validateModulesVariableTypesAndExpressions(modulesMap map[string]*Module, parserModulesMap map[string]parser.Module, universe Universe) *TypecheckError {
	for moduleName, parserModule := range parserModulesMap {
		for _, node := range parserModule.Declarations {
			var typeOfInterfaceVariableWithSameName *VariableType
		typeOfInterfaceVariableWithSameNameLoop:
			for _, implementedInterface := range modulesMap[moduleName].Implements {
				for varName, varType := range implementedInterface.Variables {
					if varName == node.Name {
						typeOfInterfaceVariableWithSameName = &varType
						break typeOfInterfaceVariableWithSameNameLoop
					}
				}
			}

			varType, err := validateVariableTypeAndExpression(node, typeOfInterfaceVariableWithSameName, universe)
			if err != nil {
				return err
			}
			if modulesMap[moduleName].Variables == nil {
				modulesMap[moduleName].Variables = map[string]VariableType{}
			}
			modulesMap[moduleName].Variables[node.Name] = varType
		}
	}
	return nil
}

func validateVariableTypeAndExpression(node parser.Declaration, typeOfInterfaceVariableWithSameName *VariableType, universe Universe) (VariableType, *TypecheckError) {
	if typeOfInterfaceVariableWithSameName == nil {
		return nonPublicDeclarationVariableType(node.Name, node.Expression, universe)
	}
	err := isExpressionOfExpectedType(node.Name, node.Expression, *typeOfInterfaceVariableWithSameName, universe)
	if err != nil {
		return nil, err
	}
	return *typeOfInterfaceVariableWithSameName, nil
}

func isExpressionOfExpectedType(variableName string, exp parser.Expression, expectedType VariableType, universe Universe) *TypecheckError {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda := exp.Cases()
	if caseLiteralExp != nil {
		return isLiteralOfExpectedType(caseLiteralExp.Literal, expectedType)
	} else if caseReferenceOrInvocation != nil {
		if caseReferenceOrInvocation.Arguments != nil {
			panic("not supported yet (caseReferenceOrInvocation.Arguments)")
		}
		if len(caseReferenceOrInvocation.DotSeparatedVars) > 1 {
			panic("not supported yet (caseReferenceOrInvocation.DotSeparatedVars)")
		}
		varName := caseReferenceOrInvocation.DotSeparatedVars[0]
		varType, ok := universe.TypeByVariableName.Get(varName)
		if !ok {
			return PtrTypeCheckErrorf("Not found reference: %s", variableName)
		}
		if !variableTypeEq(varType, expectedType) {
			return PtrTypeCheckErrorf("expected type %s but %s is %s", printableName(expectedType), varName, printableName(varType))
		}
		return nil
	} else if caseLambda != nil {
		return isLambdaSignatureOfExpectedType(*caseLambda, expectedType, universe)
	} else {
		panic(fmt.Errorf("cases on %v", exp))
	}
}

func isLambdaSignatureOfExpectedType(lambda parser.Lambda, expectedType VariableType, universe Universe) *TypecheckError {
	var expectedFunction Function
	caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.Cases()
	if caseInterface != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseFunction != nil {
		expectedFunction = *caseFunction
	} else if caseBasicType != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseVoid != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else {
		panic(fmt.Errorf("cases on %v", expectedType))
	}

	parameters, annotatedReturnType, block := parser.LambdaFields(lambda)
	_ = block
	if len(parameters) != len(expectedFunction.Arguments) {
		return PtrTypeCheckErrorf("expected same number of arguments as interface variable (%d) but found %d", len(expectedFunction.Arguments), len(parameters))
	}
	for i, parameter := range parameters {
		if parameter.Type == "" {
			continue
		}

		varType, ok := universe.TypeByTypeName.Get(parameter.Type)
		if !ok {
			return PtrTypeCheckErrorf("not found type: %s", parameter.Type)
		}

		if !variableTypeEq(varType, expectedFunction.Arguments[i].VariableType) {
			return PtrTypeCheckErrorf("in parameter position %d expected type %s but you have annotated %s", i, printableName(expectedFunction.Arguments[i].VariableType), parameter.Type)
		}
	}

	if annotatedReturnType == "" {
		return nil
	}
	varType, ok := universe.TypeByTypeName.Get(annotatedReturnType)
	if !ok {
		return PtrTypeCheckErrorf("not found type: %s", annotatedReturnType)
	}

	if !variableTypeEq(varType, expectedFunction.ReturnType) {
		return PtrTypeCheckErrorf("in return type expected type %s but you have annotated %s", printableName(expectedFunction.ReturnType), annotatedReturnType)
	}
	return nil
}

func variableTypeEq(v1 VariableType, v2 VariableType) bool {
	return reflect.DeepEqual(v1, v2)
}

func isLiteralOfExpectedType(argument parser.Literal, expectedType VariableType) *TypecheckError {
	caseInterface, caseFunction, caseBasicType, caseVoid := expectedType.Cases()
	if caseInterface != nil {
		return PtrTypeCheckErrorf("expected type %s but found an Inferface", printableName(expectedType))
	} else if caseFunction != nil {
		return PtrTypeCheckErrorf("expected type %s but found a Function", printableName(expectedType))
	} else if caseBasicType != nil {
		basicType := *caseBasicType
		expectBasicType := func(typeName string) *TypecheckError {
			if basicType.Type != typeName {
				return PtrTypeCheckErrorf("expected type %s but found %s", typeName, basicType.Type)
			}
			return nil
		}
		return parser.LiteralFold[*TypecheckError](
			argument,
			func(arg float64) *TypecheckError {
				return expectBasicType("Float")
			},
			func(arg int) *TypecheckError {
				return expectBasicType("Int")
			},
			func(arg string) *TypecheckError {
				return expectBasicType("String")
			},
			func(arg bool) *TypecheckError {
				return expectBasicType("Boolean")
			},
		)
	} else if caseVoid != nil {
		return PtrTypeCheckErrorf("expected type %s but found Void", printableName(expectedType))
	} else {
		panic(fmt.Errorf("cases on %v", expectedType))
	}
}

func nonPublicDeclarationVariableType(variableName string, expression parser.Expression, universe Universe) (VariableType, *TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda := expression.Cases()
	if caseLiteralExp != nil {
		return parser.LiteralFold(
			caseLiteralExp.Literal,
			func(arg float64) BasicType {
				return basicTypeFloat
			},
			func(arg int) BasicType {
				return basicTypeInt
			},
			func(arg string) BasicType {
				return basicTypeString
			},
			func(arg bool) BasicType {
				return basicTypeBoolean
			},
		), nil
	} else if caseReferenceOrInvocation != nil {
		return nil, PtrTypeCheckErrorf("references not supported on module variables (variable '%s')", variableName)
	} else if caseLambda != nil {
		function := Function{
			Arguments:  []FunctionArgument{},
			ReturnType: nil,
		}
		parameters, annotatedReturnType, block := parser.LambdaFields(*caseLambda)
		_ = block
		for _, parameter := range parameters {
			if parameter.Type == "" {
				return nil, PtrTypeCheckErrorf("parameter '%s' needs to be type annotated as the variable '%s' is not public", parameter.Name, variableName)
			}

			varType, ok := universe.TypeByTypeName.Get(parameter.Type)
			if !ok {
				return nil, PtrTypeCheckErrorf("not found type: %s", parameter.Type)
			}
			function.Arguments = append(function.Arguments, FunctionArgument{
				Name:         parameter.Name,
				VariableType: varType,
			})
		}
		if annotatedReturnType == "" {
			return nil, PtrTypeCheckErrorf("return type needs to be type annotated as the variable '%s' is not public", variableName)
		}
		varType, ok := universe.TypeByTypeName.Get(annotatedReturnType)
		if !ok {
			return nil, PtrTypeCheckErrorf("not found type: %s", annotatedReturnType)
		}
		function.ReturnType = varType
		return function, nil
	} else {
		panic(fmt.Errorf("cases on %v", expression))
	}
}

func validateModulesVariableFunctionBlocks(modulesMap map[string]*Module, parserModulesMap map[string]parser.Module, universe Universe) *TypecheckError {
	for moduleName, module := range modulesMap {
		for varName, varType := range module.Variables {
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			var function *Function
			if caseInterface != nil {
				continue
			} else if caseFunction != nil {
				function = caseFunction
			} else if caseBasicType != nil {
				continue
			} else if caseVoid != nil {
				continue
			} else {
				panic(fmt.Errorf("cases on %v", varType))
			}

			var parserLambda parser.Lambda
			foundParserLambda := false
			for _, declaration := range parserModulesMap[moduleName].Declarations {
				if declaration.Name == varName {
					caseLiteralExp, caseReferenceOrInvocation, caseLambda := declaration.Expression.Cases()
					if caseLiteralExp != nil {
						panic(fmt.Errorf("unexpected caseLiteralExp on %s.%s", moduleName, varName))
					} else if caseReferenceOrInvocation != nil {
						panic(fmt.Errorf("unexpected caseReferenceOrInvocation on %s.%s", moduleName, varName))
					} else if caseLambda != nil {
						parserLambda = *caseLambda
					} else {
						panic(fmt.Errorf("cases on %v", varType))
					}
					foundParserLambda = true
					break
				}
			}
			if !foundParserLambda {
				panic(fmt.Errorf("didn't foundParserLambda"))
			}

			blockUniverse, err := copyUniverseAddingFunctionArguments(universe, function.Arguments)
			if err != nil {
				return err
			}
			blockUniverse, err = copyUniverseAddingVariables(blockUniverse, module.Variables)
			if err != nil {
				return err
			}

			err = validateFunctionBlock(parserLambda.Block, blockUniverse)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func validateFunctionBlock(block []parser.ReferenceOrInvocation, universe Universe) *TypecheckError {
	for _, referenceOrInvocation := range block {
		dotSeparatedVarName, argumentsPtr := parser.ReferenceOrInvocationFields(referenceOrInvocation)
		if argumentsPtr == nil {
			panic("TODO")
		}
		arguments := *argumentsPtr

		currentUniverse := universe
		for i, varName := range dotSeparatedVarName {
			varType, ok := currentUniverse.TypeByVariableName.Get(varName)
			if !ok {
				return &TypecheckError{Message: "not found in scope: " + varName}
			}

			if i < len(dotSeparatedVarName)-1 {
				caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
				if caseInterface != nil {
					currentUniverse = NewUniverseFromInterface(*caseInterface)
				} else if caseFunction != nil {
					return PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
				} else if caseBasicType != nil {
					return PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
				} else if caseVoid != nil {
					return PtrTypeCheckErrorf("%s should be an interface to continue chained calls but found %s", varName, printableName(varType))
				} else {
					panic(fmt.Errorf("cases on %v", varType))
				}
			} else {
				caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
				if caseInterface != nil {
					return PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				} else if caseFunction != nil {
					if len(arguments) != len(caseFunction.Arguments) {
						return &TypecheckError{Message: fmt.Sprintf("Expected %d arguments but got %d", len(caseFunction.Arguments), len(arguments))}
					}
					for i2, argument := range arguments {
						expectedType := caseFunction.Arguments[i2].VariableType
						err := isExpressionOfExpectedType("", argument, expectedType, universe)
						if err != nil {
							return err
						}
					}
				} else if caseBasicType != nil {
					return PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				} else if caseVoid != nil {
					return PtrTypeCheckErrorf("%s should be a function for invocation but found %s", varName, printableName(varType))
				} else {
					panic(fmt.Errorf("cases on %v", varType))
				}
			}
		}
	}
	return nil
}

func printableName(varType VariableType) string {
	caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
	if caseInterface != nil {
		return caseInterface.Package + "." + caseInterface.Name
	} else if caseFunction != nil {
		result := "("
		for i, argumentType := range caseFunction.Arguments {
			if i > 0 {
				result = result + ", "
			}
			result = result + printableName(argumentType.VariableType)
		}
		return result + ") => " + printableName(caseFunction.ReturnType)
	} else if caseBasicType != nil {
		return caseBasicType.Type
	} else if caseVoid != nil {
		return "Void"
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}
