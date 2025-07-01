package typer

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/xplosunn/tenecs/desugar"
	"github.com/xplosunn/tenecs/typer/ast"
	"github.com/xplosunn/tenecs/typer/async"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/dependency"
	"github.com/xplosunn/tenecs/typer/expect_type"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/standard_library"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/type_of"
	"github.com/xplosunn/tenecs/typer/types"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

// TODO FIXME remove hardcoded file name
func TypecheckSingleFile(parsed desugar.FileTopLevel) (*ast.Program, error) {
	return TypecheckSinglePackage(map[string]desugar.FileTopLevel{"file.10x": parsed}, nil)
}

func TypecheckPackages(parsed map[string]desugar.FileTopLevel) (*ast.Program, error) {
	byPackage := map[string]map[string]desugar.FileTopLevel{}
	for file, parsedFile := range parsed {
		pkg := ""
		for i, name := range parsedFile.Package.DotSeparatedNames {
			if i > 0 {
				pkg += "."
			}
			pkg += name.String
		}
		if byPackage[pkg] == nil {
			byPackage[pkg] = map[string]desugar.FileTopLevel{}
		}
		byPackage[pkg][file] = parsedFile
	}
	program := ast.Program{
		Declarations:                  map[ast.Ref]ast.Expression{},
		TypeAliases:                   map[ast.Ref]ast.TypeAlias{},
		StructFunctions:               map[ast.Ref]*types.Function{},
		NativeFunctions:               map[ast.Ref]*types.Function{},
		FieldsByType:                  map[ast.Ref]map[string]types.VariableType{},
		StructTypeArgumentMatchFields: map[ast.Ref][]string{},
	}
	typedPackages := []string{}
	for len(maps.Keys(byPackage)) > len(typedPackages) {
		type PackageProgram struct {
			Program *ast.Program
			Package string
		}
		typedPackagesInThisLoop := []async.Async[PackageProgram]{}
		for pkgName, parsedPkg := range byPackage {
			if slices.Contains(typedPackages, pkgName) {
				continue
			}

			dependencies := dependency.DependenciesOfSinglePackage(parsedPkg)
			allDependenciesTyped := true
			for _, dep := range dependencies {
				if !slices.Contains(typedPackages, dep) {
					allDependenciesTyped = false
					break
				}
			}
			if !allDependenciesTyped {
				continue
			}
			otherPackageDeclarations := map[ast.Ref]types.VariableType{}
			for ref, expression := range program.Declarations {
				otherPackageDeclarations[ref] = ast.VariableTypeOfExpression(expression)
			}
			otherPackageTypeAliases := map[ast.Ref]ast.TypeAlias{}
			for ref, typeAlias := range program.TypeAliases {
				otherPackageTypeAliases[ref] = typeAlias
			}
			pkg := pkgName
			typedPackagesInThisLoop = append(typedPackagesInThisLoop, async.RunAsync(func() (PackageProgram, error) {
				program, err := TypecheckSinglePackage(parsedPkg, &OtherPackagesContext{
					Declarations:    otherPackageDeclarations,
					TypeAliases:     otherPackageTypeAliases,
					StructFunctions: program.StructFunctions,
					FieldsByType:    program.FieldsByType,
				})
				return PackageProgram{
					Program: program,
					Package: pkg,
				}, err
			}))
		}
		for _, async := range typedPackagesInThisLoop {
			packageProgramWrapper, err := async.Await()
			if err != nil {
				return nil, err
			}
			pkgProgram := packageProgramWrapper.Program
			for ref, expression := range pkgProgram.Declarations {
				program.Declarations[ref] = expression
			}
			for ref, typeAlias := range pkgProgram.TypeAliases {
				program.TypeAliases[ref] = typeAlias
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
			for ref, fields := range pkgProgram.StructTypeArgumentMatchFields {
				program.StructTypeArgumentMatchFields[ref] = fields
			}
			typedPackages = append(typedPackages, packageProgramWrapper.Package)
		}
		if len(typedPackagesInThisLoop) == 0 {
			panic("circular dependencies detected (todo nicer error here)")
		}
	}
	return &program, nil
}

type OtherPackagesContext struct {
	Declarations    map[ast.Ref]types.VariableType
	TypeAliases     map[ast.Ref]ast.TypeAlias
	StructFunctions map[ast.Ref]*types.Function
	FieldsByType    map[ast.Ref]map[string]types.VariableType
}

func TypecheckSinglePackage(parsedPackage map[string]desugar.FileTopLevel, otherPackagesContext *OtherPackagesContext) (*ast.Program, error) {
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
		desugared, err := desugar.DesugarFileTopLevel(k, parsed)
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
	scope = addAllStructFieldsToScope("", scope, standard_library.StdLib)

	program := ast.Program{
		Declarations:    map[ast.Ref]ast.Expression{},
		StructFunctions: map[ast.Ref]*types.Function{},
		NativeFunctions: map[ast.Ref]*types.Function{},
		FieldsByType:    map[ast.Ref]map[string]types.VariableType{},
	}
	for file, fileTopLevel := range parsedPackage {
		programNativeFunctions, programNativeFunctionPackages, u, err := resolveImports(fileTopLevel.Imports, standard_library.StdLib, otherPackagesContext, file, scope)
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

	structsPerFile := map[string][]desugar.Struct{}
	declarationsPerFile := map[string][]desugar.Declaration{}
	typeAliasesInAllFiles := map[string][]desugar.TypeAlias{}
	for file, fileTopLevel := range parsedPackage {
		declarations, structs, typeAliases := splitTopLevelDeclarations(fileTopLevel.TopLevelDeclarations)
		structsPerFile[file] = structs
		declarationsPerFile[file] = declarations
		typeAliasesInAllFiles[file] = typeAliases
	}

	programStructFunctions, programTypeAliases, scope, err := validateStructsAndTypeAliases(structsPerFile, typeAliasesInAllFiles, pkgName, scope)
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

	declarationsMap, err := TypecheckDeclarations(pkgName, desugar.Node{}, declarationsPerFile, scope)
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

	program.TypeAliases = map[ast.Ref]ast.TypeAlias{}
	for ref, typeAlias := range programTypeAliases {
		program.TypeAliases[ref] = typeAlias
	}

	program.StructTypeArgumentMatchFields = map[ast.Ref][]string{}
	for _, function := range program.StructFunctions {
		resolveStructFields := map[binding.Ref]map[string]types.VariableType{}
		for ref, fields := range program.FieldsByType {
			resolveStructFields[binding.Ref{
				Package: ref.Package,
				Name:    ref.Name,
			}] = fields
		}
		_, _, caseKnownType, _, _ := function.ReturnType.VariableTypeCases()
		if caseKnownType == nil {
			panic("expected struct as return of struct constructor function")
		}
		genericsMatchableByField, err := expect_type.KnownTypeGenericsMatchByField(*caseKnownType, resolveStructFields)
		if err == nil {
			program.StructTypeArgumentMatchFields[ast.Ref{
				Package: caseKnownType.Package,
				Name:    caseKnownType.Name,
			}] = genericsMatchableByField
		}

	}

	return &program, nil
}

func validatePackage(node desugar.Package, file string) *type_error.TypecheckError {
	for _, name := range node.DotSeparatedNames {
		if !unicode.IsLower(rune(name.String[0])) {
			return type_error.PtrOnNodef(file, name.Node, "package name should start with a lowercase letter")
		}
	}
	return nil
}

func splitTopLevelDeclarations(topLevelDeclarations []desugar.TopLevelDeclaration) ([]desugar.Declaration, []desugar.Struct, []desugar.TypeAlias) {
	declarations := []desugar.Declaration{}
	structs := []desugar.Struct{}
	typeAliases := []desugar.TypeAlias{}
	for _, topLevelDeclaration := range topLevelDeclarations {
		desugar.TopLevelDeclarationExhaustiveSwitch(
			topLevelDeclaration,
			func(topLevelDeclaration desugar.Declaration) {
				declarations = append(declarations, topLevelDeclaration)
			},
			func(topLevelDeclaration desugar.Struct) {
				structs = append(structs, topLevelDeclaration)
			},
			func(topLevelDeclaration desugar.TypeAlias) {
				typeAliases = append(typeAliases, topLevelDeclaration)
			},
		)
	}
	return declarations, structs, typeAliases
}

func addAllStructFieldsToScope(file string, scope binding.Scope, pkg standard_library.Package) binding.Scope {
	for structName, structWithFields := range pkg.Structs {
		var err *binding.ResolutionError
		scope, err = binding.CopyAddingFields(scope, structWithFields.Struct.Package, desugar.Name{
			String: structName,
		}, structWithFields.Fields)
		if err != nil {
			panic("failed to add standard library struct fields to scope due to: " + err.Error())
		}
	}
	for _, nestedPkg := range pkg.Packages {
		scope = addAllStructFieldsToScope(file, scope, nestedPkg)
	}
	return scope
}

func fallbackOnNil[T any](a *T, b T) T {
	if a != nil {
		return *a
	}
	return b
}

func resolveImports(nodes []desugar.Import, stdLib standard_library.Package, otherPackagesContext *OtherPackagesContext, file string, scope binding.Scope) (map[string]*types.Function, map[string]string, binding.Scope, *type_error.TypecheckError) {
	nativeFunctions := map[string]*types.Function{}
	nativeFunctionPackages := map[string]string{}
	for _, node := range nodes {
		dotSeparatedNames := node.DotSeparatedVars
		as := node.As
		if len(dotSeparatedNames) < 2 {
			errNode := node.Node
			if len(dotSeparatedNames) > 0 {
				errNode = dotSeparatedNames[0].Node
			}
			return nil, nil, nil, type_error.PtrOnNodef(file, errNode, "all interfaces belong to a package")
		}
		foundInOtherPackagesContext := false
		if otherPackagesContext != nil {
			currPackageName := ""
			for i, name := range dotSeparatedNames {
				if i < len(dotSeparatedNames)-1 {
					if i > 0 {
						currPackageName += "."
					}
					currPackageName += name.String
					continue
				}

				otherPackageTypeAlias, ok := otherPackagesContext.TypeAliases[ast.Ref{
					Package: currPackageName,
					Name:    name.String,
				}]
				if ok {
					if as != nil {
						panic("TODO FIXME importing typealias with an alias not yet supported")
					} else {
						updatedScope, err := binding.CopyAddingTypeAliasToFile(scope, file, name, otherPackageTypeAlias.Generics, otherPackageTypeAlias.VariableType)
						if err != nil {
							return nil, nil, nil, type_error.FromResolutionError(file, name.Node, err)
						}
						scope = updatedScope
					}
					foundInOtherPackagesContext = true
					continue
				}

				otherPackageDeclaration, ok := otherPackagesContext.Declarations[ast.Ref{
					Package: currPackageName,
					Name:    name.String,
				}]
				if ok {
					if as != nil {
						updatedScope, err := binding.CopyAddingFileVariable(scope, currPackageName, file, *as, &name, otherPackageDeclaration)
						if err != nil {
							return nil, nil, nil, type_error.FromResolutionError(file, as.Node, err)
						}
						scope = updatedScope
					} else {
						updatedScope, err := binding.CopyAddingFileVariable(scope, currPackageName, file, name, nil, otherPackageDeclaration)
						if err != nil {
							return nil, nil, nil, type_error.FromResolutionError(file, name.Node, err)
						}
						scope = updatedScope
					}
					foundInOtherPackagesContext = true
					continue
				}
				otherPackageStructFunction, ok := otherPackagesContext.StructFunctions[ast.Ref{
					Package: currPackageName,
					Name:    name.String,
				}]
				if ok {
					otherPackageStructFields, ok := otherPackagesContext.FieldsByType[ast.Ref{
						Package: currPackageName,
						Name:    name.String,
					}]
					if !ok {
						panic("got struct but not the fields")
					}
					updatedScope, err := binding.CopyAddingTypeToFile(scope, file, fallbackOnNil(as, name), otherPackageStructFunction.ReturnType)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(file, fallbackOnNil(as, name).Node, err)
					}
					updatedScope, err = binding.CopyAddingFields(updatedScope, currPackageName, fallbackOnNil(as, name), otherPackageStructFields)
					if err != nil {
						return nil, nil, nil, type_error.FromResolutionError(file, fallbackOnNil(as, name).Node, err)
					}
					if as != nil {
						updatedScope, err = binding.CopyAddingFileVariable(updatedScope, currPackageName, file, *as, &name, otherPackageStructFunction)
						if err != nil {
							return nil, nil, nil, type_error.FromResolutionError(file, as.Node, err)
						}
					} else {
						updatedScope, err = binding.CopyAddingFileVariable(updatedScope, currPackageName, file, name, nil, otherPackageStructFunction)
						if err != nil {
							return nil, nil, nil, type_error.FromResolutionError(file, name.Node, err)
						}
					}
					scope = updatedScope
					foundInOtherPackagesContext = true
				}
			}
		}

		if foundInOtherPackagesContext {
			continue
		}

		if dotSeparatedNames[0].String != "tenecs" {
			failedImport := ""
			for i, name := range dotSeparatedNames {
				if i > 0 {
					failedImport += "."
				}
				failedImport += name.String
			}
			return nil, nil, nil, type_error.PtrOnNodef(file, dotSeparatedNames[0].Node, "failed to import "+failedImport+" as it was not found")
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

func validateStructsAndTypeAliases(structsPerFile map[string][]desugar.Struct, typeAliasesInAllFiles map[string][]desugar.TypeAlias, pkgName string, scope binding.Scope) (map[ast.Ref]*types.Function, map[ast.Ref]ast.TypeAlias, binding.Scope, *type_error.TypecheckError) {
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
				return nil, nil, nil, type_error.FromResolutionError(file, node.Name.Node, err)
			}
		}
	}

	resultTypeAliases := map[ast.Ref]ast.TypeAlias{}
	for file, typeAliases := range typeAliasesInAllFiles {
		for _, typeAlias := range typeAliases {
			name := typeAlias.Name
			generics := typeAlias.Generics
			typ := typeAlias.Type
			genericNameStrings := []string{}
			scopeOnlyValidForTypeAlias := scope
			for _, generic := range generics {
				genericNameStrings = append(genericNameStrings, generic.String)
				u, err := binding.CopyAddingTypeToFile(scopeOnlyValidForTypeAlias, file, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return nil, nil, nil, type_error.FromResolutionError(file, generic.Node, err)
				}
				scopeOnlyValidForTypeAlias = u
			}

			varType, err := scopecheck.ValidateTypeAnnotationInScope(typ, file, scopeOnlyValidForTypeAlias)
			if err != nil {
				return nil, nil, nil, type_error.FromScopeCheckError(file, err)
			}

			u, err2 := binding.CopyAddingTypeAliasToAllFiles(scope, name, genericNameStrings, varType)
			if err2 != nil {
				return nil, nil, nil, type_error.FromResolutionError(file, name.Node, err2)
			}
			scope = u
			resultTypeAliases[ast.Ref{
				Package: pkgName,
				Name:    name.String,
			}] = ast.TypeAlias{
				Generics:     genericNameStrings,
				VariableType: varType,
			}
		}
	}

	for file, structsInFile := range structsPerFile {
		for _, node := range structsInFile {
			structName := node.Name
			generics := node.Generics
			parserVariables := node.Variables
			localScope := scope
			for _, generic := range generics {
				u, err := binding.CopyAddingTypeToAllFiles(localScope, generic, &types.TypeArgument{Name: generic.String})
				if err != nil {
					return nil, nil, nil, type_error.FromResolutionError(file, generic.Node, err)
				}
				localScope = u
			}
			constructorArgs := []types.FunctionArgument{}
			variables := map[string]types.VariableType{}
			for _, variable := range parserVariables {
				varType, err := scopecheck.ValidateTypeAnnotationInScope(variable.Type, file, localScope)
				if err != nil {
					return nil, nil, nil, type_error.FromScopeCheckError(file, err)
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
				return nil, nil, nil, type_error.FromResolutionError(file, structName.Node, err)
			}

			genericNames := []types.VariableType{}
			for _, generic := range generics {
				genericNames = append(genericNames, &types.TypeArgument{
					Name: generic.String,
				})
			}
			maybeStruc, resolutionErr := binding.GetTypeByTypeName(localScope, "", structName.String, genericNames)
			if resolutionErr != nil {
				return nil, nil, nil, type_error.FromResolutionError(file, structName.Node, resolutionErr)
			}
			struc, ok := maybeStruc.(*types.KnownType)
			if !ok {
				return nil, nil, nil, type_error.PtrOnNodef(file, structName.Node, "expected struct type in validateStructsAndTypeAliases")
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

	return constructors, resultTypeAliases, scope, nil
}

func TypecheckDeclarations(pkg string, node desugar.Node, declarationsPerFileWithUnderscores map[string][]desugar.Declaration, scope binding.Scope) (map[string]ast.Expression, *type_error.TypecheckError) {
	declarationsPerFile := map[string][]desugar.Declaration{}
	syntheticNameIterator := 0
	for file, declarations := range declarationsPerFileWithUnderscores {
		declarationsPerFile[file] = []desugar.Declaration{}
		for _, declaration := range declarations {
			if declaration.Name.String == "_" {
				declaration.Name.String = fmt.Sprintf("syntheticName_%d", syntheticNameIterator)
				syntheticNameIterator += 1
			}
			declarationsPerFile[file] = append(declarationsPerFile[file], declaration)
		}
	}

	typesByName := map[desugar.Name]types.VariableType{}
	filesByName := map[desugar.Name]string{}

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
