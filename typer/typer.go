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
		return program, err
	}
	programNativeFunctions, universe, err := resolveImports(imports, standard_library.StdLib)
	if err != nil {
		return program, err
	}
	program.NativeFunctions = programNativeFunctions
	declarations, interfaces, structs := splitTopLevelDeclarations(topLevelDeclarations)
	programStructFunctions, universe, err := validateStructs(structs, pkg, universe)
	if err != nil {
		return program, err
	}
	program.StructFunctions = programStructFunctions
	universe, err = validateInterfaces(interfaces, pkg, universe)
	if err != nil {
		return program, err
	}
	programDeclarationsMap, universe, err := validateTopLevelDeclarationsWithoutFunctionBlocks(declarations, universe)
	if err != nil {
		return program, err
	}
	programDeclarations := []*ast.Declaration{}
	for varName, varExp := range programDeclarationsMap {
		programDeclarations = append(programDeclarations, &ast.Declaration{
			VariableType: &types.BasicType{Type: "Void"},
			Name:         varName,
			Expression:   varExp,
		})
	}
	program.Declarations = programDeclarations

	for _, programDeclaration := range programDeclarations {
		caseModule, caseLiteralExp, caseReferenceAndMaybeInvocation, caseWithAccessAndMaybeInvocation, caseLambda, caseDeclaration, caseIf := programDeclaration.Expression.ExpressionCases()
		_ = caseModule
		_ = caseLiteralExp
		_ = caseReferenceAndMaybeInvocation
		_ = caseWithAccessAndMaybeInvocation
		_ = caseDeclaration
		_ = caseIf
		if caseLambda != nil {
			var parserExpBox parser.ExpressionBox
			for _, parserDec := range declarations {
				if parserDec.Name.String == programDeclaration.Name {
					parserExpBox = parserDec.ExpressionBox
					break
				}
			}
			_, exp, err := expectTypeOfExpressionBox(true, parserExpBox, caseLambda.VariableType, universe)
			if err != nil {
				return nil, err
			}
			caseLambda.Block = exp.(ast.Function).Block
			programDeclaration.Expression = caseLambda
		}
	}

	return program, nil
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

func resolveImports(nodes []parser.Import, stdLib standard_library.Package) (map[string]*types.Function, binding.Universe, *type_error.TypecheckError) {
	universe := binding.NewFromDefaults(standard_library.DefaultTypesAvailableWithoutImport)
	nativeFunctions := map[string]*types.Function{}
	for _, node := range nodes {
		dotSeparatedNames := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			errNode := node.Node
			if len(dotSeparatedNames) > 0 {
				errNode = dotSeparatedNames[0].Node
			}
			return nil, nil, type_error.PtrOnNodef(errNode, "all interfaces belong to a package")
		}
		currPackage := stdLib
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name.String]
				if !ok {
					return nil, nil, type_error.PtrOnNodef(name.Node, "no package "+name.String+" found")
				}
				currPackage = p
				continue
			}
			interf, ok := currPackage.Interfaces[name.String]
			if ok {
				updatedUniverse, err := binding.CopyAddingType(universe, name, interf)
				if err != nil {
					return nil, nil, err
				}
				universe = updatedUniverse
				continue
			}
			varTypeToImport, ok := currPackage.Variables[name.String]
			if ok {
				updatedUniverse, err := binding.CopyAddingVariable(universe, name, varTypeToImport)
				if err != nil {
					return nil, nil, err
				}
				universe = updatedUniverse
				fn, ok := varTypeToImport.(*types.Function)
				if !ok {
					panic(fmt.Sprintf("todo resolveImports not native function but %T", varTypeToImport))
				}
				nativeFunctions[name.String] = fn
				continue
			}

			return nil, nil, type_error.PtrOnNodef(name.Node, "didn't find "+name.String+" while importing")
		}
	}
	return nativeFunctions, universe, nil
}

func validateTypeAnnotationInUniverse(typeAnnotation parser.TypeAnnotation, universe binding.Universe) (types.VariableType, *type_error.TypecheckError) {
	var varType types.VariableType
	var err *type_error.TypecheckError
	parser.TypeAnnotationExhaustiveSwitch(
		typeAnnotation,
		func(typeAnnotation parser.SingleNameType) {
			var ok bool
			varType, ok = binding.GetTypeByTypeName(universe, typeAnnotation.TypeName.String)
			if !ok {
				err = type_error.PtrOnNodef(typeAnnotation.TypeName.Node, "not found type: %s", typeAnnotation.TypeName.String)
			}
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

func validateStructs(nodes []parser.Struct, pkg parser.Package, universe binding.Universe) (map[string]*types.Function, binding.Universe, *type_error.TypecheckError) {
	constructors := map[string]*types.Function{}
	var err *type_error.TypecheckError
	for _, node := range nodes {
		universe, err = binding.CopyAddingType(universe, node.Name, &types.Struct{
			Package: pkg.Identifier.String,
			Name:    node.Name.String,
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
		maybeStruc, ok := binding.GetTypeByTypeName(universe, structName.String)
		if !ok {
			return nil, nil, type_error.PtrOnNodef(structName.Node, "expected to find type name in validateStructs")
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
		universe, err = binding.CopyAddingVariable(universe, structName, constructorVarType)
		constructors[structName.String] = constructorVarType
	}
	return constructors, universe, nil
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
		maybeInterf, ok := binding.GetTypeByTypeName(updatedUniverse, name.String)
		if !ok {
			return nil, type_error.PtrOnNodef(node.Name.Node, "expected to find type name in validateInterfaces")
		}
		interf, ok := maybeInterf.(*types.Interface)
		if !ok {
			panic("expected interface type in validateInterfaces")
		}
		interf.Variables = variables
	}
	return updatedUniverse, nil
}

func validateTopLevelDeclarationsWithoutFunctionBlocks(parserDeclarations []parser.Declaration, universe binding.Universe) (map[string]ast.Expression, binding.Universe, *type_error.TypecheckError) {
	expressions := map[string]ast.Expression{}

	for _, declaration := range parserDeclarations {
		u, expression, err := determineTypeOfExpressionBox(false, declaration.ExpressionBox, universe)
		if err != nil {
			return nil, nil, err
		}
		expressions[declaration.Name.String] = expression
		universe = u
		universe, err = binding.CopyAddingVariable(universe, declaration.Name, ast.VariableTypeOfExpression(expression))
	}

	return expressions, universe, nil
}

func printableNameOfTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	var result string
	parser.TypeAnnotationExhaustiveSwitch(
		typeAnnotation,
		func(typeAnnotation parser.SingleNameType) {
			result = typeAnnotation.TypeName.String
		},
		func(typeAnnotation parser.FunctionType) {
			result = "("
			for i, argument := range typeAnnotation.Arguments {
				if i > 0 {
					result += ", "
				}
				result += printableNameOfTypeAnnotation(argument)
			}
			result = result + ") -> " + printableNameOfTypeAnnotation(typeAnnotation.ReturnType)
		},
	)
	return result
}

func printableName(varType types.VariableType) string {
	caseTypeArgument, caseStruct, caseInterface, caseFunction, caseBasicType, caseVoid := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		return "<" + caseTypeArgument.Name + ">"
	} else if caseStruct != nil {
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
		panic(fmt.Errorf("code on %v", varType))
	}
}
