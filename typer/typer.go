package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"unicode"
)

func Typecheck(parsed parser.FileTopLevel) (*ast.Program, error) {
	program := &ast.Program{}
	pkg, imports, topLevelDeclarations := parser.FileTopLevelFields(parsed)
	err := validatePackage(pkg)
	if err != nil {
		return program, err
	}
	universe, err := resolveImports(imports, StdLib, StdLibInterfaceVariables)
	if err != nil {
		return program, err
	}
	modules, interfaces, structs := splitTopLevelDeclarations(topLevelDeclarations)
	universe, err = validateStructs(structs, pkg, universe)
	if err != nil {
		return program, err
	}
	universe, err = validateInterfaces(interfaces, pkg, universe)
	if err != nil {
		return program, err
	}
	modulesMap, parserModulesMap, err := validateModulesImplements(modules, universe)
	if err != nil {
		return program, err
	}
	universeByModuleName, err := validateModulesVariableTypesAndExpressionsWithoutFunctionBlocks(modulesMap, parserModulesMap, universe)
	if err != nil {
		return program, err
	}
	for moduleName, universe := range universeByModuleName {
		module := modulesMap[moduleName]
		for varName, varExp := range module.Variables {
			caseLiteralExp, caseReferenceOrInvocation, caseLambda, caseDeclaration, caseIf := varExp.Cases()
			_ = caseLiteralExp
			_ = caseReferenceOrInvocation
			_ = caseDeclaration
			_ = caseIf
			if caseLambda != nil {
				var parserExp parser.Expression
				for _, parserDec := range parserModulesMap[moduleName].Declarations {
					if parserDec.Name == varName {
						parserExp = parserDec.Expression
						break
					}
				}
				_, exp, err := expectTypeOfExpression(true, parserExp, caseLambda.VariableType, universe)
				if err != nil {
					return nil, err
				}
				caseLambda.Block = exp.(ast.Function).Block
				module.Variables[varName] = caseLambda
			}
		}
	}
	addModulesToProgram(program, modulesMap)

	return program, nil
}

func addModulesToProgram(program *ast.Program, modulesMap map[string]*ast.Module) {
	modules := []*ast.Module{}
	for _, interf := range modulesMap {
		modules = append(modules, interf)
	}
	program.Modules = modules
}

func splitTopLevelDeclarations(topLevelDeclarations []parser.TopLevelDeclaration) ([]parser.Module, []parser.Interface, []parser.Struct) {
	modules := []parser.Module{}
	interfaces := []parser.Interface{}
	structs := []parser.Struct{}
	for _, topLevelDeclaration := range topLevelDeclarations {
		caseModule, caseInterface, caseStruct := topLevelDeclaration.Cases()
		if caseModule != nil {
			modules = append(modules, *caseModule)
		} else if caseInterface != nil {
			interfaces = append(interfaces, *caseInterface)
		} else if caseStruct != nil {
			structs = append(structs, *caseStruct)
		} else {
			panic("cases on topLevelDeclaration")
		}
	}
	return modules, interfaces, structs
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
		u2, err := binding.CopyAddingGlobalInterfaceRefVariables(universe, interfaceRef, variables)
		if err != nil {
			return nil, err
		}
		universe = u2
	}
	for _, node := range nodes {
		dotSeparatedNames := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			return nil, type_error.PtrTypeCheckErrorf("all interfaces belong to a package")
		}
		currPackage := stdLib
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name]
				if !ok {
					return nil, type_error.PtrTypeCheckErrorf("no package " + name + " found")
				}
				currPackage = p
				continue
			}
			interf, ok := currPackage.Interfaces[name]
			if !ok {
				return nil, type_error.PtrTypeCheckErrorf("no interface " + name + " found")
			}
			updatedUniverse, err := binding.CopyAddingType(universe, name, interf)
			if err != nil {
				return nil, err
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

func validateStructs(nodes []parser.Struct, pkg parser.Package, universe binding.Universe) (binding.Universe, *type_error.TypecheckError) {
	var err *type_error.TypecheckError
	for _, node := range nodes {
		universe, err = binding.CopyAddingType(universe, node.Name, types.Struct{
			Package: pkg.Identifier,
			Name:    node.Name,
		})
		if err != nil {
			return nil, err
		}
	}
	for _, node := range nodes {
		structName, parserVariables := parser.StructFields(node)
		constructorArgs := []types.FunctionArgument{}
		variables := map[string]types.StructVariableType{}
		for _, variable := range parserVariables {
			varType, err := validateTypeAnnotationInUniverse(variable.Type, universe)
			if err != nil {
				return nil, type_error.PtrTypeCheckErrorf("%s (are you using an incomparable type?)", err.Error())
			}
			structVarType, ok := types.StructVariableTypeFromVariableType(varType)
			if !ok {
				return nil, type_error.PtrTypeCheckErrorf("not a valid struct var type %s", printableName(varType))
			}
			constructorArgs = append(constructorArgs, types.FunctionArgument{
				Name:         variable.Name,
				VariableType: varType,
			})
			variables[variable.Name] = structVarType
		}
		struc := types.Struct{
			Package: pkg.Identifier,
			Name:    structName,
		}
		universe, err = binding.CopyAddingGlobalStructVariables(universe, struc, variables)
		if err != nil {
			return nil, err
		}
		universe, err = binding.CopyAddingConstructor(universe, structName, binding.Constructor{
			Arguments:  constructorArgs,
			ReturnType: struc,
		})
	}
	return universe, nil
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

func validateModulesImplements(nodes []parser.Module, universe binding.Universe) (map[string]*ast.Module, map[string]parser.Module, *type_error.TypecheckError) {
	modulesMap := map[string]*ast.Module{}
	parserModulesMap := map[string]parser.Module{}
	for _, node := range nodes {
		implementing, name, constructorArgs, declarations := parser.ModuleFields(node)
		_ = declarations
		_ = constructorArgs
		_, ok := modulesMap[name]
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
		modulesMap[name] = &ast.Module{
			Name:       name,
			Implements: implementedInterface,
			Variables:  nil,
		}
		parserModulesMap[name] = node
	}
	return modulesMap, parserModulesMap, nil
}

func validateImplementedInterfaces(implements string, universe binding.Universe) (types.Interface, *type_error.TypecheckError) {
	emptyInterface := types.Interface{}
	varType, ok := binding.GetTypeByTypeName(universe, implements)
	if !ok {
		return emptyInterface, type_error.PtrTypeCheckErrorf("not found interface with name %s", implements)
	}
	caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
	if caseStruct != nil {
		return emptyInterface, type_error.PtrTypeCheckErrorf("only interfaces can be implemented but %s is %s", implements, printableName(varType))
	} else if caseInterface != nil {
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

func validateModulesVariableTypesAndExpressionsWithoutFunctionBlocks(modulesMap map[string]*ast.Module, parserModulesMap map[string]parser.Module, universe binding.Universe) (map[string]binding.Universe, *type_error.TypecheckError) {
	universeByModuleName := map[string]binding.Universe{}

	for moduleName, parserModule := range parserModulesMap {
		universeByModuleName[moduleName] = universe
		implementedInterface := modulesMap[moduleName].Implements
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
			ReturnType: modulesMap[moduleName].Implements,
		}
		for _, constructorArg := range parserModule.ConstructorArgs {
			if constructorArg.Name == moduleName {
				return nil, type_error.PtrTypeCheckErrorf("variable %s cannot have the same name as the module", constructorArg.Name)
			}
			varType, err := validateTypeAnnotationInUniverse(constructorArg.Type, universeByModuleName[moduleName])
			if err != nil {
				return nil, err
			}
			typeOfInterfaceVariableWithSameName, err := getVariableWithSameNameInInterface(constructorArg.Public, constructorArg.Name, modulesMap[moduleName].Implements, universe)
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
			_, ok := modulesMap[moduleName].Variables[node.Name]
			if ok {
				return nil, type_error.PtrTypeCheckErrorf("two variables declared in module %s with name %s", moduleName, node.Name)
			}
			typeOfInterfaceVariableWithSameName, err := getVariableWithSameNameInInterface(node.Public, node.Name, modulesMap[moduleName].Implements, universe)
			if err != nil {
				return nil, err
			}

			u2, varType, err := validateModuleVariableTypeAndExpression(node, typeOfInterfaceVariableWithSameName, universeByModuleName[moduleName])
			if err != nil {
				return nil, err
			}
			if modulesMap[moduleName].Variables == nil {
				modulesMap[moduleName].Variables = map[string]ast.Expression{}
			}
			universeByModuleName[moduleName] = u2
			modulesMap[moduleName].Variables[node.Name] = varType
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

func validateModuleVariableTypeAndExpression(node parser.ModuleDeclaration, typeOfInterfaceVariableWithSameName *types.VariableType, universe binding.Universe) (binding.Universe, ast.Expression, *type_error.TypecheckError) {
	var programExp ast.Expression
	var err *type_error.TypecheckError
	if typeOfInterfaceVariableWithSameName == nil {
		universe, programExp, err = determineTypeOfExpression(false, node.Name, node.Expression, universe)
	} else {
		universe, programExp, err = expectTypeOfExpression(false, node.Expression, *typeOfInterfaceVariableWithSameName, universe)
	}
	if err != nil {
		return nil, nil, err
	}
	universe, err = binding.CopyAddingVariable(universe, node.Name, ast.VariableTypeOfExpression(programExp))
	if err != nil {
		return nil, nil, err
	}
	return universe, programExp, nil
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
	caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.Cases()
	if caseStruct != nil {
		return "struct " + caseStruct.Package + "." + caseStruct.Name
	} else if caseInterface != nil {
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
