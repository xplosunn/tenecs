package binding

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/types"
)

type Scope interface {
	impl() *scopeImpl
}

type packageAndAliasFor struct {
	pkg      string
	aliasFor *string
}

type typeAlias struct {
	generics     []string
	variableType types.VariableType
}

type scopeImpl struct {
	TypeAliasByTypeName        *immutable.Map[string, typeAlias]
	TypeByTypeName             TwoLevelMap[string, string, types.VariableType]
	FieldsByTypeName           *immutable.Map[string, map[string]types.VariableType]
	TypeByVariableName         TwoLevelMap[string, string, types.VariableType]
	PackageLevelByVariableName TwoLevelMap[string, string, packageAndAliasFor]
}

func (u scopeImpl) impl() *scopeImpl {
	return &u
}

func NewFromDefaults(defaultTypesWithoutImport map[string]types.VariableType) Scope {
	mapBuilder := NewTwoLevelMap[string, string, types.VariableType]()
	var ok bool
	for key, value := range defaultTypesWithoutImport {
		mapBuilder, ok = mapBuilder.SetGlobalIfAbsent(key, value)
		if !ok {
			panic("repeat type in std lib " + key)
		}
	}
	return scopeImpl{
		TypeAliasByTypeName:        immutable.NewMap[string, typeAlias](nil),
		TypeByTypeName:             mapBuilder,
		FieldsByTypeName:           immutable.NewMap[string, map[string]types.VariableType](nil),
		TypeByVariableName:         NewTwoLevelMap[string, string, types.VariableType](),
		PackageLevelByVariableName: NewTwoLevelMap[string, string, packageAndAliasFor](),
	}
}

func GetTypeByTypeName(scope Scope, file string, typeName string, generics []types.VariableType) (types.VariableType, *ResolutionError) {
	u := scope.impl()

	alias, ok := u.TypeAliasByTypeName.Get(typeName)
	if ok {
		if len(generics) != len(alias.generics) {
			return nil, ResolutionErrorWrongNumberOfGenerics(alias.variableType, len(alias.generics), len(generics))
		}
		varType := alias.variableType
		for i, generic := range alias.generics {
			resolved, err := ResolveGeneric(varType, generic, generics[i])
			if err != nil {
				return nil, err
			}
			varType = resolved
		}
		return varType, nil
	}

	varType, ok := u.TypeByTypeName.Get(file, typeName)
	if ok {
		return applyGenerics(varType, generics)

	}
	return nil, ResolutionErrorCouldNotResolve(typeName)
}

func applyGenerics(varType types.VariableType, genericArgs []types.VariableType) (types.VariableType, *ResolutionError) {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		if len(genericArgs) != 0 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 0, len(genericArgs))
		}
		return varType, nil
	} else if caseList != nil {
		if len(genericArgs) != 1 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 1, len(genericArgs))
		}
		return &types.List{
			Generic: genericArgs[0],
		}, nil
	} else if caseKnownType != nil {
		if len(genericArgs) != len(caseKnownType.Generics) {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, len(caseKnownType.Generics), len(genericArgs))
		}

		if len(genericArgs) == 0 {
			genericArgs = nil
		}

		return &types.KnownType{
			Package:          caseKnownType.Package,
			Name:             caseKnownType.Name,
			DeclaredGenerics: caseKnownType.DeclaredGenerics,
			Generics:         genericArgs,
		}, nil
	} else if caseFunction != nil {
		if len(caseFunction.Generics) != len(genericArgs) {
			panic("TODO nicer error message")
		}
		resolvedArguments := []types.FunctionArgument{}
		for _, argument := range caseFunction.Arguments {
			resolvedArguments = append(resolvedArguments, argument)
		}
		for i, argument := range resolvedArguments {
			resolvedArgumentType := argument.VariableType
			for i, generic := range caseFunction.Generics {
				resolved, err := ResolveGeneric(resolvedArgumentType, generic, genericArgs[i])
				if err != nil {
					return nil, err
				}
				resolvedArgumentType = resolved
			}
			resolvedArguments[i] = types.FunctionArgument{
				Name:         argument.Name,
				VariableType: resolvedArgumentType,
			}
		}

		resolvedReturnType := caseFunction.ReturnType
		for i, generic := range caseFunction.Generics {
			resolved, err := ResolveGeneric(resolvedReturnType, generic, genericArgs[i])
			if err != nil {
				return nil, err
			}
			resolvedReturnType = resolved
		}
		return &types.Function{
			Generics:   nil,
			Arguments:  resolvedArguments,
			ReturnType: resolvedReturnType,
		}, nil
	} else if caseOr != nil {
		panic("unexpected applyGenerics caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func ResolveGeneric(over types.VariableType, genericName string, resolveWith types.VariableType) (types.VariableType, *ResolutionError) {
	caseTypeArgument, caseList, caseKnownType, caseFunction, caseOr := over.VariableTypeCases()
	if caseTypeArgument != nil {
		if caseTypeArgument.Name == genericName {
			return resolveWith, nil
		}
		return caseTypeArgument, nil
	} else if caseList != nil {
		resolvedGeneric, err := ResolveGeneric(caseList.Generic, genericName, resolveWith)
		if err != nil {
			return nil, err
		}
		return &types.List{
			Generic: resolvedGeneric,
		}, nil
	} else if caseKnownType != nil {
		newGenerics := []types.VariableType{}
		for _, genericVarType := range caseKnownType.Generics {
			resolvedGeneric, err := ResolveGeneric(genericVarType, genericName, resolveWith)
			if err != nil {
				return nil, err
			}
			newGenerics = append(newGenerics, resolvedGeneric)
		}
		newKnownType := &types.KnownType{
			Package:          caseKnownType.Package,
			Name:             caseKnownType.Name,
			DeclaredGenerics: caseKnownType.DeclaredGenerics,
			Generics:         newGenerics,
		}
		return newKnownType, nil
	} else if caseFunction != nil {
		arguments := []types.FunctionArgument{}
		for _, argument := range caseFunction.Arguments {
			varType, err := ResolveGeneric(argument.VariableType, genericName, resolveWith)
			if err != nil {
				return nil, err
			}
			arguments = append(arguments, types.FunctionArgument{
				Name:         argument.Name,
				VariableType: varType,
			})
		}
		returnType, err := ResolveGeneric(caseFunction.ReturnType, genericName, resolveWith)
		if err != nil {
			return nil, err
		}
		return &types.Function{
			Generics:   caseFunction.Generics,
			Arguments:  arguments,
			ReturnType: returnType,
		}, nil
	} else if caseOr != nil {
		resolvedOr := &types.OrVariableType{Elements: []types.VariableType{}}
		for _, elem := range caseOr.Elements {
			resolved, err := ResolveGeneric(elem, genericName, resolveWith)
			if err != nil {
				return nil, err
			}
			types.VariableTypeAddToOr(resolved, resolvedOr)
		}
		return resolvedOr, nil
	} else {
		panic(fmt.Errorf("cases on %v", over))
	}
}

func GetFields(scope Scope, knownType *types.KnownType) (map[string]types.VariableType, *ResolutionError) {
	u := scope.impl()
	fields, ok := u.FieldsByTypeName.Get(knownType.Package + "~>" + knownType.Name)
	if !ok {
		return map[string]types.VariableType{}, nil
	}
	fieldsWithResolvedGenerics := map[string]types.VariableType{}
	for k, v := range fields {
		fieldsWithResolvedGenerics[k] = v
	}

	for i, resolveWith := range knownType.Generics {
		for fieldName, fieldVarType := range fieldsWithResolvedGenerics {
			resolved, err := ResolveGeneric(fieldVarType, knownType.DeclaredGenerics[i], resolveWith)
			if err != nil {
				return nil, err
			}
			fieldsWithResolvedGenerics[fieldName] = resolved
		}
	}

	return fieldsWithResolvedGenerics, nil
}

func GetAllFields(scope Scope) map[string]map[string]types.VariableType {
	u := scope.impl()

	result := map[string]map[string]types.VariableType{}
	iterator := u.FieldsByTypeName.Iterator()
	for !iterator.Done() {
		key, value, _ := iterator.Next()
		result[key] = value
	}
	return result
}

func GetTypeByVariableName(scope Scope, file string, variableName string) (types.VariableType, bool) {
	u := scope.impl()
	return u.TypeByVariableName.Get(file, variableName)
}

func CopyAddingTypeToFile(scope Scope, file string, typeName parser.Name, varType types.VariableType) (Scope, *ResolutionError) {
	u := scope.impl()
	m, ok := u.TypeByTypeName.SetScopedIfAbsent(file, typeName.String, varType)
	if !ok {
		return nil, ResolutionErrorTypeAlreadyExists(varType)
	}
	return scopeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             m,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func CopyAddingTypeToAllFiles(scope Scope, typeName parser.Name, varType types.VariableType) (Scope, *ResolutionError) {
	u := scope.impl()
	m, ok := u.TypeByTypeName.SetGlobalIfAbsent(typeName.String, varType)
	if !ok {
		return nil, ResolutionErrorTypeAlreadyExists(varType)
	}
	return scopeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             m,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func CopyAddingTypeAliasToAllFiles(scope Scope, typeName parser.Name, generics []string, varType types.VariableType) (Scope, *ResolutionError) {
	u := scope.impl()
	_, ok := u.TypeAliasByTypeName.Get(typeName.String)
	if ok {
		return nil, ResolutionErrorTypeAlreadyExists(varType)
	}
	return scopeImpl{
		TypeAliasByTypeName: u.TypeAliasByTypeName.Set(typeName.String, typeAlias{
			generics:     generics,
			variableType: varType,
		}),
		TypeByTypeName:             u.TypeByTypeName,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func CopyAddingFields(scope Scope, packageName string, typeName parser.Name, fields map[string]types.VariableType) (Scope, *ResolutionError) {
	u := scope.impl()
	_, ok := u.FieldsByTypeName.Get(typeName.String)
	if ok {
		return nil, ResolutionErrorTypeFieldsAlreadyExists(typeName.String)
	}
	return scopeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             u.TypeByTypeName,
		FieldsByTypeName:           u.FieldsByTypeName.Set(packageName+"~>"+typeName.String, fields),
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func copyAddingVariable(isPackageLevel *string, isFileLevel *string, scope Scope, variableName parser.Name, aliasFor *parser.Name, varType types.VariableType) (Scope, *ResolutionError) {
	if isFileLevel != nil && isPackageLevel == nil {
		panic("misuse of copyAddingVariable")
	}
	if variableName.String == "_" {
		return scope, nil
	}
	u := scope.impl()

	if aliasFor != nil && isFileLevel == nil {
		panic("copyAddingVariable with alias should be done on file level")
	}

	packageLevelByVariableName := u.PackageLevelByVariableName
	if isPackageLevel != nil {
		var aliasForStr *string = nil
		if aliasFor != nil {
			aliasForStr = &aliasFor.String
		}
		ok := true
		if isFileLevel != nil {
			file := *isFileLevel
			packageLevelByVariableName, ok = u.PackageLevelByVariableName.SetScopedIfAbsent(file, variableName.String, packageAndAliasFor{
				pkg:      *isPackageLevel,
				aliasFor: aliasForStr,
			})
		} else {
			packageLevelByVariableName, ok = u.PackageLevelByVariableName.SetGlobalIfAbsent(variableName.String, packageAndAliasFor{
				pkg:      *isPackageLevel,
				aliasFor: aliasForStr,
			})
		}
		if !ok {
			return nil, ResolutionErrorVariableAlreadyExists(varType, variableName.String)
		}

	}
	typeByVariableName := u.TypeByVariableName
	ok := true
	if isFileLevel != nil {
		file := *isFileLevel
		typeByVariableName, ok = u.TypeByVariableName.SetScopedIfAbsent(file, variableName.String, varType)
	} else {
		typeByVariableName, ok = u.TypeByVariableName.SetGlobalIfAbsent(variableName.String, varType)
	}
	if !ok {
		return nil, ResolutionErrorVariableAlreadyExists(varType, variableName.String)
	}
	return scopeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             u.TypeByTypeName,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         typeByVariableName,
		PackageLevelByVariableName: packageLevelByVariableName,
	}, nil
}

func CopyAddingPackageVariable(scope Scope, pkgName string, variableName parser.Name, varType types.VariableType) (Scope, *ResolutionError) {
	return copyAddingVariable(&pkgName, nil, scope, variableName, nil, varType)
}

func CopyAddingFileVariable(scope Scope, pkgName string, file string, variableName parser.Name, aliasFor *parser.Name, varType types.VariableType) (Scope, *ResolutionError) {
	return copyAddingVariable(&pkgName, &file, scope, variableName, aliasFor, varType)
}

func CopyAddingLocalVariable(scope Scope, variableName parser.Name, varType types.VariableType) (Scope, *ResolutionError) {
	return copyAddingVariable(nil, nil, scope, variableName, nil, varType)
}

func GetPackageLevelAndUnaliasedNameOfVariable(scope Scope, file string, variableName parser.Name) (*string, string) {
	u := scope.impl()
	result, ok := u.PackageLevelByVariableName.Get(file, variableName.String)
	if ok {
		if result.aliasFor != nil {
			return &result.pkg, *result.aliasFor
		} else {
			return &result.pkg, variableName.String
		}
	} else {
		return nil, variableName.String
	}
}
