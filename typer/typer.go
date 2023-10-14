package typer

import (
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"unicode"
)

func TypecheckSingleFile(parsed parser.FileTopLevel) (*ast.Program, error) {
	pkgName := ""
	for i, name := range parsed.Package.DotSeparatedNames {
		if i > 0 {
			pkgName += "."
		}
		pkgName += name.String
	}
	return TypecheckPackage(pkgName, map[string]parser.FileTopLevel{"file.10x": parsed})
}

func TypecheckPackage(pkgName string, parsedPackage map[string]parser.FileTopLevel) (*ast.Program, error) {
	if len(parsedPackage) == 0 {
		return nil, errors.New("no files provided for typechecking")
	}

	if len(parsedPackage) == 0 {
		panic("no files in package when typechecking " + pkgName)
	}

	for _, topLevel := range parsedPackage {
		fileDeclaredPackage := ""
		for i, name := range topLevel.Package.DotSeparatedNames {
			if i > 0 {
				fileDeclaredPackage += "."
			}
			fileDeclaredPackage += name.String
		}
		if pkgName != fileDeclaredPackage {
			panic("tried to typecheck files from different packages as if they belonged to the same package")
		}
		err := validatePackage(topLevel.Package)
		if err != nil {
			return nil, err
		}
	}

	universe := binding.NewFromDefaults(standard_library.DefaultTypesAvailableWithoutImport)
	universe, err := addAllInterfaceFieldsToUniverse(universe, standard_library.StdLib)
	if err != nil {
		return nil, err
	}

	program := ast.Program{
		NativeFunctions:        map[string]*types.Function{},
		NativeFunctionPackages: map[string]string{},
	}
	for file, fileTopLevel := range parsedPackage {
		programNativeFunctions, programNativeFunctionPackages, u, err := resolveImports(fileTopLevel.Imports, standard_library.StdLib, file, universe)
		if err != nil {
			return nil, err
		}
		universe = u
		for functionName, function := range programNativeFunctions {
			if program.NativeFunctions[functionName] != nil && program.NativeFunctionPackages[functionName] != programNativeFunctionPackages[functionName] {
				return nil, type_error.PtrOnNodef(fileTopLevel.Package.DotSeparatedNames[0].Node, "TODO: unsupported imports of different functions from standard library with same name on different files of same package")
			}
			program.NativeFunctions[functionName] = function
			program.NativeFunctionPackages[functionName] = programNativeFunctionPackages[functionName]
		}
	}

	structsInAllFiles := []parser.Struct{}
	interfacesInAllFiles := []parser.Interface{}
	declarationsPerFile := map[string][]parser.Declaration{}
	for file, fileTopLevel := range parsedPackage {
		declarations, interfaces, structs := splitTopLevelDeclarations(fileTopLevel.TopLevelDeclarations)
		structsInAllFiles = append(structsInAllFiles, structs...)
		interfacesInAllFiles = append(interfacesInAllFiles, interfaces...)
		declarationsPerFile[file] = declarations
	}

	programStructFunctions, universe, err := validateStructs(structsInAllFiles, pkgName, universe)
	if err != nil {
		return nil, err
	}
	program.StructFunctions = programStructFunctions
	universe, err = validateInterfaces(interfacesInAllFiles, pkgName, universe)
	if err != nil {
		return nil, err
	}
	program.FieldsByType = binding.GetAllFields(universe)

	declarationsMap, err := TypecheckDeclarations(nil, parser.Node{}, declarationsPerFile, universe)
	if err != nil {
		return nil, err
	}
	programDeclarations := []*ast.Declaration{}
	for varName, varExp := range declarationsMap {
		programDeclarations = append(programDeclarations, &ast.Declaration{
			Name:       varName,
			Expression: varExp,
		})
	}
	program.Declarations = programDeclarations

	return &program, nil
}

func validatePackage(node parser.Package) *type_error.TypecheckError {
	for _, name := range node.DotSeparatedNames {
		if !unicode.IsLower(rune(name.String[0])) {
			return type_error.PtrOnNodef(name.Node, "package name should start with a lowercase letter")
		}
	}
	return nil
}

func validateInterfaces(nodes []parser.Interface, pkgName string, universe binding.Universe) (binding.Universe, *type_error.TypecheckError) {
	updatedUniverse := universe
	var err *type_error.TypecheckError
	for _, node := range nodes {
		variables := map[string]types.VariableType{}
		for _, variable := range node.Variables {
			variables[variable.Name.String] = nil
		}
		genericNames := []string{}
		generics := []types.VariableType{}
		for _, generic := range node.Generics {
			genericNames = append(genericNames, generic.String)
			generics = append(generics, &types.TypeArgument{Name: generic.String})
		}
		updatedUniverse, err = binding.CopyAddingTypeToAllFiles(updatedUniverse, node.Name, &types.KnownType{
			Package:          pkgName,
			Name:             node.Name.String,
			DeclaredGenerics: genericNames,
			Generics:         generics,
			ValidStructField: false,
		})
		if err != nil {
			return nil, err
		}
	}
	for _, node := range nodes {
		name, generics, parserVariables := parser.InterfaceFields(node)
		_ = name
		localUniverse := updatedUniverse
		for _, generic := range generics {
			u, err := binding.CopyAddingTypeToAllFiles(localUniverse, generic, &types.TypeArgument{Name: generic.String})
			if err != nil {
				return nil, err
			}
			localUniverse = u
		}
		variables := map[string]types.VariableType{}
		for _, variable := range parserVariables {
			varType, err := validateTypeAnnotationInUniverse(variable.Type, "", localUniverse)
			if err != nil {
				return nil, err
			}
			_, ok := variables[variable.Name.String]
			if ok {
				return nil, type_error.PtrOnNodef(variable.Name.Node, "more than one variable with name '%s'", variable.Name.String)
			}
			variables[variable.Name.String] = varType
		}
		updatedUniverse, err = binding.CopyAddingFields(updatedUniverse, pkgName, node.Name, variables)
		if err != nil {
			return nil, err
		}
	}
	return updatedUniverse, nil
}

func splitTopLevelDeclarations(topLevelDeclarations []parser.TopLevelDeclaration) ([]parser.Declaration, []parser.Interface, []parser.Struct) {
	declarations := []parser.Declaration{}
	interfaces := []parser.Interface{}
	structs := []parser.Struct{}
	for _, topLevelDeclaration := range topLevelDeclarations {
		parser.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration parser.Declaration) {
				declarations = append(declarations, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Interface) {
				interfaces = append(interfaces, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Struct) {
				structs = append(structs, topLevelDeclaration)
			},
		)
	}
	return declarations, interfaces, structs
}

func addAllInterfaceFieldsToUniverse(universe binding.Universe, pkg standard_library.Package) (binding.Universe, *type_error.TypecheckError) {
	var err *type_error.TypecheckError
	for interfaceName, interfaceWithFields := range pkg.Interfaces {
		universe, err = binding.CopyAddingFields(universe, interfaceWithFields.Interface.Package, parser.Name{
			String: interfaceName,
		}, interfaceWithFields.Fields)
		if err != nil {
			return nil, err
		}
	}
	for _, nestedPkg := range pkg.Packages {
		universe, err = addAllInterfaceFieldsToUniverse(universe, nestedPkg)
		if err != nil {
			return nil, err
		}
	}
	return universe, nil
}

func resolveImports(nodes []parser.Import, stdLib standard_library.Package, file string, universe binding.Universe) (map[string]*types.Function, map[string]string, binding.Universe, *type_error.TypecheckError) {
	nativeFunctions := map[string]*types.Function{}
	nativeFunctionPackages := map[string]string{}
	for _, node := range nodes {
		dotSeparatedNames := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			errNode := node.Node
			if len(dotSeparatedNames) > 0 {
				errNode = dotSeparatedNames[0].Node
			}
			return nil, nil, nil, type_error.PtrOnNodef(errNode, "all interfaces belong to a package")
		}
		currPackage := stdLib
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name.String]
				if !ok {
					return nil, nil, nil, type_error.PtrOnNodef(name.Node, "no package "+name.String+" found")
				}
				currPackage = p
				continue
			}
			interf, ok := currPackage.Interfaces[name.String]
			if ok {
				updatedUniverse, err := binding.CopyAddingTypeToFile(universe, file, name, interf.Interface)
				if err != nil {
					return nil, nil, nil, err
				}
				universe = updatedUniverse
				continue
			}
			varTypeToImport, ok := currPackage.Variables[name.String]
			if ok {
				updatedUniverse, err := binding.CopyAddingVariable(universe, name, varTypeToImport)
				if err != nil {
					return nil, nil, nil, err
				}
				universe = updatedUniverse
				fn, ok := varTypeToImport.(*types.Function)
				if !ok {
					panic(fmt.Sprintf("todo resolveImports not native function but %T", varTypeToImport))
				}
				nativeFunctions[name.String] = fn
				pkg := ""
				for i, name := range dotSeparatedNames {
					if i > 0 {
						pkg += "_"
					}
					pkg += name.String
				}
				nativeFunctionPackages[name.String] = pkg
				continue
			}

			return nil, nil, nil, type_error.PtrOnNodef(name.Node, "didn't find "+name.String+" while importing")
		}
	}
	return nativeFunctions, nativeFunctionPackages, universe, nil
}

func validateStructs(nodes []parser.Struct, pkgName string, universe binding.Universe) (map[string]*types.Function, binding.Universe, *type_error.TypecheckError) {
	constructors := map[string]*types.Function{}
	var err *type_error.TypecheckError
	for _, node := range nodes {
		genericNames := []string{}
		genericTypeArgs := []types.VariableType{}
		for _, generic := range node.Generics {
			genericNames = append(genericNames, generic.String)
			genericTypeArgs = append(genericTypeArgs, &types.TypeArgument{Name: generic.String})
		}
		universe, err = binding.CopyAddingTypeToAllFiles(universe, node.Name, &types.KnownType{
			Package:          pkgName,
			Name:             node.Name.String,
			DeclaredGenerics: genericNames,
			Generics:         genericTypeArgs,
			ValidStructField: true,
		})
		if err != nil {
			return nil, nil, err
		}
	}
	for _, node := range nodes {
		structName, generics, parserVariables := parser.StructFields(node)
		localUniverse := universe
		for _, generic := range generics {
			u, err := binding.CopyAddingTypeToAllFiles(localUniverse, generic, &types.TypeArgument{Name: generic.String})
			if err != nil {
				return nil, nil, err
			}
			localUniverse = u
		}
		constructorArgs := []types.FunctionArgument{}
		variables := map[string]types.VariableType{}
		for _, variable := range parserVariables {
			varType, err := validateTypeAnnotationInUniverse(variable.Type, "", localUniverse)
			if err != nil {
				return nil, nil, type_error.PtrOnNodef(variable.Name.Node, "%s (are you using an incomparable type?)", err.Error())
			}
			if !varType.CanBeStructField() {
				return nil, nil, type_error.PtrOnNodef(variable.Name.Node, "not a valid struct var type %s", printableName(varType))
			}
			constructorArgs = append(constructorArgs, types.FunctionArgument{
				Name:         variable.Name.String,
				VariableType: varType,
			})
			variables[variable.Name.String] = varType
		}
		universe, err = binding.CopyAddingFields(universe, pkgName, structName, variables)

		genericNames := []types.VariableType{}
		for _, generic := range generics {
			genericNames = append(genericNames, &types.TypeArgument{
				Name: generic.String,
			})
		}
		maybeStruc, resolutionErr := binding.GetTypeByTypeName(localUniverse, "", structName.String, genericNames)
		if resolutionErr != nil {
			return nil, nil, TypecheckErrorFromResolutionError(structName.Node, resolutionErr)
		}
		struc, ok := maybeStruc.(*types.KnownType)
		if !ok {
			return nil, nil, type_error.PtrOnNodef(structName.Node, "expected struct type in validateStructs")
		}

		genericStrings := []string{}
		for _, generic := range generics {
			genericStrings = append(genericStrings, generic.String)
		}
		if generics == nil {
			genericStrings = nil
		}
		constructorVarType := &types.Function{
			Generics:   genericStrings,
			Arguments:  constructorArgs,
			ReturnType: struc,
		}
		universe, err = binding.CopyAddingVariable(universe, structName, constructorVarType)
		constructors[structName.String] = constructorVarType
	}
	return constructors, universe, nil
}

func TypecheckDeclarations(expectedTypes *map[string]types.VariableType, node parser.Node, declarationsPerFile map[string][]parser.Declaration, universe binding.Universe) (map[string]ast.Expression, *type_error.TypecheckError) {
	typesByName := map[parser.Name]types.VariableType{}

	for file, declarations := range declarationsPerFile {
		for _, declaration := range declarations {
			if expectedTypes != nil {
				typesByName[declaration.Name] = (*expectedTypes)[declaration.Name.String]
			}
			if declaration.TypeAnnotation != nil {
				annotatedVarType, err := validateTypeAnnotationInUniverse(*declaration.TypeAnnotation, file, universe)
				if err != nil {
					return nil, err
				}
				if typesByName[declaration.Name] == nil {
					typesByName[declaration.Name] = annotatedVarType
				} else if !types.VariableTypeEq(typesByName[declaration.Name], annotatedVarType) {
					return nil, type_error.PtrOnNodef(node, "annotated type %s doesn't match the expected %s", printableName(annotatedVarType), printableName(typesByName[declaration.Name]))
				}
			}
			if typesByName[declaration.Name] == nil {
				varType, err := typeOfExpressionBox(declaration.ExpressionBox, file, universe)
				if err != nil {
					return nil, err
				}
				typesByName[declaration.Name] = varType
			}
		}
	}

	if expectedTypes != nil {
		for expectedVarName, _ := range *expectedTypes {
			found := false
			for varName, _ := range typesByName {
				if varName.String == expectedVarName {
					found = true
					break
				}
			}
			if !found {
				return nil, type_error.PtrOnNodef(node, "missing declaration for variable %s", expectedVarName)
			}
		}
	}

	for varName, varType := range typesByName {
		var err *type_error.TypecheckError
		universe, err = binding.CopyAddingVariable(universe, varName, varType)
		if err != nil {
			return nil, err
		}
	}

	result := map[string]ast.Expression{}

	for file, declarations := range declarationsPerFile {
		for _, declaration := range declarations {
			expectedType := typesByName[declaration.Name]
			if expectedType == nil {
				panic("nil expectedType on TypecheckDeclarations")
			}
			astExp, err := expectTypeOfExpressionBox(expectedType, declaration.ExpressionBox, file, universe)
			if err != nil {
				return nil, err
			}
			result[declaration.Name.String] = astExp
		}
	}

	return result, nil
}

func validateTypeAnnotationInUniverse(typeAnnotation parser.TypeAnnotation, file string, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	switch len(typeAnnotation.OrTypes) {
	case 0:
		return nil, type_error.PtrOnNodef(typeAnnotation.Node, "unexpected error validateTypeAnnotationInUniverse no types found")
	case 1:
		elem := typeAnnotation.OrTypes[0]
		return validateTypeAnnotationElementInUniverse(elem, file, universe)
	default:
		elements := []types.VariableType{}
		for _, element := range typeAnnotation.OrTypes {
			newElement, err := validateTypeAnnotationElementInUniverse(element, file, universe)
			if err != nil {
				return nil, err
			}
			elements = append(elements, newElement)
		}
		return &types.OrVariableType{
			Elements: elements,
		}, nil
	}
}

func validateTypeAnnotationElementInUniverse(typeAnnotationElement parser.TypeAnnotationElement, file string, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	var varType types.VariableType
	var err *type_error.TypecheckError
	parser.TypeAnnotationElementExhaustiveSwitch(
		typeAnnotationElement,
		func(typeAnnotation parser.SingleNameType) {
			genericTypes := []types.VariableType{}
			for _, generic := range typeAnnotation.Generics {
				genericVarType, err2 := validateTypeAnnotationInUniverse(generic, file, universe)
				if err2 != nil {
					err = err2
					return
				}
				if !genericVarType.CanBeStructField() {
					err = type_error.PtrOnNodef(generic.Node, "not a valid generic: %s", printableName(varType))
					return
				}
				genericTypes = append(genericTypes, genericVarType)
			}
			varType2, err2 := binding.GetTypeByTypeName(universe, file, typeAnnotation.TypeName.String, genericTypes)
			varType = varType2
			err = TypecheckErrorFromResolutionError(typeAnnotation.TypeName.Node, err2)
		},
		func(typeAnnotation parser.FunctionType) {
			localUniverse := universe
			for _, generic := range typeAnnotation.Generics {
				localUniverse, err = binding.CopyAddingTypeToFile(localUniverse, file, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return
				}
			}
			arguments := []types.FunctionArgument{}
			for _, argAnnotatedType := range typeAnnotation.Arguments {
				varType, err = validateTypeAnnotationInUniverse(argAnnotatedType, file, localUniverse)
				if err != nil {
					return
				}
				arguments = append(arguments, types.FunctionArgument{
					Name:         "?",
					VariableType: varType,
				})
			}
			var returnType types.VariableType
			returnType, err = validateTypeAnnotationInUniverse(typeAnnotation.ReturnType, file, localUniverse)
			if err != nil {
				return
			}
			generics := []string{}
			for _, generic := range typeAnnotation.Generics {
				generics = append(generics, generic.String)
			}
			if typeAnnotation.Generics == nil {
				generics = nil
			}
			varType = &types.Function{
				Generics:   generics,
				Arguments:  arguments,
				ReturnType: returnType,
			}
		},
	)
	return varType, err
}
