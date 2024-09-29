package standard_library

import "github.com/xplosunn/tenecs/typer/types"

type Package struct {
	Packages  map[string]Package
	Structs   map[string]*StructWithFields
	Variables map[string]types.VariableType
}

type StructWithFields struct {
	Struct           *types.KnownType
	Fields           map[string]types.VariableType
	FieldNamesSorted []string
}

func packageWith(opts ...func(*Package)) Package {
	pkg := &Package{
		Packages:  map[string]Package{},
		Structs:   map[string]*StructWithFields{},
		Variables: map[string]types.VariableType{},
	}
	for _, opt := range opts {
		opt(pkg)
	}
	return *pkg
}

func withPackage(name string, pack Package) func(pkg *Package) {
	return func(pkg *Package) {
		pkg.Packages[name] = pack
	}
}

func withStruct(name string, struc *types.KnownType, fieldFuncs ...func(*StructWithFields)) func(pkg *Package) {
	return func(pkg *Package) {
		result := &StructWithFields{
			Struct:           struc,
			Fields:           map[string]types.VariableType{},
			FieldNamesSorted: []string{},
		}
		for _, f := range fieldFuncs {
			f(result)
		}
		pkg.Structs[name] = result
	}
}

func structField(name string, varType types.VariableType) func(*StructWithFields) {
	return func(structWithFields *StructWithFields) {
		structWithFields.FieldNamesSorted = append(structWithFields.FieldNamesSorted, name)
		structWithFields.Fields[name] = varType
	}
}

func withFunction(name string, function *types.Function) func(pkg *Package) {
	return func(pkg *Package) {
		pkg.Variables[name] = function
	}
}

type NamedFunction struct {
	name     string
	function *types.Function
}

func withFunctions(functions []NamedFunction) func(pkg *Package) {
	return func(pkg *Package) {
		for _, f := range functions {
			pkg.Variables[f.name] = f.function
		}
	}
}
