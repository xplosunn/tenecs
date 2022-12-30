package types

type VariableType interface {
	sealedVariableType()
	Cases() (*Interface, *Function, *BasicType, *Void)
}

type Interface struct {
	Package string
	Name    string
}

func (i Interface) sealedVariableType() {}
func (i Interface) Cases() (*Interface, *Function, *BasicType, *Void) {
	return &i, nil, nil, nil
}

type Function struct {
	Arguments  []FunctionArgument
	ReturnType VariableType
}

func (f Function) sealedVariableType() {}
func (f Function) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, &f, nil, nil
}

type FunctionArgument struct {
	Name         string
	VariableType VariableType
}

type BasicType struct {
	Type string
}

func (b BasicType) sealedVariableType() {}
func (b BasicType) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, nil, &b, nil
}

type Void struct {
}

func (v Void) sealedVariableType() {}
func (v Void) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, nil, nil, &v
}
