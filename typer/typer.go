package typer

import (
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/slices"
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
	for k, parsed := range parsedPackage {
		desugared, err := DesugarFileTopLevel(parsed)
		if err != nil {
			return nil, err
		}
		parsedPackage[k] = desugared
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

	scope := binding.NewFromDefaults(standard_library.DefaultTypesAvailableWithoutImport)
	scope, err := addAllStructFieldsToScope(scope, standard_library.StdLib)
	if err != nil {
		return nil, err
	}

	program := ast.Program{
		Package:                pkgName,
		NativeFunctions:        map[string]*types.Function{},
		NativeFunctionPackages: map[string]string{},
	}
	for file, fileTopLevel := range parsedPackage {
		programNativeFunctions, programNativeFunctionPackages, u, err := resolveImports(fileTopLevel.Imports, standard_library.StdLib, file, scope)
		if err != nil {
			return nil, err
		}
		scope = u
		for functionName, function := range programNativeFunctions {
			if program.NativeFunctions[functionName] != nil && program.NativeFunctionPackages[functionName] != programNativeFunctionPackages[functionName] {
				return nil, type_error.PtrOnNodef(fileTopLevel.Package.DotSeparatedNames[0].Node, "TODO: unsupported imports of different functions from standard library with same name on different files of same package")
			}
			program.NativeFunctions[functionName] = function
			program.NativeFunctionPackages[functionName] = programNativeFunctionPackages[functionName]
		}
	}

	structsInAllFiles := []parser.Struct{}
	declarationsPerFile := map[string][]parser.Declaration{}
	typeAliasesInAllFiles := map[string][]parser.TypeAlias{}
	for file, fileTopLevel := range parsedPackage {
		declarations, structs, typeAliases := splitTopLevelDeclarations(fileTopLevel.TopLevelDeclarations)
		structsInAllFiles = append(structsInAllFiles, structs...)
		declarationsPerFile[file] = declarations
		typeAliasesInAllFiles[file] = typeAliases
	}

	programStructFunctions, scope, err := validateStructs(structsInAllFiles, pkgName, scope)
	if err != nil {
		return nil, err
	}
	program.StructFunctions = programStructFunctions
	program.FieldsByType = binding.GetAllFields(scope)

	for file, typeAliases := range typeAliasesInAllFiles {
		for _, typeAlias := range typeAliases {
			name, generics, typ := parser.TypeAliasFields(typeAlias)
			genericNameStrings := []string{}
			scopeOnlyValidForTypeAlias := scope
			for _, generic := range generics {
				genericNameStrings = append(genericNameStrings, generic.String)
				u, err := binding.CopyAddingTypeToFile(scopeOnlyValidForTypeAlias, file, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return nil, err
				}
				scopeOnlyValidForTypeAlias = u
			}

			varType, err := scopecheck.ValidateTypeAnnotationInScope(typ, file, scopeOnlyValidForTypeAlias)
			if err != nil {
				return nil, type_error.FromScopeCheckError(err)
			}

			u, err2 := binding.CopyAddingTypeAliasToAllFiles(scope, name, genericNameStrings, varType)
			if err2 != nil {
				return nil, type_error.FromResolutionError(name.Node, err2)
			}
			scope = u
		}
	}

	declarationsMap, err := TypecheckDeclarations(nil, &pkgName, parser.Node{}, declarationsPerFile, scope)
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
	slices.SortFunc(program.Declarations, func(a *ast.Declaration, b *ast.Declaration) bool {
		return slices.IsSorted([]string{a.Name, b.Name})
	})

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

func splitTopLevelDeclarations(topLevelDeclarations []parser.TopLevelDeclaration) ([]parser.Declaration, []parser.Struct, []parser.TypeAlias) {
	declarations := []parser.Declaration{}
	structs := []parser.Struct{}
	typeAliases := []parser.TypeAlias{}
	for _, topLevelDeclaration := range topLevelDeclarations {
		parser.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration parser.Declaration) {
				declarations = append(declarations, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.Struct) {
				structs = append(structs, topLevelDeclaration)
			},
			func(topLevelDeclaration parser.TypeAlias) {
				typeAliases = append(typeAliases, topLevelDeclaration)
			},
		)
	}
	return declarations, structs, typeAliases
}

func addAllStructFieldsToScope(scope binding.Scope, pkg standard_library.Package) (binding.Scope, *type_error.TypecheckError) {
	for structName, structWithFields := range pkg.Structs {
		var err *binding.ResolutionError
		scope, err = binding.CopyAddingFields(scope, structWithFields.Struct.Package, parser.Name{
			String: structName,
		}, structWithFields.Fields)
		if err != nil {
			// TODO FIXME shouldn't convert with an empty Node
			return nil, type_error.FromResolutionError(parser.Node{}, err)
		}
	}
	for _, nestedPkg := range pkg.Packages {
		var err *type_error.TypecheckError
		scope, err = addAllStructFieldsToScope(scope, nestedPkg)
		if err != nil {
			return nil, err
		}
	}
	return scope, nil
}

func fallbackOnNil[T any](a *T, b T) T {
	if a != nil {
		return *a
	}
	return b
}

func resolveImports(nodes []parser.Import, stdLib standard_library.Package, file string, scope binding.Scope) (map[string]*types.Function, map[string]string, binding.Scope, *type_error.TypecheckError) {
	nativeFunctions := map[string]*types.Function{}
	nativeFunctionPackages := map[string]string{}
	for _, node := range nodes {
		dotSeparatedNames, as := parser.ImportFields(node)
		if len(dotSeparatedNames) < 2 {
			errNode := node.Node
			if len(dotSeparatedNames) > 0 {
				errNode = dotSeparatedNames[0].Node
			}
			return nil, nil, nil, type_error.PtrOnNodef(errNode, "all interfaces belong to a package")
		}
		currPackage := stdLib
		currPackageName := ""
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name.String]
				if !ok {
					return nil, nil, nil, type_error.PtrOnNodef(name.Node, "no package "+name.String+" found")
				}
				currPackage = p
				if i > 0 {
					currPackageName += "."
				}
				currPackageName += name.String
				continue
			}
			struc, ok := currPackage.Structs[name.String]
			if ok {
				updatedScope, err := binding.CopyAddingTypeToFile(scope, file, fallbackOnNil(as, name), struc.Struct)
				if err != nil {
					return nil, nil, nil, type_error.FromResolutionError(fallbackOnNil(as, name).Node, err)
				}
				updatedScope, err = binding.CopyAddingFields(updatedScope, currPackageName, fallbackOnNil(as, name), struc.Fields)
				if err != nil {
					return nil, nil, nil, type_error.FromResolutionError(fallbackOnNil(as, name).Node, err)
				}
				constructorArguments := []types.FunctionArgument{}
				for _, structFieldName := range struc.FieldNamesSorted {
					constructorArguments = append(constructorArguments, types.FunctionArgument{
						Name:         structFieldName,
						VariableType: struc.Fields[structFieldName],
					})
				}
				constructorVarType := &types.Function{
					Generics:   struc.Struct.DeclaredGenerics,
					Arguments:  constructorArguments,
					ReturnType: struc.Struct,
				}
				if as != nil {
					updatedScope, err = binding.CopyAddingPackageVariable(updatedScope, struc.Struct.Package, *as, &name, constructorVarType)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(as.Node, err)
					}
				} else {
					updatedScope, err = binding.CopyAddingPackageVariable(updatedScope, struc.Struct.Package, name, nil, constructorVarType)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(name.Node, err)
					}
				}
				scope = updatedScope
				nativeFunctions[name.String] = constructorVarType
				pkg := ""
				for i, name := range dotSeparatedNames {
					if i < len(dotSeparatedNames)-1 {
						if i > 0 {
							pkg += "_"
						}
						pkg += name.String
					}
				}
				nativeFunctionPackages[name.String] = pkg
				continue
			}
			varTypeToImport, ok := currPackage.Variables[name.String]
			if ok {
				if as != nil {
					updatedScope, err := binding.CopyAddingPackageVariable(scope, currPackageName, *as, &name, varTypeToImport)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(as.Node, err)
					}
					scope = updatedScope
				} else {
					updatedScope, err := binding.CopyAddingPackageVariable(scope, currPackageName, name, nil, varTypeToImport)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(name.Node, err)
					}
					scope = updatedScope
				}
				fn, ok := varTypeToImport.(*types.Function)
				if !ok {
					panic(fmt.Sprintf("todo resolveImports not native function but %T", varTypeToImport))
				}
				nativeFunctions[name.String] = fn
				pkg := ""
				for i, name := range dotSeparatedNames {
					if i < len(dotSeparatedNames)-1 {
						if i > 0 {
							pkg += "_"
						}
						pkg += name.String
					}
				}
				nativeFunctionPackages[name.String] = pkg
				continue
			}

			return nil, nil, nil, type_error.PtrOnNodef(name.Node, "didn't find "+name.String+" while importing")
		}
	}
	return nativeFunctions, nativeFunctionPackages, scope, nil
}

func validateStructs(nodes []parser.Struct, pkgName string, scope binding.Scope) (map[string]*types.Function, binding.Scope, *type_error.TypecheckError) {
	constructors := map[string]*types.Function{}
	for _, node := range nodes {
		var err *binding.ResolutionError
		genericNames := []string{}
		genericTypeArgs := []types.VariableType{}
		for _, generic := range node.Generics {
			genericNames = append(genericNames, generic.String)
			genericTypeArgs = append(genericTypeArgs, &types.TypeArgument{Name: generic.String})
		}
		scope, err = binding.CopyAddingTypeToAllFiles(scope, node.Name, &types.KnownType{
			Package:          pkgName,
			Name:             node.Name.String,
			DeclaredGenerics: genericNames,
			Generics:         genericTypeArgs,
		})
		if err != nil {
			return nil, nil, type_error.FromResolutionError(node.Name.Node, err)
		}
	}
	for _, node := range nodes {
		structName, generics, parserVariables := parser.StructFields(node)
		localScope := scope
		for _, generic := range generics {
			u, err := binding.CopyAddingTypeToAllFiles(localScope, generic, &types.TypeArgument{Name: generic.String})
			if err != nil {
				return nil, nil, type_error.FromResolutionError(generic.Node, err)
			}
			localScope = u
		}
		constructorArgs := []types.FunctionArgument{}
		variables := map[string]types.VariableType{}
		for _, variable := range parserVariables {
			varType, err := scopecheck.ValidateTypeAnnotationInScope(variable.Type, "", localScope)
			if err != nil {
				return nil, nil, type_error.FromScopeCheckError(err)
			}
			constructorArgs = append(constructorArgs, types.FunctionArgument{
				Name:         variable.Name.String,
				VariableType: varType,
			})
			variables[variable.Name.String] = varType
		}
		var err *binding.ResolutionError
		scope, err = binding.CopyAddingFields(scope, pkgName, structName, variables)
		if err != nil {
			return nil, nil, type_error.FromResolutionError(structName.Node, err)
		}

		genericNames := []types.VariableType{}
		for _, generic := range generics {
			genericNames = append(genericNames, &types.TypeArgument{
				Name: generic.String,
			})
		}
		maybeStruc, resolutionErr := binding.GetTypeByTypeName(localScope, "", structName.String, genericNames)
		if resolutionErr != nil {
			return nil, nil, type_error.FromResolutionError(structName.Node, resolutionErr)
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
		scope, err = binding.CopyAddingPackageVariable(scope, pkgName, structName, nil, constructorVarType)
		constructors[structName.String] = constructorVarType
	}
	return constructors, scope, nil
}

func TypecheckDeclarations(expectedTypes *map[string]types.VariableType, pkg *string, node parser.Node, declarationsPerFile map[string][]parser.Declaration, scope binding.Scope) (map[string]ast.Expression, *type_error.TypecheckError) {
	if (expectedTypes == nil) == (pkg == nil) {
		panic("TypecheckDeclarations should have either expectedTypes or pkg")
	}
	typesByName := map[parser.Name]types.VariableType{}

	for file, declarations := range declarationsPerFile {
		for _, declaration := range declarations {
			if expectedTypes != nil {
				typesByName[declaration.Name] = (*expectedTypes)[declaration.Name.String]
			}
			if declaration.TypeAnnotation != nil {
				annotatedVarType, err := scopecheck.ValidateTypeAnnotationInScope(*declaration.TypeAnnotation, file, scope)
				if err != nil {
					return nil, type_error.FromScopeCheckError(err)
				}
				if typesByName[declaration.Name] == nil {
					typesByName[declaration.Name] = annotatedVarType
				} else if !types.VariableTypeEq(typesByName[declaration.Name], annotatedVarType) {
					return nil, type_error.PtrOnNodef(node, "annotated type %s doesn't match the expected %s", types.PrintableName(annotatedVarType), types.PrintableName(typesByName[declaration.Name]))
				}
			}
			if typesByName[declaration.Name] == nil {
				varType, err := typeOfExpressionBox(declaration.ExpressionBox, file, scope)
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
		var err *binding.ResolutionError
		if pkg != nil {
			scope, err = binding.CopyAddingPackageVariable(scope, *pkg, varName, nil, varType)
		} else {
			scope, err = binding.CopyAddingLocalVariable(scope, varName, varType)
		}
		if err != nil {
			return nil, type_error.FromResolutionError(varName.Node, err)
		}
	}

	result := map[string]ast.Expression{}

	for file, declarations := range declarationsPerFile {
		for i, declaration := range declarations {
			expectedType := typesByName[declaration.Name]
			if expectedType == nil {
				panic("nil expectedType on TypecheckDeclarations")
			}
			astExp, err := expectTypeOfExpressionBox(expectedType, declaration.ExpressionBox, file, scope)
			if err != nil {
				return nil, err
			}
			name := declaration.Name.String
			if name == "_" {
				name = fmt.Sprintf("syntheticName_%d", i)
			}
			result[name] = astExp
		}
	}

	return result, nil
}
