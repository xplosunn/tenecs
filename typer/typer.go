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
	universe, err := resolveImports(imports, StdLib)
	if err != nil {
		return program, err
	}
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
				if parserDec.Name == programDeclaration.Name {
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
		caseDeclaration, caseInterface, caseStruct := topLevelDeclaration.TopLevelDeclarationCases()
		if caseDeclaration != nil {
			declarations = append(declarations, *caseDeclaration)
		} else if caseInterface != nil {
			interfaces = append(interfaces, *caseInterface)
		} else if caseStruct != nil {
			structs = append(structs, *caseStruct)
		} else {
			panic("code on topLevelDeclaration")
		}
	}
	return declarations, interfaces, structs
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

func resolveImports(nodes []parser.Import, stdLib Package) (binding.Universe, *type_error.TypecheckError) {
	universe := binding.NewFromDefaults(DefaultTypesAvailableWithoutImport)
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
	caseSingleNameType, caseFunctionType := typeAnnotation.TypeAnnotationCases()
	if caseSingleNameType != nil {
		varType, ok := binding.GetTypeByTypeName(universe, caseSingleNameType.TypeName)
		if !ok {
			return nil, type_error.PtrTypeCheckErrorf("not found type: %s", caseSingleNameType.TypeName)
		}
		return varType, nil
	} else if caseFunctionType != nil {
		localUniverse := universe
		for _, generic := range caseFunctionType.Generics {
			u, err := binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic})
			if err != nil {
				return nil, err
			}
			localUniverse = u
		}
		arguments := []types.FunctionArgument{}
		for _, argAnnotatedType := range caseFunctionType.Arguments {
			varType, err := validateTypeAnnotationInUniverse(argAnnotatedType, localUniverse)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, types.FunctionArgument{
				Name:         "?",
				VariableType: varType,
			})
		}
		returnType, err := validateTypeAnnotationInUniverse(caseFunctionType.ReturnType, localUniverse)
		if err != nil {
			return nil, err
		}
		return &types.Function{
			Generics:   caseFunctionType.Generics,
			Arguments:  arguments,
			ReturnType: returnType,
		}, nil
	} else {
		panic("cases on typeAnnotation")
	}
}

func validateStructs(nodes []parser.Struct, pkg parser.Package, universe binding.Universe) (map[string]*types.Function, binding.Universe, *type_error.TypecheckError) {
	constructors := map[string]*types.Function{}
	var err *type_error.TypecheckError
	for _, node := range nodes {
		universe, err = binding.CopyAddingType(universe, node.Name, &types.Struct{
			Package: pkg.Identifier,
			Name:    node.Name,
		})
		if err != nil {
			return nil, nil, err
		}
	}
	for _, node := range nodes {
		structName, generics, parserVariables := parser.StructFields(node)
		localUniverse := universe
		for _, generic := range generics {
			u, err := binding.CopyAddingType(localUniverse, generic, &types.TypeArgument{Name: generic})
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
				return nil, nil, type_error.PtrTypeCheckErrorf("%s (are you using an incomparable type?)", err.Error())
			}
			structVarType, ok := types.StructFieldVariableTypeFromVariableType(varType)
			if !ok {
				return nil, nil, type_error.PtrTypeCheckErrorf("not a valid struct var type %s", printableName(varType))
			}
			constructorArgs = append(constructorArgs, types.FunctionArgument{
				Name:         variable.Name,
				VariableType: varType,
			})
			variables[variable.Name] = structVarType
		}
		maybeStruc, ok := binding.GetTypeByTypeName(universe, structName)
		if !ok {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected to find type name in validateStructs")
		}
		struc, ok := maybeStruc.(*types.Struct)
		if !ok {
			return nil, nil, type_error.PtrTypeCheckErrorf("expected struct type in validateStructs")
		}
		struc.Fields = variables

		constructorVarType := &types.Function{
			Generics:   generics,
			Arguments:  constructorArgs,
			ReturnType: struc,
		}
		universe, err = binding.CopyAddingVariable(universe, structName, constructorVarType)
		constructors[structName] = constructorVarType
	}
	return constructors, universe, nil
}

func validateInterfaces(nodes []parser.Interface, pkg parser.Package, universe binding.Universe) (binding.Universe, *type_error.TypecheckError) {
	updatedUniverse := universe
	var err *type_error.TypecheckError
	for _, node := range nodes {
		variables := map[string]types.VariableType{}
		for _, variable := range node.Variables {
			variables[variable.Name] = nil
		}
		updatedUniverse, err = binding.CopyAddingType(updatedUniverse, node.Name, &types.Interface{
			Package:   pkg.Identifier,
			Name:      node.Name,
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
			_, ok := variables[variable.Name]
			if ok {
				return updatedUniverse, type_error.PtrTypeCheckErrorf("more than one variable with name '%s'", variable.Name)
			}
			variables[variable.Name] = varType
		}
		maybeInterf, ok := binding.GetTypeByTypeName(updatedUniverse, name)
		if !ok {
			return nil, type_error.PtrTypeCheckErrorf("expected to find type name in validateInterfaces")
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
		expressions[declaration.Name] = expression
		universe = u
		universe, err = binding.CopyAddingVariable(universe, declaration.Name, ast.VariableTypeOfExpression(expression))
	}

	return expressions, universe, nil
}

func printableNameOfTypeAnnotation(typeAnnotation parser.TypeAnnotation) string {
	caseSingleNameType, caseFunctionType := typeAnnotation.TypeAnnotationCases()
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
		panic("code on typeAnnotation")
	}
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
