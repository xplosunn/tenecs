package typer

type VariableType interface {
	sealedVariableType()
	Cases() (*Interface, *Function, *BasicType, *Void)
}

type Interface struct {
	Package   string
	Name      string
	Variables map[string]VariableType
}

func (i Interface) sealedVariableType() {}
func (i Interface) Cases() (*Interface, *Function, *BasicType, *Void) {
	return &i, nil, nil, nil
}

type Function struct {
	ArgumentTypes []VariableType
	ReturnType    VariableType
}

func (f Function) sealedVariableType() {}
func (f Function) Cases() (*Interface, *Function, *BasicType, *Void) {
	return nil, &f, nil, nil
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
