package typer

import (
	"errors"
	"fmt"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/expect_type"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/type_of"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/slices"
	"unicode"
)

// TODO FIXME remove hardcoded file name
func TypecheckSingleFile(parsed parser.FileTopLevel) (*ast.Program, error) {
	return TypecheckSinglePackage(map[string]parser.FileTopLevel{"file.10x": parsed})
}

func TypecheckPackages(parsed map[string]parser.FileTopLevel) (*ast.Program, error) {
	byPackage := map[string]map[string]parser.FileTopLevel{}
	for file, parsedFile := range parsed {
		pkg := ""
		for i, name := range parsedFile.Package.DotSeparatedNames {
			if i > 0 {
				pkg += "."
			}
			pkg += name.String
		}
		if byPackage[pkg] == nil {
			byPackage[pkg] = map[string]parser.FileTopLevel{}
		}
		byPackage[pkg][file] = parsedFile
	}
	program := ast.Program{
		Declarations:    map[ast.Ref]ast.Expression{},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{},
		FieldsByType:    map[ast.Ref]map[string]types.VariableType{},
	}
	for _, parsedPkg := range byPackage {
		pkgProgram, err := TypecheckSinglePackage(parsedPkg)
		if err != nil {
			return nil, err
		}
		for ref, expression := range pkgProgram.Declarations {
			program.Declarations[ref] = expression
		}
		for ref, function := range pkgProgram.StructFunctions {
			program.StructFunctions[ref] = function
		}
		for ref, function := range pkgProgram.NativeFunctions {
			program.NativeFunctions[ref] = function
		}
		for ref, fields := range pkgProgram.FieldsByType {
			program.FieldsByType[ref] = fields
		}
	}
	return &program, nil
}

func TypecheckSinglePackage(parsedPackage map[string]parser.FileTopLevel) (*ast.Program, error) {
	if len(parsedPackage) == 0 {
		return nil, errors.New("no files provided for typechecking")
	}
	pkgName := ""
	for _, parsed := range parsedPackage {
		pkgNameInThisFile := ""
		for i, name := range parsed.Package.DotSeparatedNames {
			if i > 0 {
				pkgNameInThisFile += "."
			}
			pkgNameInThisFile += name.String
		}
		if pkgName == "" {
			pkgName = pkgNameInThisFile
		} else if pkgName != pkgNameInThisFile {
			panic("typecheck package should be called with files on same package")
		}

	}
	for k, parsed := range parsedPackage {
		desugared, err := DesugarFileTopLevel(k, parsed)
		if err != nil {
			return nil, err
		}
		parsedPackage[k] = desugared
	}

	for file, topLevel := range parsedPackage {
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
		err := validatePackage(topLevel.Package, file)
		if err != nil {
			return nil, err
		}
	}

	scope := binding.NewFromDefaults(standard_library.DefaultTypesAvailableWithoutImport)
	scope, err := addAllStructFieldsToScope("", scope, standard_library.StdLib)
	if err != nil {
		return nil, err
	}

	program := ast.Program{
		Declarations:    map[ast.Ref]ast.Expression{},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{},
		FieldsByType:    map[ast.Ref]map[string]types.VariableType{},
	}
	for file, fileTopLevel := range parsedPackage {
		programNativeFunctions, programNativeFunctionPackages, u, err := resolveImports(fileTopLevel.Imports, standard_library.StdLib, file, scope)
		if err != nil {
			return nil, err
		}
		scope = u
		for functionName, function := range programNativeFunctions {
			program.NativeFunctions[ast.Ref{
				Package: programNativeFunctionPackages[functionName],
				Name:    functionName,
			}] = function
		}
	}

	structsPerFile := map[string][]parser.Struct{}
	declarationsPerFile := map[string][]parser.Declaration{}
	typeAliasesInAllFiles := map[string][]parser.TypeAlias{}
	for file, fileTopLevel := range parsedPackage {
		declarations, structs, typeAliases := splitTopLevelDeclarations(fileTopLevel.TopLevelDeclarations)
		structsPerFile[file] = structs
		declarationsPerFile[file] = declarations
		typeAliasesInAllFiles[file] = typeAliases
	}

	programStructFunctions, scope, err := validateStructs(structsPerFile, pkgName, scope)
	if err != nil {
		return nil, err
	}
	program.StructFunctions = programStructFunctions
	programFieldsByType := binding.GetAllFields(scope)
	for name, fieldsMap := range programFieldsByType {
		program.FieldsByType[ast.Ref{
			Package: pkgName,
			Name:    name,
		}] = fieldsMap
	}

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
				return nil, type_error.FromScopeCheckError(file, err)
			}

			u, err2 := binding.CopyAddingTypeAliasToAllFiles(scope, name, genericNameStrings, varType)
			if err2 != nil {
				return nil, type_error.FromResolutionError(file, name.Node, err2)
			}
			scope = u
		}
	}

	declarationsMap, err := TypecheckDeclarations(pkgName, parser.Node{}, declarationsPerFile, scope)
	if err != nil {
		return nil, err
	}
	program.Declarations = map[ast.Ref]ast.Expression{}
	for varName, varExp := range declarationsMap {
		program.Declarations[ast.Ref{
			Package: pkgName,
			Name:    varName,
		}] = varExp
	}

	return &program, nil
}

func validatePackage(node parser.Package, file string) *type_error.TypecheckError {
	for _, name := range node.DotSeparatedNames {
		if !unicode.IsLower(rune(name.String[0])) {
			return type_error.PtrOnNodef(file, name.Node, "package name should start with a lowercase letter")
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

func addAllStructFieldsToScope(file string, scope binding.Scope, pkg standard_library.Package) (binding.Scope, *type_error.TypecheckError) {
	for structName, structWithFields := range pkg.Structs {
		var err *binding.ResolutionError
		scope, err = binding.CopyAddingFields(scope, structWithFields.Struct.Package, parser.Name{
			String: structName,
		}, structWithFields.Fields)
		if err != nil {
			// TODO FIXME shouldn't convert with an empty Node
			return nil, type_error.FromResolutionError(file, parser.Node{}, err)
		}
	}
	for _, nestedPkg := range pkg.Packages {
		var err *type_error.TypecheckError
		scope, err = addAllStructFieldsToScope(file, scope, nestedPkg)
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
			return nil, nil, nil, type_error.PtrOnNodef(file, errNode, "all interfaces belong to a package")
		}
		currPackage := stdLib
		currPackageName := ""
		for i, name := range dotSeparatedNames {
			if i < len(dotSeparatedNames)-1 {
				p, ok := currPackage.Packages[name.String]
				if !ok {
					return nil, nil, nil, type_error.PtrOnNodef(file, name.Node, "no package "+name.String+" found")
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
					return nil, nil, nil, type_error.FromResolutionError(file, fallbackOnNil(as, name).Node, err)
				}
				updatedScope, err = binding.CopyAddingFields(updatedScope, currPackageName, fallbackOnNil(as, name), struc.Fields)
				if err != nil {
					return nil, nil, nil, type_error.FromResolutionError(file, fallbackOnNil(as, name).Node, err)
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
					updatedScope, err = binding.CopyAddingFileVariable(updatedScope, struc.Struct.Package, file, *as, &name, constructorVarType)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(file, as.Node, err)
					}
				} else {
					updatedScope, err = binding.CopyAddingFileVariable(updatedScope, struc.Struct.Package, file, name, nil, constructorVarType)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(file, name.Node, err)
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
					updatedScope, err := binding.CopyAddingFileVariable(scope, currPackageName, file, *as, &name, varTypeToImport)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(file, as.Node, err)
					}
					scope = updatedScope
				} else {
					updatedScope, err := binding.CopyAddingFileVariable(scope, currPackageName, file, name, nil, varTypeToImport)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(file, name.Node, err)
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

			return nil, nil, nil, type_error.PtrOnNodef(file, name.Node, "didn't find "+name.String+" while importing")
		}
	}
	return nativeFunctions, nativeFunctionPackages, scope, nil
}

func validateStructs(structsPerFile map[string][]parser.Struct, pkgName string, scope binding.Scope) (map[ast.Ref]*types.Function, binding.Scope, *type_error.TypecheckError) {
	constructors := map[ast.Ref]*types.Function{}
	for file, structsInFile := range structsPerFile {
		for _, node := range structsInFile {
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
				return nil, nil, type_error.FromResolutionError(file, node.Name.Node, err)
			}
		}
	}
	for file, structsInFile := range structsPerFile {
		for _, node := range structsInFile {
			structName, generics, parserVariables := parser.StructFields(node)
			localScope := scope
			for _, generic := range generics {
				u, err := binding.CopyAddingTypeToAllFiles(localScope, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return nil, nil, type_error.FromResolutionError(file, generic.Node, err)
				}
				localScope = u
			}
			constructorArgs := []types.FunctionArgument{}
			variables := map[string]types.VariableType{}
			for _, variable := range parserVariables {

				varType, err := scopecheck.ValidateTypeAnnotationInScope(variable.Type, file, localScope)
				if err != nil {
					return nil, nil, type_error.FromScopeCheckError(file, err)
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
				return nil, nil, type_error.FromResolutionError(file, structName.Node, err)
			}

			genericNames := []types.VariableType{}
			for _, generic := range generics {
				genericNames = append(genericNames, &types.TypeArgument{
					Name: generic.String,
				})
			}
			maybeStruc, resolutionErr := binding.GetTypeByTypeName(localScope, "", structName.String, genericNames)
			if resolutionErr != nil {
				return nil, nil, type_error.FromResolutionError(file, structName.Node, resolutionErr)
			}
			struc, ok := maybeStruc.(*types.KnownType)
			if !ok {
				return nil, nil, type_error.PtrOnNodef(file, structName.Node, "expected struct type in validateStructs")
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
			scope, err = binding.CopyAddingPackageVariable(scope, pkgName, structName, constructorVarType)
			constructors[ast.Ref{
				Package: pkgName,
				Name:    structName.String,
			}] = constructorVarType
		}
	}
	return constructors, scope, nil
}

func TypecheckDeclarations(pkg string, node parser.Node, declarationsPerFileWithUnderscores map[string][]parser.Declaration, scope binding.Scope) (map[string]ast.Expression, *type_error.TypecheckError) {
	declarationsPerFile := map[string][]parser.Declaration{}
	syntheticNameIterator := 0
	for file, declarations := range declarationsPerFileWithUnderscores {
		declarationsPerFile[file] = []parser.Declaration{}
		for _, declaration := range declarations {
			if declaration.Name.String == "_" {
				declaration.Name.String = fmt.Sprintf("syntheticName_%d", syntheticNameIterator)
				syntheticNameIterator += 1
			}
			declarationsPerFile[file] = append(declarationsPerFile[file], declaration)
		}
	}

	typesByName := map[parser.Name]types.VariableType{}
	filesByName := map[parser.Name]string{}

	for file, declarations := range declarationsPerFile {
		for _, declaration := range declarations {
			if slices.Contains(expect_type.ForbiddenVariableNames, declaration.Name.String) {
				return nil, type_error.PtrOnNodef(file, declaration.Name.Node, "Variable can't be named '%s'", declaration.Name.String)
			}
			if declaration.TypeAnnotation != nil {
				annotatedVarType, err := scopecheck.ValidateTypeAnnotationInScope(*declaration.TypeAnnotation, file, scope)
				if err != nil {
					return nil, type_error.FromScopeCheckError(file, err)
				}
				if typesByName[declaration.Name] == nil {
					typesByName[declaration.Name] = annotatedVarType
					filesByName[declaration.Name] = file
				} else if !types.VariableTypeEq(typesByName[declaration.Name], annotatedVarType) {
					return nil, type_error.PtrOnNodef(file, node, "annotated type %s doesn't match the expected %s", types.PrintableName(annotatedVarType), types.PrintableName(typesByName[declaration.Name]))
				}
			}
			if typesByName[declaration.Name] == nil {
				varType, err := type_of.TypeOfExpressionBox(declaration.ExpressionBox, file, scope)
				if err != nil {
					return nil, err
				}
				typesByName[declaration.Name] = varType
				filesByName[declaration.Name] = file
			}
		}
	}

	for varName, varType := range typesByName {
		var err *binding.ResolutionError
		scope, err = binding.CopyAddingPackageVariable(scope, pkg, varName, varType)
		if err != nil {
			return nil, type_error.FromResolutionError(filesByName[varName], varName.Node, err)
		}
	}

	result := map[string]ast.Expression{}

	for file, declarations := range declarationsPerFile {
		for _, declaration := range declarations {
			expectedType := typesByName[declaration.Name]
			if expectedType == nil {
				panic("nil expectedType on TypecheckDeclarations")
			}
			astExp, err := expect_type.ExpectTypeOfExpressionBox(expectedType, declaration.ExpressionBox, file, scope)
			if err != nil {
				return nil, err
			}
			result[declaration.Name.String] = astExp
		}
	}

	return result, nil
}
