package types

type VariableType interface {
	sealedVariableType()
	VariableTypeCases() (*TypeArgument, *List, *KnownType, *Function, *OrVariableType)
}

type TypeArgument struct {
	Name string
}

func (t *TypeArgument) sealedVariableType() {}
func (t *TypeArgument) VariableTypeCases() (*TypeArgument, *List, *KnownType, *Function, *OrVariableType) {
	return t, nil, nil, nil, nil
}

type List struct {
	Generic VariableType
}

func (l *List) sealedVariableType() {}
func (l *List) VariableTypeCases() (*TypeArgument, *List, *KnownType, *Function, *OrVariableType) {
	return nil, l, nil, nil, nil
}

type KnownType struct {
	Package          string
	Name             string
	DeclaredGenerics []string
	Generics         []VariableType
}

func (k *KnownType) sealedVariableType() {}
func (k *KnownType) VariableTypeCases() (*TypeArgument, *List, *KnownType, *Function, *OrVariableType) {
	return nil, nil, k, nil, nil
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
func (f *Function) VariableTypeCases() (*TypeArgument, *List, *KnownType, *Function, *OrVariableType) {
	return nil, nil, nil, f, nil
}

type OrVariableType struct {
	Elements []VariableType
}

func (o *OrVariableType) sealedVariableType() {}
func (o *OrVariableType) VariableTypeCases() (*TypeArgument, *List, *KnownType, *Function, *OrVariableType) {
	return nil, nil, nil, nil, o
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
	}
}

func Struct(pkg string, name string, generics []string) *KnownType {
	genericVarTypes := []VariableType{}
	for _, generic := range generics {
		genericVarTypes = append(genericVarTypes, &TypeArgument{Name: generic})
	}
	return &KnownType{
		Package:          pkg,
		Name:             name,
		DeclaredGenerics: generics,
		Generics:         genericVarTypes,
	}
}

func UncheckedApplyGenerics(to *KnownType, generics []VariableType) *KnownType {
	if len(generics) != len(to.DeclaredGenerics) {
		panic("Tried UncheckedApplyGenerics but provided wrong number of generics")
	}
	return &KnownType{
		Package:          to.Package,
		Name:             to.Name,
		DeclaredGenerics: to.DeclaredGenerics,
		Generics:         generics,
	}
}
