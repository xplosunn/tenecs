package binding

import (
	"fmt"
	"github.com/benbjohnson/immutable"
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/type_error"
	"github.com/xplosunn/tenecs/typer/types"
)

type Universe interface {
	impl() *universeImpl
}

type packageAndAliasFor struct {
	pkg      string
	aliasFor *string
}

type typeAlias struct {
	generics     []string
	variableType types.VariableType
}

type universeImpl struct {
	TypeAliasByTypeName        immutable.Map[string, typeAlias]
	TypeByTypeName             TwoLevelMap[string, string, types.VariableType]
	FieldsByTypeName           immutable.Map[string, map[string]types.VariableType]
	TypeByVariableName         immutable.Map[string, types.VariableType]
	PackageLevelByVariableName immutable.Map[string, packageAndAliasFor]
}

func (u universeImpl) impl() *universeImpl {
	return &u
}

func NewFromDefaults(defaultTypesWithoutImport map[string]types.VariableType) Universe {
	mapBuilder := NewTwoLevelMap[string, string, types.VariableType]()
	var ok bool
	for key, value := range defaultTypesWithoutImport {
		mapBuilder, ok = mapBuilder.SetGlobalIfAbsent(key, value)
		if !ok {
			panic("repeat type in std lib " + key)
		}
	}
	return universeImpl{
		TypeAliasByTypeName:        *immutable.NewMap[string, typeAlias](nil),
		TypeByTypeName:             mapBuilder,
		FieldsByTypeName:           *immutable.NewMap[string, map[string]types.VariableType](nil),
		TypeByVariableName:         *immutable.NewMap[string, types.VariableType](nil),
		PackageLevelByVariableName: *immutable.NewMap[string, packageAndAliasFor](nil),
	}
}

func GetTypeByTypeName(universe Universe, file string, typeName string, generics []types.VariableType) (types.VariableType, *ResolutionError) {
	u := universe.impl()

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
		return ApplyGenerics(varType, generics)

	}
	return nil, ResolutionErrorCouldNotResolve(typeName)
}

func ApplyGenerics(varType types.VariableType, genericArgs []types.VariableType) (types.VariableType, *ResolutionError) {
	caseTypeArgument, caseKnownType, caseFunction, caseOr := varType.VariableTypeCases()
	if caseTypeArgument != nil {
		if len(genericArgs) != 0 {
			return nil, ResolutionErrorWrongNumberOfGenerics(varType, 0, len(genericArgs))
		}
		return varType, nil
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
			ValidStructField: caseKnownType.ValidStructField,
		}, nil
	} else if caseFunction != nil {
		panic("TODO ApplyGenerics caseFunction")
	} else if caseOr != nil {
		panic("unexpected ApplyGenerics caseOr")
	} else {
		panic(fmt.Errorf("cases on %v", varType))
	}
}

func ResolveGeneric(over types.VariableType, genericName string, resolveWith types.VariableType) (types.VariableType, *ResolutionError) {
	if !resolveWith.CanBeStructField() {
		return nil, ResolutionErrorNotAValidGeneric(resolveWith)
	}
	caseTypeArgument, caseKnownType, caseFunction, caseOr := over.VariableTypeCases()
	if caseTypeArgument != nil {
		if caseTypeArgument.Name == genericName {
			return resolveWith, nil
		}
		return caseTypeArgument, nil
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
			ValidStructField: caseKnownType.ValidStructField,
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

func GetFields(universe Universe, knownType *types.KnownType) (map[string]types.VariableType, *ResolutionError) {
	u := universe.impl()
	fields, ok := u.FieldsByTypeName.Get(knownType.Package + "->" + knownType.Name)
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

func GetAllFields(universe Universe) map[string]map[string]types.VariableType {
	u := universe.impl()

	result := map[string]map[string]types.VariableType{}
	iterator := u.FieldsByTypeName.Iterator()
	for !iterator.Done() {
		key, value, _ := iterator.Next()
		result[key] = value
	}
	return result
}

func GetTypeByVariableName(universe Universe, variableName string) (types.VariableType, bool) {
	u := universe.impl()
	return u.TypeByVariableName.Get(variableName)
}

func CopyAddingTypeToFile(universe Universe, file string, typeName parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	m, ok := u.TypeByTypeName.SetScopedIfAbsent(file, typeName.String, varType)
	if !ok {
		return nil, type_error.PtrOnNodef(typeName.Node, "type already exists %s", typeName.String)
	}
	return universeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             m,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func CopyAddingTypeToAllFiles(universe Universe, typeName parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	m, ok := u.TypeByTypeName.SetGlobalIfAbsent(typeName.String, varType)
	if !ok {
		return nil, type_error.PtrOnNodef(typeName.Node, "type already exists %s", typeName.String)
	}
	return universeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             m,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func CopyAddingTypeAliasToAllFiles(universe Universe, typeName parser.Name, generics []string, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeAliasByTypeName.Get(typeName.String)
	if ok {
		return nil, type_error.PtrOnNodef(typeName.Node, "type already exists %s", typeName.String)
	}
	return universeImpl{
		TypeAliasByTypeName: *u.TypeAliasByTypeName.Set(typeName.String, typeAlias{
			generics:     generics,
			variableType: varType,
		}),
		TypeByTypeName:             u.TypeByTypeName,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func CopyAddingFields(universe Universe, packageName string, typeName parser.Name, fields map[string]types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.FieldsByTypeName.Get(typeName.String)
	if ok {
		return nil, type_error.PtrOnNodef(typeName.Node, "type fields already exist: %s", typeName.String)
	}
	return universeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             u.TypeByTypeName,
		FieldsByTypeName:           *u.FieldsByTypeName.Set(packageName+"->"+typeName.String, fields),
		TypeByVariableName:         u.TypeByVariableName,
		PackageLevelByVariableName: u.PackageLevelByVariableName,
	}, nil
}

func copyAddingVariable(isPackageLevel *string, universe Universe, variableName parser.Name, aliasFor *parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	u := universe.impl()
	_, ok := u.TypeByVariableName.Get(variableName.String)
	if ok {
		return nil, type_error.PtrOnNodef(variableName.Node, "duplicate variable '%s'", variableName.String)
	}
	if aliasFor != nil && isPackageLevel == nil {
		panic("copyAddingVariable with alias should be done on package level")
	}
	packageLevelByVariableName := u.PackageLevelByVariableName
	if isPackageLevel != nil {
		var aliasForStr *string = nil
		if aliasFor != nil {
			aliasForStr = &aliasFor.String
		}
		packageLevelByVariableName = *u.PackageLevelByVariableName.Set(variableName.String, packageAndAliasFor{
			pkg:      *isPackageLevel,
			aliasFor: aliasForStr,
		})
	}
	return universeImpl{
		TypeAliasByTypeName:        u.TypeAliasByTypeName,
		TypeByTypeName:             u.TypeByTypeName,
		FieldsByTypeName:           u.FieldsByTypeName,
		TypeByVariableName:         *u.TypeByVariableName.Set(variableName.String, varType),
		PackageLevelByVariableName: packageLevelByVariableName,
	}, nil
}

func CopyAddingPackageVariable(universe Universe, pkgName string, variableName parser.Name, aliasFor *parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	return copyAddingVariable(&pkgName, universe, variableName, aliasFor, varType)
}

func CopyAddingLocalVariable(universe Universe, variableName parser.Name, varType types.VariableType) (Universe, *type_error.TypecheckError) {
	return copyAddingVariable(nil, universe, variableName, nil, varType)
}

func GetPackageLevelAndUnaliasedNameOfVariable(universe Universe, variableName parser.Name) (*string, string) {
	u := universe.impl()
	result, ok := u.PackageLevelByVariableName.Get(variableName.String)
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
