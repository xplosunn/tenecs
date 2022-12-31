package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/program"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"unicode"
)

func Typecheck(parsed parser.FileTopLevel) (*program.Information, error) {
	program := &program.Information{}
	pkg, imports, topLevelDeclarations := parser.FileTopLevelFields(parsed)
	err := validatePackage(pkg)
	if err != nil {
		return program, err
	}
	universe, err := resolveImports(imports, StdLib, StdLibInterfaceVariables)
	if err != nil {
		return program, err
	}
	modules, interfaces := splitTopLevelDeclarations(topLevelDeclarations)
	universe, err = validateInterfaces(interfaces, pkg, universe)
	if err != nil {
		return program, err
	}
	modulesMap, parserModulesMap, err := validateModulesImplements(modules, universe)
	if err != nil {
		return program, err
	}
	addImplementedInterfacesToProgramInformation(program, modulesMap)
	universeByModuleName, err := validateModulesVariableTypesAndExpressions(modulesMap, parserModulesMap, universe)
	if err != nil {
		return program, err
	}
	err = validateModulesVariableFunctionBlocks(modulesMap, parserModulesMap, universeByModuleName)
	if err != nil {
		return program, err
	}

	return program, nil
}

func addImplementedInterfacesToProgramInformation(programInfo *program.Information, programInterfacesMap map[string]*program.Interface) {
	interfaces := []*program.Interface{}
	for _, interf := range programInterfacesMap {
		interfaces = append(interfaces, interf)
	}
	programInfo.Interfaces = interfaces
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

func validatePackage(node parser.Package) *type_error.TypecheckError {
	identifier := parser.PackageFields(node)
	for _, r := range identifier {
		if !unicode.IsLower(r) {
			return type_error.PtrTypeCheckErrorf("package name should start with a lowercase letter")
		} else {
			return nil
		}
	}
	return nil
}

func resolveImports(nodes []parser.Import, stdLib Package, stdLibInterfaceVariables map[string]map[string]types.VariableType) (binding.Universe, *type_error.TypecheckError) {
	universe := binding.NewFromDefaults(DefaultTypesAvailableWithoutImport)
	for interfaceRef, variables := range stdLibInterfaceVariables {
		updatedUniverse, err := binding.CopyAddingGlobalInterfaceRefVariables(universe, interfaceRef, variables)
		if err != nil {
			return universe, err
		}
		universe = updatedUniverse
	}
	for _, node := range nodes {
		dotSeparatedNames := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			return universe, type_error.PtrTypeCheckErrorf("all interfaces belong to a package")
		}
		currPackage := stdLib
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name]
				if !ok {
					return universe, type_error.PtrTypeCheckErrorf("no package " + name + " found")
				}
				currPackage = p
				continue
			}
			interf, ok := currPackage.Interfaces[name]
			if !ok {
				return universe, type_error.PtrTypeCheckErrorf("no interface " + name + " found")
			}
			updatedUniverse, err := binding.CopyAddingType(universe, name, interf)
			if err != nil {
				return universe, err
			}
			universe = updatedUniverse
		}
	}
	return universe, nil
}

func validateTypeAnnotationInUniverse(typeAnnotation parser.TypeAnnotation, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	caseSingleNameType, caseFunctionType := typeAnnotation.Cases()
	if caseSingleNameType != nil {
		varType, ok := binding.GetTypeByTypeName(universe, caseSingleNameType.TypeName)
		if !ok {
			return nil, type_error.PtrTypeCheckErrorf("not found type: %s", caseSingleNameType.TypeName)
		}
		return varType, nil
	} else if caseFunctionType != nil {
		arguments := []types.FunctionArgument{}
		for _, argAnnotatedType := range caseFunctionType.Arguments {
			varType, err := validateTypeAnnotationInUniverse(argAnnotatedType, universe)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, types.FunctionArgument{
				Name:         "?",
				VariableType: varType,
			})
		}
		returnType, err := validateTypeAnnotationInUniverse(caseFunctionType.ReturnType, universe)
		if err != nil {
			return nil, err
		}
		return types.Function{
			Arguments:  arguments,
			ReturnType: returnType,
		}, nil
	} else {
		panic("Cases on typeAnnotation")
	}
}

func validateInterfaces(nodes []parser.Interface, pkg parser.Package, universe binding.Universe) (binding.Universe, *type_error.TypecheckError) {
	updatedUniverse := universe
	var err *type_error.TypecheckError
	for _, node := range nodes {
		updatedUniverse, err = binding.CopyAddingType(updatedUniverse, node.Name, types.Interface{
			Package: pkg.Identifier,
			Name:    node.Name,
		})
		if err != nil {
			return updatedUniverse, err
		}
	}
	for _, node := range nodes {
		name, parserVariables := parser.InterfaceFields(node)
		variables := map[string]types.VariableType{}
		for _, variable := range parserVariables {
			varType, err := validateTypeAnnotationInUniverse(variable.Type, updatedUniverse)
			if err != nil {
				return updatedUniverse, err
			}
			_, ok := variables[variable.Name]
			if ok {
				return updatedUniverse, type_error.PtrTypeCheckErrorf("more than one variable with name '%s'", variable.Name)
			}
			variables[variable.Name] = varType
		}
		interf := types.Interface{
			Package: pkg.Identifier,
			Name:    name,
		}
		updatedUniverse, err = binding.CopyAddingGlobalInterfaceVariables(updatedUniverse, interf, variables)
		if err != nil {
			return updatedUniverse, err
		}
	}
	return updatedUniverse, nil
}

func validateModulesImplements(nodes []parser.Module, universe binding.Universe) (map[string]*program.Interface, map[string]parser.Module, *type_error.TypecheckError) {
	implementedInterfacesMap := map[string]*program.Interface{}
	parserModulesMap := map[string]parser.Module{}
	for _, node := range nodes {
		implementing, name, constructorArgs, declarations := parser.ModuleFields(node)
		_ = declarations
		_ = constructorArgs
		_, ok := implementedInterfacesMap[name]
		if ok {
			return nil, nil, type_error.PtrTypeCheckErrorf("another module declared with name %s", name)
		}
		if implementing == "" {
			return nil, nil, type_error.PtrTypeCheckErrorf("module %s needs to implement some interface", name)
		}
		implementedInterface, err := validateImplementedInterfaces(implementing, universe)
		if err != nil {
			return nil, nil, err
		}
		implementedInterfacesMap[name] = &program.Interface{
			Implements: implementedInterface,
			Variables:  nil,
		}
		parserModulesMap[name] = node
	}
	return implementedInterfacesMap, parserModulesMap, nil
}

func validateImplementedInterfaces(implements string, universe binding.Universe) (types.Interface, *type_error.TypecheckError) {
	emptyInterface := types.Interface{}
	varType, ok := binding.GetTypeByTypeName(universe, implements)
	if !ok {
		return emptyInterface, type_error.PtrTypeCheckErrorf("not found interface with name %s", implements)
	}
	caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
	if caseInterface != nil {
		return *caseInterface, nil
	} else if caseFunction != nil {
		return emptyInterface, type_error.PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implements, printableName(varType))
	} else if caseBasicType != nil {
		return emptyInterface, type_error.PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implements, printableName(varType))
	} else if caseVoid != nil {
		return emptyInterface, type_error.PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implements, printableName(varType))
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func validateModulesVariableTypesAndExpressions(implementedInterfacesMap map[string]*program.Interface, parserModulesMap map[string]parser.Module, universe binding.Universe) (map[string]binding.Universe, *type_error.TypecheckError) {
	universeByModuleName := map[string]binding.Universe{}

	for moduleName, parserModule := range parserModulesMap {
		universeByModuleName[moduleName] = universe
		implementedInterface := implementedInterfacesMap[moduleName].Implements
		implementedInterfaceVariables, err := binding.GetGlobalInterfaceVariables(universe, implementedInterface)
		if err != nil {
			return nil, err
		}
		for interfaceVarName, _ := range implementedInterfaceVariables {
			found := false
			for _, declaration := range parserModule.Declarations {
				if declaration.Name == interfaceVarName {
					found = true
					break
				}
			}
			for _, constructorArg := range parserModule.ConstructorArgs {
				if constructorArg.Name == interfaceVarName {
					found = true
					break
				}
			}
			if !found {
				return nil, type_error.PtrTypeCheckErrorf("variable %s of interface %s missing in module %s", interfaceVarName, implementedInterface.Name, moduleName)
			}
		}
	}

	for moduleName, parserModule := range parserModulesMap {
		moduleConstructor := binding.Constructor{
			Arguments:  []types.FunctionArgument{},
			ReturnType: implementedInterfacesMap[moduleName].Implements,
		}
		for _, constructorArg := range parserModule.ConstructorArgs {
			if constructorArg.Name == moduleName {
				return nil, type_error.PtrTypeCheckErrorf("variable %s cannot have the same name as the module", constructorArg.Name)
			}
			varType, err := validateTypeAnnotationInUniverse(constructorArg.Type, universeByModuleName[moduleName])
			if err != nil {
				return nil, err
			}
			typeOfInterfaceVariableWithSameName, err := getVariableWithSameNameInInterface(constructorArg.Public, constructorArg.Name, implementedInterfacesMap[moduleName].Implements, universe)
			if err != nil {
				return nil, err
			}
			if typeOfInterfaceVariableWithSameName != nil {
				if !variableTypeEq(varType, *typeOfInterfaceVariableWithSameName) {
					return nil, type_error.PtrTypeCheckErrorf("variable %s should be of type %s but is of type %s", constructorArg.Name, printableName(*typeOfInterfaceVariableWithSameName), printableName(varType))
				}
			}
			updatedUniverse, err := binding.CopyAddingVariable(universeByModuleName[moduleName], constructorArg.Name, varType)
			if err != nil {
				return nil, err
			}
			universeByModuleName[moduleName] = updatedUniverse
			moduleConstructor.Arguments = append(moduleConstructor.Arguments, types.FunctionArgument{
				Name:         constructorArg.Name,
				VariableType: varType,
			})
		}

		for moduleNameWithUniverse, _ := range universeByModuleName {
			updatedUniverse, err := binding.CopyAddingConstructor(universeByModuleName[moduleNameWithUniverse], moduleName, moduleConstructor)
			if err != nil {
				return nil, err
			}
			universeByModuleName[moduleNameWithUniverse] = updatedUniverse
		}
	}

	for moduleName, parserModule := range parserModulesMap {
		for _, node := range parserModule.Declarations {
			if node.Name == moduleName {
				return nil, type_error.PtrTypeCheckErrorf("variable %s cannot have the same name as the module", node.Name)
			}
			typeOfInterfaceVariableWithSameName, err := getVariableWithSameNameInInterface(node.Public, node.Name, implementedInterfacesMap[moduleName].Implements, universe)
			if err != nil {
				return nil, err
			}

			u2, varType, err := validateModuleVariableTypeAndExpression(node, typeOfInterfaceVariableWithSameName, universeByModuleName[moduleName])
			if err != nil {
				return nil, err
			}
			universe = u2
			if implementedInterfacesMap[moduleName].Variables == nil {
				implementedInterfacesMap[moduleName].Variables = map[string]program.Expression{}
			}
			_, ok := implementedInterfacesMap[moduleName].Variables[node.Name]
			if ok {
				return nil, type_error.PtrTypeCheckErrorf("two variables declared in module %s with name %s", moduleName, node.Name)
			}
			implementedInterfacesMap[moduleName].Variables[node.Name] = varType
		}
	}

	return universeByModuleName, nil
}

func getVariableWithSameNameInInterface(varIsPublic bool, varNameToSearch string, implements types.Interface, universe binding.Universe) (*types.VariableType, *type_error.TypecheckError) {
	var nameOfInterfaceWithVariableWithSameName string
	var typeOfInterfaceVariableWithSameName *types.VariableType
	implementedInterfaceVariables, err := binding.GetGlobalInterfaceVariables(universe, implements)
	if err != nil {
		return nil, err
	}
	for varName, varType := range implementedInterfaceVariables {
		if varName == varNameToSearch {
			typeOfInterfaceVariableWithSameName = &varType
			nameOfInterfaceWithVariableWithSameName = implements.Name
			break
		}
	}

	if typeOfInterfaceVariableWithSameName == nil && varIsPublic {
		return nil, type_error.PtrTypeCheckErrorf("variable %s can't be public as no implemented interface has a variable with the same name", varNameToSearch)
	}

	if typeOfInterfaceVariableWithSameName != nil && !varIsPublic {
		return nil, type_error.PtrTypeCheckErrorf("variable %s should be public as it's in implemented interface %s", varNameToSearch, nameOfInterfaceWithVariableWithSameName)
	}

	return typeOfInterfaceVariableWithSameName, nil
}

func validateModuleVariableTypeAndExpression(node parser.ModuleDeclaration, typeOfInterfaceVariableWithSameName *types.VariableType, universe binding.Universe) (binding.Universe, program.Expression, *type_error.TypecheckError) {
	if typeOfInterfaceVariableWithSameName == nil {
		u2, programExp, err := determineVariableTypeOfExpression(node.Name, node.Expression, universe)
		universe = u2
		return universe, programExp, err
	}
	u2, programExp, err := expectVariableTypeOfExpression(node.Expression, *typeOfInterfaceVariableWithSameName, universe)
	if err != nil {
		return nil, nil, err
	}
	universe = u2
	return universe, programExp, nil
}

func validateModulesVariableFunctionBlocks(implementedInterfacesMap map[string]*program.Interface, parserModulesMap map[string]parser.Module, universeByModuleName map[string]binding.Universe) *type_error.TypecheckError {
	for moduleName, module := range implementedInterfacesMap {
		for varName, programExpression := range module.Variables {
			varType := program.VariableTypeOfExpression(programExpression)
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			var function *types.Function
			if caseInterface != nil {
				continue
			} else if caseFunction != nil {
				function = caseFunction
			} else if caseBasicType != nil {
				continue
			} else if caseVoid != nil {
				continue
			} else {
				panic(fmt.Errorf("cases on %v", program.VariableTypeOfExpression(programExpression)))
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

			blockUniverse, err := binding.CopyAddingFunctionArguments(universeByModuleName[moduleName], function.Arguments)
			if err != nil {
				return err
			}
			blockUniverse, err = binding.CopyAddingVariables(blockUniverse, module.Variables)
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

func searchAndValidateFunctionBlocks(expression parser.Expression, universe binding.Universe, inferredFunction *types.Function) (binding.Universe, *type_error.TypecheckError) {
	caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := expression.Cases()
	if caseLiteralExp != nil {
		return universe, nil
	} else if caseReferenceOrInvocation != nil {
		if caseReferenceOrInvocation.Arguments != nil {
			_, programExp, err := determineVariableTypeOfExpression("--", parser.ReferenceOrInvocation{
				DotSeparatedVars: caseReferenceOrInvocation.DotSeparatedVars,
				Arguments:        nil,
			}, universe)
			if err != nil {
				return universe, err
			}
			variableType := program.VariableTypeOfExpression(programExp)
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
			return universe, nil
		} else {
			_, _, err := determineVariableTypeOfExpression("--", parser.ReferenceOrInvocation{
				DotSeparatedVars: caseReferenceOrInvocation.DotSeparatedVars,
				Arguments:        nil,
			}, universe)
			if err != nil {
				return universe, err
			}
			return universe, nil
		}
	} else if caseLambda != nil {
		var function types.Function
		if inferredFunction != nil {
			_, _, err := expectVariableTypeOfExpression(expression, *inferredFunction, universe)
			if err != nil {
				return universe, err
			}
			function = *inferredFunction
		} else {
			u2, programExp, err := determineVariableTypeOfExpression("<>", expression, universe)
			universe = u2
			if err != nil {
				return universe, err
			}
			varType := program.VariableTypeOfExpression(programExp)
			caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
			_ = caseInterface
			_ = caseBasicType
			_ = caseVoid
			if caseFunction == nil {
				panic("expected caseFunction on lambda")
			}
			function = *caseFunction
		}
		blockUniverse, err := binding.CopyAddingFunctionArguments(universe, function.Arguments)
		if err != nil {
			return universe, err
		}
		err = validateFunctionBlock(caseLambda.Block, function.ReturnType, blockUniverse)
		if err != nil {
			return universe, err
		}
		return universe, nil
	} else if caseDeclaration != nil {
		u, programExp, err := determineVariableTypeOfExpression(caseDeclaration.Name, caseDeclaration.Expression, universe)
		if err != nil {
			return universe, err
		}
		universe = u
		u, err = searchAndValidateFunctionBlocks(caseDeclaration.Expression, universe, nil)
		if err != nil {
			return universe, err
		}
		universe = u
		varType := program.VariableTypeOfExpression(programExp)
		universe, err = binding.CopyAddingVariable(universe, caseDeclaration.Name, varType)
		if err != nil {
			return universe, err
		}
		return universe, nil
	} else if caseIf != nil {
		scopedBlock := func(expressions []parser.Expression) *type_error.TypecheckError {
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

func validateFunctionBlock(block []parser.Expression, functionReturnType types.VariableType, universe binding.Universe) *type_error.TypecheckError {
	if len(block) == 0 {
		if !variableTypeEq(functionReturnType, void) {
			return type_error.PtrTypeCheckErrorf("Function has return type of %s but has empty body", printableName(functionReturnType))
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
			_, _, err := expectVariableTypeOfExpression(expression, functionReturnType, updatedUniverse)
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

func printableName(varType types.VariableType) string {
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
