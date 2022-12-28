package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"unicode"
)

func Typecheck(parsed parser.FileTopLevel) error {
	pkg, imports, topLevelDeclarations := parser.FileTopLevelFields(parsed)
	err := validatePackage(pkg)
	if err != nil {
		return err
	}
	universe, err := resolveImports(imports, StdLib)
	if err != nil {
		return err
	}
	modules, interfaces := splitTopLevelDeclarations(topLevelDeclarations)
	universe, err = validateInterfaces(interfaces, pkg, universe)
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

func splitTopLevelDeclarations(topLevelDeclarations []parser.TopLevelDeclaration) ([]parser.Module, []parser.Interface) {
	modules := []parser.Module{}
	interfaces := []parser.Interface{}
	for _, topLevelDeclaration := range topLevelDeclarations {
		caseModule, caseInterface := topLevelDeclaration.Cases()
		if caseModule != nil {
			modules = append(modules, *caseModule)
		} else if caseInterface != nil {
			interfaces = append(interfaces, *caseInterface)
		} else {
			panic("cases on topLevelDeclaration")
		}
	}
	return modules, interfaces
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

func validateTypeAnnotationInUniverse(typeAnnotation parser.TypeAnnotation, universe Universe) (VariableType, *TypecheckError) {
	caseSingleNameType, caseFunctionType := typeAnnotation.Cases()
	if caseSingleNameType != nil {
		varType, ok := universe.TypeByTypeName.Get(caseSingleNameType.TypeName)
		if !ok {
			return nil, PtrTypeCheckErrorf("not found type: %s", caseSingleNameType.TypeName)
		}
		return varType, nil
	} else if caseFunctionType != nil {
		arguments := []FunctionArgument{}
		for _, argAnnotatedType := range caseFunctionType.Arguments {
			varType, err := validateTypeAnnotationInUniverse(argAnnotatedType, universe)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, FunctionArgument{
				Name:         "?",
				VariableType: varType,
			})
		}
		returnType, err := validateTypeAnnotationInUniverse(caseFunctionType.ReturnType, universe)
		if err != nil {
			return nil, err
		}
		return Function{
			Arguments:  arguments,
			ReturnType: returnType,
		}, nil
	} else {
		panic("Cases on typeAnnotation")
	}
}

func validateInterfaces(nodes []parser.Interface, pkg parser.Package, universe Universe) (Universe, *TypecheckError) {
	updatedUniverse := universe
	var err *TypecheckError
	for _, node := range nodes {
		name, parserVariables := parser.InterfaceFields(node)
		variables := map[string]VariableType{}
		for _, variable := range parserVariables {
			varType, err := validateTypeAnnotationInUniverse(variable.Type, universe)
			if err != nil {
				return updatedUniverse, err
			}
			_, ok := variables[variable.Name]
			if ok {
				return updatedUniverse, PtrTypeCheckErrorf("more than one variable with name '%s'", variable.Name)
			}
			variables[variable.Name] = varType
		}
		varType := Interface{
			Package:   pkg.Identifier,
			Name:      name,
			Variables: variables,
		}
		updatedUniverse, err = copyUniverseAddingType(updatedUniverse, name, varType)
		if err != nil {
			return updatedUniverse, err
		}
	}
	return updatedUniverse, nil
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

			varType, err := validateModuleVariableTypeAndExpression(node, typeOfInterfaceVariableWithSameName, universe)
			if err != nil {
				return err
			}
			if modulesMap[moduleName].Variables == nil {
				modulesMap[moduleName].Variables = map[string]VariableType{}
			}
			_, ok := modulesMap[moduleName].Variables[node.Name]
			if ok {
				return PtrTypeCheckErrorf("two variables declared in module %s with name %s", moduleName, node.Name)
			}
			modulesMap[moduleName].Variables[node.Name] = varType
		}
	}
	return nil
}

func validateModuleVariableTypeAndExpression(node parser.ModuleDeclaration, typeOfInterfaceVariableWithSameName *VariableType, universe Universe) (VariableType, *TypecheckError) {
	if typeOfInterfaceVariableWithSameName == nil {
		updatedUniverse, varType, err := determineVariableTypeOfExpression(node.Name, node.Expression, universe)
		_ = updatedUniverse
		return varType, err
	}
	err := expectVariableTypeOfExpression(node.Expression, *typeOfInterfaceVariableWithSameName, universe)
	if err != nil {
		return nil, err
	}
	return *typeOfInterfaceVariableWithSameName, nil
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
					caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := declaration.Expression.Cases()
					if caseLiteralExp != nil {
						panic(fmt.Errorf("unexpected caseLiteralExp on %s.%s", moduleName, varName))
					} else if caseReferenceOrInvocation != nil {
						panic(fmt.Errorf("unexpected caseReferenceOrInvocation on %s.%s", moduleName, varName))
					} else if caseLambda != nil {
						parserLambda = *caseLambda
					} else if caseDeclaration != nil {
						panic(fmt.Errorf("unexpected caseDeclaration on %s.%s", moduleName, varName))
					} else if caseIf != nil {
						panic(fmt.Errorf("unexpected if on %s.%s", moduleName, varName))
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

			err = validateFunctionBlock(parserLambda.Block, function.ReturnType, blockUniverse)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func searchAndValidateFunctionBlocks(expression parser.Expression, universe Universe, inferredFunction *Function) (Universe, *TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.Cases()
	if caseLiteralExp != nil {
		return universe, nil
	} else if caseReferenceOrInvocation != nil {
		if caseReferenceOrInvocation.Arguments != nil {
			_, variableType, err := determineVariableTypeOfExpression("--", parser.ReferenceOrInvocation{
				DotSeparatedVars: caseReferenceOrInvocation.DotSeparatedVars,
				Arguments:        nil,
			}, universe)
			if err != nil {
				return universe, err
			}
			caseInterface, caseFunction, caseBasicType, caseVoid := variableType.Cases()
			_ = caseInterface
			_ = caseBasicType
			_ = caseVoid
			if caseFunction == nil {
				panic(fmt.Sprintf("should be a function: %+v", variableType))
			}
			for i, arg := range caseReferenceOrInvocation.Arguments.Arguments {
				caseInterface, caseFunction, caseBasicType, caseVoid := caseFunction.Arguments[i].VariableType.Cases()
				_ = caseInterface
				_ = caseBasicType
				_ = caseVoid
				_, err := searchAndValidateFunctionBlocks(arg, universe, caseFunction)
				if err != nil {
					return universe, err
				}
			}
		}
		return universe, nil
	} else if caseLambda != nil {
		var function Function
		if inferredFunction != nil {
			err := expectVariableTypeOfExpression(expression, *inferredFunction, universe)
			if err != nil {
				return universe, err
			}
			function = *inferredFunction
		} else {
			u2, varType, err := determineVariableTypeOfExpression("<>", expression, universe)
			universe = u2
			if err != nil {
				return universe, err
			}
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			_ = caseInterface
			_ = caseBasicType
			_ = caseVoid
			if caseFunction == nil {
				panic("expected caseFunction on lambda")
			}
			function = *caseFunction
		}
		blockUniverse, err := copyUniverseAddingFunctionArguments(universe, function.Arguments)
		if err != nil {
			return universe, err
		}
		err = validateFunctionBlock(caseLambda.Block, function.ReturnType, blockUniverse)
		if err != nil {
			return universe, err
		}
		return universe, nil
	} else if caseDeclaration != nil {
		u, varType, err := determineVariableTypeOfExpression(caseDeclaration.Name, caseDeclaration.Expression, universe)
		if err != nil {
			return universe, err
		}
		universe = u
		u, err = searchAndValidateFunctionBlocks(caseDeclaration.Expression, universe, nil)
		if err != nil {
			return universe, err
		}
		universe = u
		universe, err = copyUniverseAddingVariable(universe, caseDeclaration.Name, varType)
		if err != nil {
			return universe, err
		}
		return universe, nil
	} else if caseIf != nil {
		scopedBlock := func(expressions []parser.Expression) *TypecheckError {
			scopeUniverse := universe
			for _, exp := range expressions {
				u2, err := searchAndValidateFunctionBlocks(exp, scopeUniverse, nil)
				if err != nil {
					return err
				}
				scopeUniverse = u2
			}
			return nil
		}
		err := scopedBlock(caseIf.ThenBlock)
		if err != nil {
			return universe, err
		}
		err = scopedBlock(caseIf.ElseBlock)
		if err != nil {
			return universe, err
		}
		return universe, nil
	} else {
		panic("cases on expression")
	}
}

func validateFunctionBlock(block []parser.Expression, functionReturnType VariableType, universe Universe) *TypecheckError {
	if len(block) == 0 {
		if !variableTypeEq(functionReturnType, void) {
			return PtrTypeCheckErrorf("Function has return type of %s but has empty body", printableName(functionReturnType))
		}
		return nil
	}
	updatedUniverse := universe
	for i, expression := range block {
		u2, err := searchAndValidateFunctionBlocks(expression, updatedUniverse, nil)
		if err != nil {
			return err
		}
		if i < len(block)-1 {
			updatedUniverse = u2
		} else {
			err := expectVariableTypeOfExpression(expression, functionReturnType, updatedUniverse)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func printableNameOfTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	caseSingleNameType, caseFunctionType := typeAnnotation.Cases()
	if caseSingleNameType != nil {
		return caseSingleNameType.TypeName
	} else if caseFunctionType != nil {
		result := "("
		for i, argument := range caseFunctionType.Arguments {
			if i > 0 {
				result += ", "
			}
			result += printableNameOfTypeAnnotation(argument)
		}
		return result + ") -> " + printableNameOfTypeAnnotation(caseFunctionType.ReturnType)
	} else {
		panic("cases on typeAnnotation")
	}
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
		return result + ") -> " + printableName(caseFunction.ReturnType)
	} else if caseBasicType != nil {
		return caseBasicType.Type
	} else if caseVoid != nil {
		return "Void"
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}
