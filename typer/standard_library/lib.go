package standard_library

import (
	"github.com/xplosunn/tenecs/parser"
	"github.com/xplosunn/tenecs/typer/binding"
	"github.com/xplosunn/tenecs/typer/scopecheck"
	"github.com/xplosunn/tenecs/typer/types"
	"reflect"
	"strings"
	"testing"
)

var StdLib = Package{
	Packages: topLevelPackages,
}

var DefaultTypesAvailableWithoutImport = map[string]types.VariableType{
	"String":  types.String(),
	"Float":   types.Float(),
	"Int":     types.Int(),
	"Boolean": types.Boolean(),
	"Void":    types.Void(),
	"List": &types.List{
		Generic: &types.TypeArgument{
			Name: "T",
		},
	},
}

var topLevelPackages = map[string]Package{
	"tenecs": packageWith(
		withPackage("list", tenecs_list),
		withPackage("boolean", tenecs_boolean),
		withPackage("compare", tenecs_compare),
		withPackage("error", tenecs_error),
		withPackage("int", tenecs_int),
		withPackage("json", tenecs_json),
		withPackage("go", tenecs_go),
		withPackage("ref", tenecs_ref),
		withPackage("string", tenecs_string),
		withPackage("test", tenecs_test),
		withPackage("time", tenecs_time),
		withPackage("web", tenecs_web),
	),
}

func StdLibGetOrPanic(t *testing.T, ref string) *types.KnownType {
	pkg := StdLib
	split := strings.Split(ref, ".")
	var finalName string
	for i, name := range split {
		if i < len(split)-1 {
			pkg = pkg.Packages[name]
		} else {
			finalName = name
		}
	}
	if pkg.Structs[finalName] == nil {
		t.Fatal("StdLibGetOrPanic" + ref)
	}
	return pkg.Structs[finalName].Struct
}

func StdLibGetFunctionOrPanic(t *testing.T, ref string) *types.Function {
	pkg := StdLib
	split := strings.Split(ref, ".")
	var finalName string
	for i, name := range split {
		if i < len(split)-1 {
			pkg = pkg.Packages[name]
		} else {
			finalName = name
		}
	}
	if pkg.Structs[finalName] == nil {
		t.Fatal("StdLibGetOrPanic" + ref)
	}
	arguments := []types.FunctionArgument{}
	for _, fieldName := range pkg.Structs[finalName].FieldNamesSorted {
		arguments = append(arguments, types.FunctionArgument{
			Name:         fieldName,
			VariableType: pkg.Structs[finalName].Fields[fieldName],
		})
	}
	if len(pkg.Structs[finalName].Struct.Generics) > 0 {
		panic("todo StdLibGetFunctionOrPanic with generics")
	}
	return &types.Function{
		Generics:   nil,
		Arguments:  arguments,
		ReturnType: pkg.Structs[finalName].Struct,
	}
}

// When a function returns a pointer to a specific struct implementing 'error' and it's assigned to a variable of type 'error'
// then the nil check will no longer work.
//
// Example:
//
//	parsed, err := parser.ParseFunctionTypeString(signature) // returns (..., error)
//	if err != nil { // this one is ok
//		panic(err)
//	}
//	returnType, err := scopecheck.ValidateTypeAnnotationInScope(*parsed.ReturnType, "non_existing_file", scope) // returns (..., *ScopeCheckError)
//	if err != nil { // this one will return true even when err is nil
//		panic(err)
//	}
//
// See https://www.reddit.com/r/golang/comments/1bu5r72/subtle_and_surprising_behavior_when_interface/
func isNil(val any) bool {
	if val == nil {
		return true
	}

	v := reflect.ValueOf(val)
	k := v.Kind()
	switch k {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer,
		reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return v.IsNil()
	}

	return false
}

func functionFromType(signature string, structsUsedInSignature ...*StructWithFields) *types.Function {
	parsed, err := parser.ParseFunctionTypeString(signature)
	if err != nil {
		panic("functionFromType failed to parse '" + signature + "' due to " + err.Error())
	}
	scope := binding.NewFromDefaults(DefaultTypesAvailableWithoutImport)
	for _, structUsed := range structsUsedInSignature {
		scope, err = binding.CopyAddingTypeToAllFiles(scope, parser.Name{String: structUsed.Struct.Name}, structUsed.Struct)
		if !isNil(err) {
			panic("functionFromType failed to add '" + structUsed.Struct.Name + "' struct to scope due to " + err.Error())
		}
	}

	generics := []string{}
	for _, generic := range parsed.Generics {
		generics = append(generics, generic.String)
		scope, err = binding.CopyAddingTypeToAllFiles(scope, generic, &types.TypeArgument{Name: generic.String})
		if !isNil(err) {
			panic("functionFromType failed to add '" + generic.String + "' generic to scope due to " + err.Error())
		}
	}

	arguments := []types.FunctionArgument{}
	for _, param := range parsed.Arguments {
		varType, err := scopecheck.ValidateTypeAnnotationInScope(param.Type, "non_existing_file", scope)
		if err != nil {
			panic("functionFromType failed to ValidateTypeAnnotationInScope for '" + signature + "' due to " + err.Error())
		}
		arguments = append(arguments, types.FunctionArgument{
			Name:         param.Name.String,
			VariableType: varType,
		})
	}
	
	returnType, err := scopecheck.ValidateTypeAnnotationInScope(parsed.ReturnType, "non_existing_file", scope)
	if err != nil {
		panic("functionFromType failed to ValidateTypeAnnotationInScope for '" + signature + "' due to " + err.Error())
	}

	return &types.Function{
		Generics:   generics,
		Arguments:  arguments,
		ReturnType: returnType,
	}
}
