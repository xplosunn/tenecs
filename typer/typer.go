package typer

import (
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"unicode"
)

func Typecheck(parsed parser.FileTopLevel) (*ast.Program, error) {
	program := &ast.Program{}
	pkg, imports, topLevelDeclarations := parser.FileTopLevelFields(parsed)
	err := validatePackage(pkg)
	if err != nil {
		return nil, err
	}
	programNativeFunctions, programNativeFunctionPackages, universe, err := resolveImports(imports, standard_library.StdLib)
	if err != nil {
		return nil, err
	}
	program.NativeFunctions = programNativeFunctions
	program.NativeFunctionPackages = programNativeFunctionPackages
	declarations, interfaces, structs := splitTopLevelDeclarations(topLevelDeclarations)
	programStructFunctions, universe, err := validateStructs(structs, pkg, universe)
	if err != nil {
		return nil, err
	}
	program.StructFunctions = programStructFunctions
	universe, err = validateInterfaces(interfaces, pkg, universe)
	if err != nil {
		return nil, err
	}

	declarationsMap, err := TypecheckDeclarations(nil, parser.Node{}, declarations, universe)
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

	return program, nil
}

func validatePackage(node parser.Package) *type_error.TypecheckError {
	for _, r := range node.Identifier.String {
		if !unicode.IsLower(r) {
			return type_error.PtrOnNodef(node.Identifier.Node, "package name should start with a lowercase letter")
		} else {
			return nil
		}
	}
	return nil
}

func validateInterfaces(nodes []parser.Interface, pkg parser.Package, universe binding.Universe) (binding.Universe, *type_error.TypecheckError) {
	updatedUniverse := universe
	var err *type_error.TypecheckError
	for _, node := range nodes {
		variables := map[string]types.VariableType{}
		for _, variable := range node.Variables {
			variables[variable.Name.String] = nil
		}
		updatedUniverse, err = binding.CopyAddingType(updatedUniverse, node.Name, &types.Interface{
			Package:   pkg.Identifier.String,
			Name:      node.Name.String,
			Variables: map[string]types.VariableType{},
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
			_, ok := variables[variable.Name.String]
			if ok {
				return updatedUniverse, type_error.PtrOnNodef(variable.Name.Node, "more than one variable with name '%s'", variable.Name.String)
			}
			variables[variable.Name.String] = varType
		}
		maybeInterf, err := binding.GetTypeByTypeName(updatedUniverse, name.String, []string{})
		if err != nil {
			return nil, TypecheckErrorFromResolutionError(node.Name.Node, err)
		}
		interf, ok := maybeInterf.(*types.Interface)
		if !ok {
			panic("expected interface type in validateInterfaces")
		}
		interf.Variables = variables
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

func resolveImports(nodes []parser.Import, stdLib standard_library.Package) (map[string]*types.Function, map[string]string, binding.Universe, *type_error.TypecheckError) {
	universe := binding.NewFromDefaults(standard_library.DefaultTypesAvailableWithoutImport)
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
				updatedUniverse, err := binding.CopyAddingType(universe, name, interf)
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

func validateStructs(nodes []parser.Struct, pkg parser.Package, universe binding.Universe) (map[string]*types.Function, binding.Universe, *type_error.TypecheckError) {
	constructors := map[string]*types.Function{}
	var err *type_error.TypecheckError
	for _, node := range nodes {
		genericNames := []string{}
		for _, generic := range node.Generics {
			genericNames = append(genericNames, generic.String)
		}
		universe, err = binding.CopyAddingType(universe, node.Name, &types.Struct{
			Package:      pkg.Identifier.String,
			Name:         node.Name.String,
			GenericCount: len(genericNames),
		})
		if err != nil {
			return nil, nil, err
		}
	}
	for _, node := range nodes {
		structName, generics, parserVariables := parser.StructFields(node)
		localUniverse := universe
		for _, generic := range generics {
			u, err := binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic.String})
			if err != nil {
				return nil, nil, err
			}
			localUniverse = u
		}
		constructorArgs := []types.FunctionArgument{}
		variables := map[string]types.StructFieldVariableType{}
		for _, variable := range parserVariables {
			varType, err := validateTypeAnnotationInUniverse(variable.Type, localUniverse)
			if err != nil {
				return nil, nil, type_error.PtrOnNodef(variable.Name.Node, "%s (are you using an incomparable type?)", err.Error())
			}
			structVarType, ok := types.StructFieldVariableTypeFromVariableType(varType)
			if !ok {
				return nil, nil, type_error.PtrOnNodef(variable.Name.Node, "not a valid struct var type %s", printableName(varType))
			}
			constructorArgs = append(constructorArgs, types.FunctionArgument{
				Name:         variable.Name.String,
				VariableType: varType,
			})
			variables[variable.Name.String] = structVarType
		}
		genericNames := []string{}
		for _, generic := range generics {
			genericNames = append(genericNames, generic.String)
		}
		maybeStruc, resolutionErr := binding.GetTypeByTypeName(localUniverse, structName.String, genericNames)
		if resolutionErr != nil {
			return nil, nil, TypecheckErrorFromResolutionError(structName.Node, resolutionErr)
		}
		struc, ok := maybeStruc.(*types.Struct)
		if !ok {
			return nil, nil, type_error.PtrOnNodef(structName.Node, "expected struct type in validateStructs")
		}
		struc.Fields = variables

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
		universe, err = binding.CopyOverridingType(universe, structName.String, struc)
		universe, err = binding.CopyAddingVariable(universe, structName, constructorVarType)
		constructors[structName.String] = constructorVarType
	}
	return constructors, universe, nil
}

func TypecheckDeclarations(expectedTypes *map[string]types.VariableType, node parser.Node, declarations []parser.Declaration, universe binding.Universe) (map[string]ast.Expression, *type_error.TypecheckError) {
	typesByName := map[parser.Name]types.VariableType{}

	for _, declaration := range declarations {
		if expectedTypes != nil {
			typesByName[declaration.Name] = (*expectedTypes)[declaration.Name.String]
		}
		if typesByName[declaration.Name] == nil {
			varType, err := typeOfExpressionBox(declaration.ExpressionBox, universe)
			if err != nil {
				return nil, err
			}
			typesByName[declaration.Name] = varType
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

	for _, declaration := range declarations {
		expectedType := typesByName[declaration.Name]
		if expectedType == nil {
			panic("nil expectedType on TypecheckDeclarations")
		}
		astExp, err := expectTypeOfExpressionBox(expectedType, declaration.ExpressionBox, universe)
		if err != nil {
			return nil, err
		}
		result[declaration.Name.String] = astExp
	}

	return result, nil
}

func validateTypeAnnotationInUniverse(typeAnnotation parser.TypeAnnotation, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	switch len(typeAnnotation.OrTypes) {
	case 0:
		return nil, type_error.PtrOnNodef(typeAnnotation.Node, "unexpected error validateTypeAnnotationInUniverse no types found")
	case 1:
		elem := typeAnnotation.OrTypes[0]
		return validateTypeAnnotationElementInUniverse(elem, universe)
	default:
		elements := []types.VariableType{}
		for _, element := range typeAnnotation.OrTypes {
			newElement, err := validateTypeAnnotationElementInUniverse(element, universe)
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

func validateTypeAnnotationElementInUniverse(typeAnnotationElement parser.TypeAnnotationElement, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	var varType types.VariableType
	var err *type_error.TypecheckError
	parser.TypeAnnotationElementExhaustiveSwitch(
		typeAnnotationElement,
		func(typeAnnotation parser.SingleNameType) {
			genericNames := []string{}
			for _, generic := range typeAnnotation.Generics {
				genericNames = append(genericNames, generic.String)
			}
			varType2, err2 := binding.GetTypeByTypeName(universe, typeAnnotation.TypeName.String, genericNames)
			varType = varType2
			err = TypecheckErrorFromResolutionError(typeAnnotation.TypeName.Node, err2)
		},
		func(typeAnnotation parser.FunctionType) {
			localUniverse := universe
			for _, generic := range typeAnnotation.Generics {
				localUniverse, err = binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return
				}
			}
			arguments := []types.FunctionArgument{}
			for _, argAnnotatedType := range typeAnnotation.Arguments {
				varType, err = validateTypeAnnotationInUniverse(argAnnotatedType, localUniverse)
				if err != nil {
					return
				}
				arguments = append(arguments, types.FunctionArgument{
					Name:         "?",
					VariableType: varType,
				})
			}
			var returnType types.VariableType
			returnType, err = validateTypeAnnotationInUniverse(typeAnnotation.ReturnType, localUniverse)
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
