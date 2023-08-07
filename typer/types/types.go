package types

/*
There are different categories of types we care about:
1. Functions -> can construct types
2. TypeArgument -> can only happen when there's an unresolved generic in scope
3. "concrete" types
*/

type VariableType interface {
	sealedVariableType()
	VariableTypeCases() (*TypeArgument, *KnownType, *Function, *OrVariableType)
	CanBeStructField() bool
}

type TypeArgument struct {
	Name string
}

func (t *TypeArgument) sealedVariableType() {}
func (t *TypeArgument) VariableTypeCases() (*TypeArgument, *KnownType, *Function, *OrVariableType) {
	return t, nil, nil, nil
}
func (t *TypeArgument) CanBeStructField() bool {
	return true
}

type KnownType struct {
	Package          string
	Name             string
	DeclaredGenerics []string
	Generics         []VariableType
	ValidStructField bool
}

func (k *KnownType) sealedVariableType() {}
func (k *KnownType) VariableTypeCases() (*TypeArgument, *KnownType, *Function, *OrVariableType) {
	return nil, k, nil, nil
}
func (k *KnownType) CanBeStructField() bool {
	return k.ValidStructField
}

type Function struct {
	Generics   []string
	Arguments  []FunctionArgument
	ReturnType VariableType
}

type FunctionArgument struct {
	Name         string
	VariableType VariableType
}

func (f *Function) sealedVariableType() {}
func (f *Function) VariableTypeCases() (*TypeArgument, *KnownType, *Function, *OrVariableType) {
	return nil, nil, f, nil
}
func (f *Function) CanBeStructField() bool {
	return false
}

type OrVariableType struct {
	Elements []VariableType
}

func (o *OrVariableType) sealedVariableType() {}
func (o *OrVariableType) VariableTypeCases() (*TypeArgument, *KnownType, *Function, *OrVariableType) {
	return nil, nil, nil, o
}
func (o *OrVariableType) CanBeStructField() bool {
	for _, elem := range o.Elements {
		if !elem.CanBeStructField() {
			return false
		}
	}
	return true
}

func String() *KnownType  { return basicType("String") }
func Float() *KnownType   { return basicType("Float") }
func Int() *KnownType     { return basicType("Int") }
func Boolean() *KnownType { return basicType("Boolean") }
func Void() *KnownType    { return basicType("Void") }

func basicType(name string) *KnownType {
	return &KnownType{
		Package:          "",
		Name:             name,
		DeclaredGenerics: nil,
		Generics:         nil,
		ValidStructField: true,
	}
}

func Interface(pkg string, name string, generics []string) *KnownType {
	genericVarTypes := []VariableType{}
	for _, generic := range generics {
		genericVarTypes = append(genericVarTypes, &TypeArgument{Name: generic})
	}
	return &KnownType{
		Package:          pkg,
		Name:             name,
		DeclaredGenerics: generics,
		Generics:         genericVarTypes,
		ValidStructField: false,
	}
}

func UncheckedApplyGenerics(to *KnownType, generics []VariableType) *KnownType {
	if len(generics) != len(to.DeclaredGenerics) {
		panic("Tried UncheckedApplyGenerics but provided wrong number of generics")
	}
	for _, generic := range generics {
		if !generic.CanBeStructField() {
			panic("Tried UncheckedApplyGenerics but provided an invalid generic")
		}
	}
	return &KnownType{
		Package:          to.Package,
		Name:             to.Name,
		DeclaredGenerics: to.DeclaredGenerics,
		Generics:         generics,
		ValidStructField: to.ValidStructField,
	}
}

func UncheckedArray(of VariableType) *KnownType {
	array, ok := Array(of)
	if !ok {
		panic("Tried to make unchecked array")
	}
	return array
}

func Array(of VariableType) (*KnownType, bool) {
	if !of.CanBeStructField() {
		return nil, false
	}
	return &KnownType{
		Package:          "",
		Name:             "Array",
		DeclaredGenerics: []string{"T"},
		Generics: []VariableType{
			of,
		},
		ValidStructField: true,
	}, true
}
